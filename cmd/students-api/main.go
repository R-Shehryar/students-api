package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/R-Shehryar/students-api/internal/config"
)

func main() {
	// Load configuration
	cfg := config.MustLoadConfig()
	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Students API!"))
	})
	// setup server
	server := &http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}
	log.Println("Server is running on", cfg.HttpServer.Address)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-done
	slog.Info("Server is shutting down...")
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(context)
	if err != nil {
		slog.Error("Failed to gracefully shutdown the server:", slog.Any("error", err))
	}
	slog.Info("Server has been shut down.")
}
