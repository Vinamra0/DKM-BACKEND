package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadHandler struct {
	BaseURL   string
	UploadDir string
	MaxSize   int64
}

func NewUploadHandler(baseURL, uploadDir string, maxSize int64) *UploadHandler {
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = "uploads"
	}
	if maxSize <= 0 {
		maxSize = 20 << 20 // 20MB default
	}
	return &UploadHandler{BaseURL: baseURL, UploadDir: uploadDir, MaxSize: maxSize}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.MaxSize+1024)
	if err := r.ParseMultipartForm(h.MaxSize); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}
	file, header, err := firstFormFile(r, "file", "image", "upload")
	if err != nil {
		http.Error(w, "file required", 400)
		return
	}
	defer file.Close()

	_ = os.MkdirAll(h.UploadDir, 0755)
	name := sanitizeFilename(header.Filename)
	ts := time.Now().UnixNano()
	fname := fmt.Sprintf("%d_%s", ts, name)
	path := filepath.Join(h.UploadDir, fname)
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "save error", 500)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "save error", 500)
		return
	}

	base := strings.TrimSpace(h.BaseURL)
	if base == "" {
		scheme := r.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = "http"
		}
		base = scheme + "://" + r.Host
	}
	url := fmt.Sprintf("%s/uploads/%s", strings.TrimRight(base, "/"), fname)
	writeJSON(w, http.StatusOK, map[string]string{"path": url})
}

func firstFormFile(r *http.Request, keys ...string) (multipart.File, *multipart.FileHeader, error) {
	for _, key := range keys {
		f, hdr, err := r.FormFile(key)
		if err == nil {
			return f, hdr, nil
		}
	}
	return nil, nil, http.ErrMissingFile
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
