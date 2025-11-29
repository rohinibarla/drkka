package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticFileHandler serves static files from the specified directory
type StaticFileHandler struct {
	staticDir string
}

// NewStaticFileHandler creates a new static file handler
func NewStaticFileHandler(staticDir string) *StaticFileHandler {
	// Get absolute path
	absPath, err := filepath.Abs(staticDir)
	if err != nil {
		log.Printf("Warning: Could not resolve absolute path for %s: %v", staticDir, err)
		absPath = staticDir
	}

	log.Printf("📁 Static files directory: %s", absPath)

	return &StaticFileHandler{
		staticDir: absPath,
	}
}

// ServeHTTP handles static file requests
func (h *StaticFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Security: prevent directory traversal attacks
	cleanPath := filepath.Clean(r.URL.Path)
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Determine file path
	var filePath string
	if r.URL.Path == "/" {
		// Serve exam.html as the default page
		filePath = filepath.Join(h.staticDir, "exam.html")
	} else {
		// Serve the requested file
		filePath = filepath.Join(h.staticDir, cleanPath)
	}

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Don't serve directories
	if fileInfo.IsDir() {
		http.Error(w, "Directory listing not allowed", http.StatusForbidden)
		return
	}

	// Set appropriate content type based on file extension
	contentType := getContentType(filePath)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
	log.Printf("📄 Served: %s -> %s", r.URL.Path, filePath)
}

// getContentType returns the appropriate content type for a file
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	contentTypes := map[string]string{
		".html": "text/html; charset=utf-8",
		".js":   "application/javascript; charset=utf-8",
		".json": "application/json; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}

	return ""
}
