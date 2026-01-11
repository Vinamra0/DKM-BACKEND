package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler struct{ BaseURL string }

func NewUploadHandler(baseURL string) *UploadHandler { return &UploadHandler{BaseURL: baseURL} }

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
		http.Error(w, "invalid form", 400)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", 400)
		return
	}
	defer file.Close()

	_ = os.MkdirAll("uploads", 0755)
	name := sanitizeFilename(header.Filename)
	ts := time.Now().UnixNano()
	fname := fmt.Sprintf("%d_%s", ts, name)
	path := filepath.Join("uploads", fname)
	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "save error", 500)
		return
	}
	defer dst.Close()
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "save error", 500)
		return
	}

	url := fmt.Sprintf("%s/uploads/%s", strings.TrimRight(h.BaseURL, "/"), fname)
	writeJSON(w, http.StatusOK, map[string]string{"path": url})
}

func sanitizeFilename(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, s)
	return s
}
