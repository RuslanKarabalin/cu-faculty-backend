package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
		http.ServeFile(w, r, "static/index.html")
	}))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Server starting on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Server starting failed", slog.Any("error", err))
		os.Exit(1)
	}
}
