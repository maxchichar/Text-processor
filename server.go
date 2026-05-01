package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"text-processor/api"
)

func main() {
	// Create API server
	server, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create HTTP multiplexer
	mux := http.NewServeMux()

	// Register API routes
	server.RegisterRoutes(mux)

	// Serve static files from frontend/dist if it exists
	frontendDist := "frontend/dist"
	if _, err := os.Stat(frontendDist); err == nil {
		fs := http.FileServer(http.Dir(frontendDist))
		mux.Handle("/", fs)
		log.Println("Serving frontend from", frontendDist)
	} else {
		// Fallback: serve a simple HTML page
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, `<!DOCTYPE html>
<html>
<head>
	<title>Text Processor</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
		h1 { color: #333; }
		.info { background: #f4f4f4; padding: 15px; border-radius: 5px; margin: 20px 0; }
		.api-endpoint { background: #e8f4f8; padding: 10px; margin: 5px 0; border-radius: 3px; }
		code { background: #f0f0f0; padding: 2px 6px; border-radius: 3px; }
	</style>
</head>
<body>
	<h1>Text Processor API</h1>
	<div class="info">
		<p>The Text Processor API is running! Build the frontend to access the web interface.</p>
	</div>
	<h2>API Endpoints</h2>
	<div class="api-endpoint"><code>POST /api/upload</code> - Upload a .txt file</div>
	<div class="api-endpoint"><code>POST /api/process/{fileId}</code> - Process uploaded file</div>
	<div class="api-endpoint"><code>GET /api/preview/{fileId}</code> - Preview processed text</div>
	<div class="api-endpoint"><code>GET /api/download/{fileId}</code> - Download processed file</div>
	<div class="api-endpoint"><code>GET /api/history</code> - Get processing history</div>
	<div class="api-endpoint"><code>DELETE /api/files/{fileId}</code> - Delete files</div>
</body>
</html>`)
		})
	}

	// Configure HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}