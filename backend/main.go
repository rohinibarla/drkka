package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exam42/drkka/backend/handlers"
	"github.com/exam42/drkka/backend/middleware"
	"github.com/exam42/drkka/backend/storage"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./drkka.db"
	}

	// Get static files directory from environment or use default
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "../" // Parent directory contains exam.html, main.js, etc.
	}

	// Initialize SQLite storage
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer store.Close()

	log.Printf("✅ Database initialized: %s", dbPath)

	// Initialize handlers
	submitHandler := handlers.NewSubmitHandler(store)
	staticHandler := handlers.NewStaticFileHandler(staticDir)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/submit", submitHandler.HandleSubmit)

	// Serve static files (HTML, JS, JSON) - this should be last
	mux.Handle("/", staticHandler)

	// Wrap with CORS middleware
	handler := middleware.CORS(mux)

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		// Increase max header size for large submissions
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine for graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", port)
		log.Printf("📊 Health check: http://localhost:%s/health", port)
		log.Printf("📝 Submit endpoint: http://localhost:%s/submit", port)
		log.Printf("📄 Exam page: http://localhost:%s/exam.html", port)
		log.Printf("📄 Review page: http://localhost:%s/review.html", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("🛑 Shutdown signal received: %v", sig)

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("⚠️  Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("❌ Failed to close server: %v", err)
			}
		}

		log.Println("✅ Server stopped gracefully")
	}
}
