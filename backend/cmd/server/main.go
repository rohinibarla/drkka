package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize SQLite storage
	store, err := storage.NewSQLiteStorage(cfg.DB.Path)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer store.Close()

	log.Printf("✅ Database initialized: %s", cfg.DB.Path)

	// Initialize handlers
	submitHandler := handlers.NewSubmitHandler(store)
	submissionsHandler := handlers.NewSubmissionsHandler(store)
	staticHandler := handlers.NewStaticFileHandler(cfg.Static.Dir)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/submit", submitHandler.HandleSubmit)
	mux.HandleFunc("/submissions", submissionsHandler.HandleListSubmissions)

	// Serve static files (HTML, JS, JSON) - this should be last
	mux.Handle("/", staticHandler)

	// Wrap with CORS middleware
	handler := middleware.CORS(&cfg.CORS)(mux)

	// Configure server
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        handler,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// Start server in goroutine for graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", cfg.Server.Port)
		log.Printf("📊 Health check: http://localhost:%s/health", cfg.Server.Port)
		log.Printf("📝 Submit endpoint: http://localhost:%s/submit", cfg.Server.Port)
		log.Printf("📋 Submissions list: http://localhost:%s/submissions", cfg.Server.Port)
		log.Printf("📄 Exam page: http://localhost:%s/exam.html", cfg.Server.Port)
		log.Printf("📄 Review page: http://localhost:%s/review.html", cfg.Server.Port)
		log.Printf("📄 Submissions page: http://localhost:%s/submissions.html", cfg.Server.Port)
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
