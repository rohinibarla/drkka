package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware adds CORS headers to responses
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment variable or use default
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:3000,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:8080"
		}

		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if origin != "" && isOriginAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if allowedOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin, allowedOrigins string) bool {
	if allowedOrigins == "*" {
		return true
	}

	origins := strings.Split(allowedOrigins, ",")
	for _, allowed := range origins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}

	return false
}
