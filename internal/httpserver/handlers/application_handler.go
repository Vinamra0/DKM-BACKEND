package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"dkmbackend/internal/models"
)

type ApplicationHandler struct {
	mu       sync.Mutex
	storage  string
	metaFile string
}

func NewApplicationHandler() *ApplicationHandler {
	h := &ApplicationHandler{
		storage:  "storage/cvs",
		metaFile: "storage/applications.json",
	}
	_ = os.MkdirAll(h.storage, 0o755)
	return h
}

var allowedMIMEs = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

var allowedExt = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
}

var sanitizeRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeFileName(name string) string {
	base := filepath.Base(name)
	return sanitizeRe.ReplaceAllString(base, "-")
}

func (h *ApplicationHandler) loadAll() ([]models.Application, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := os.Stat(h.metaFile); errors.Is(err, os.ErrNotExist) {
		return []models.Application{}, nil
	}
	b, err := os.ReadFile(h.metaFile)
	if err != nil {
		return nil, err
	}
	var apps []models.Application
	if err := json.Unmarshal(b, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func (h *ApplicationHandler) saveAll(apps []models.Application) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(h.metaFile), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.metaFile, b, 0o644)
}

// POST /api/applications
// multipart/form-data fields: name,email,phone,education,experience,location,coverLetter,jobId and file field cv
func (h *ApplicationHandler) Submit(w http.ResponseWriter, r *http.Request) {
	// limit overall request size to ~12MB
	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
	if err := r.ParseMultipartForm(12 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	if name == "" || email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "name and email are required"})
		return
	}

	file, fh, err := r.FormFile("cv")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "cv is required"})
		return
	}
	defer file.Close()

	// sniff MIME from first 512 bytes
	head := make([]byte, 512)
	n, _ := file.Read(head)
	mime := http.DetectContentType(head[:n])
	if !allowedMIMEs[mime] {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "cv must be PDF or Word document"})
		return
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExt[ext] {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid file extension"})
		return
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	id := uuid.NewString()
	safe := sanitizeFileName(fh.Filename)
	stored := ts + "-" + id + "-" + safe
	destPath := filepath.Join(h.storage, stored)

	out, err := os.Create(destPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to save file"})
		return
	}
	defer out.Close()

	// include head bytes then the rest, but cap at 10MB
	reader := io.MultiReader(bytes.NewReader(head[:n]), file)
	lr := io.LimitReader(reader, 10<<20+1)
	written, err := io.Copy(out, lr)
	if err != nil {
		_ = os.Remove(destPath)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to write file"})
		return
	}
	if written > 10<<20 {
		_ = os.Remove(destPath)
		writeJSON(w, http.StatusRequestEntityTooLarge, map[string]interface{}{"success": false, "message": "cv exceeds 10MB"})
		return
	}

	app := models.Application{
		ID:             id,
		Name:           name,
		Email:          email,
		Phone:          r.FormValue("phone"),
		Education:      r.FormValue("education"),
		Experience:     r.FormValue("experience"),
		Location:       r.FormValue("location"),
		CoverLetter:    r.FormValue("coverLetter"),
		JobID:          r.FormValue("jobId"),
		CVOriginalName: fh.Filename,
		CVStoredName:   stored,
		CVMimeType:     mime,
		CVSize:         written,
		AppliedAt:      time.Now().UTC().Format(time.RFC3339),
		IP:             clientIP(r),
	}

	apps, err := h.loadAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to persist metadata"})
		return
	}
	apps = append(apps, app)
	if err := h.saveAll(apps); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to persist metadata"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "id": app.ID, "cvStoredName": app.CVStoredName, "message": "Application received"})
}

func clientIP(r *http.Request) string {
	if h := r.Header.Get("X-Forwarded-For"); h != "" {
		parts := strings.Split(h, ",")
		return strings.TrimSpace(parts[0])
	}
	if h := r.Header.Get("X-Real-IP"); h != "" {
		return h
	}
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// GET /api/admin/applications?page=1&limit=50&jobId=&search=
func (h *ApplicationHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	apps, err := h.loadAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to load"})
		return
	}
	q := r.URL.Query()
	page := atoiWithDefault(q.Get("page"), 1)
	limit := atoiWithDefault(q.Get("limit"), 50)
	jobId := q.Get("jobId")
	search := strings.ToLower(q.Get("search"))

	var filtered []models.Application
	for _, a := range apps {
		if jobId != "" && a.JobID != jobId {
			continue
		}
		if search != "" {
			if !strings.Contains(strings.ToLower(a.Name), search) &&
				!strings.Contains(strings.ToLower(a.Email), search) &&
				!strings.Contains(strings.ToLower(a.JobID), search) {
				continue
			}
		}
		filtered = append(filtered, a)
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].AppliedAt > filtered[j].AppliedAt })

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	pageItems := filtered[start:end]

	var out []map[string]interface{}
	for _, a := range pageItems {
		out = append(out, map[string]interface{}{
			"id":             a.ID,
			"name":           a.Name,
			"email":          a.Email,
			"phone":          a.Phone,
			"jobId":          a.JobID,
			"appliedAt":      a.AppliedAt,
			"cvOriginalName": a.CVOriginalName,
			"cvStoredName":   a.CVStoredName,
			"downloadUrl":    "/api/admin/applications/cv/" + a.CVStoredName,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": out, "total": total, "page": page, "limit": limit})
}

func atoiWithDefault(s string, d int) int {
	if s == "" {
		return d
	}
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return d
}

// GET /api/admin/applications/{id}
func (h *ApplicationHandler) AdminGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "missing id"})
		return
	}
	apps, err := h.loadAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to load"})
		return
	}
	for _, a := range apps {
		if a.ID == id {
			writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": a})
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "not found"})
}

// DELETE /api/admin/applications/{id}
func (h *ApplicationHandler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "missing id"})
		return
	}
	apps, err := h.loadAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to load"})
		return
	}
	found := false
	var newApps []models.Application
	for _, a := range apps {
		if a.ID == id {
			found = true
			_ = os.Remove(filepath.Join(h.storage, a.CVStoredName))
			continue
		}
		newApps = append(newApps, a)
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "not found"})
		return
	}
	if err := h.saveAll(newApps); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to save"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true})
}

// GET /api/admin/applications/cv/{filename}
func (h *ApplicationHandler) AdminGetCV(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.NotFound(w, r)
		return
	}
	// prevent path traversal: filename must equal its base and not contain '..'
	if strings.Contains(filename, "..") {
		log.Printf("AdminGetCV: rejected filename with '..': %q", filename)
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}
	if filepath.Base(filename) != filename {
		log.Printf("AdminGetCV: rejected filename with path separator: %q (base=%q)", filename, filepath.Base(filename))
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}
	apps, err := h.loadAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to load"})
		return
	}
	var a *models.Application
	for _, ap := range apps {
		if ap.CVStoredName == filename {
			tmp := ap
			a = &tmp
			break
		}
	}
	if a == nil {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(h.storage, filename)
	f, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	fi, _ := f.Stat()

	// Content-Disposition safe filename fallback
	dispName := strings.ReplaceAll(a.CVOriginalName, "\"", "'")
	w.Header().Set("Content-Type", a.CVMimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+dispName+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	http.ServeContent(w, r, a.CVOriginalName, time.Time{}, f)
}
