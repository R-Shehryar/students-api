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
	student "github.com/R-Shehryar/students-api/internal/http/handlers/students"
	"github.com/R-Shehryar/students-api/internal/storage/sqlite"
)

func main() {
	// Load configuration
	cfg := config.MustLoadConfig()
	// database setup
	storage, er := sqlite.New(cfg)
	if er != nil {
		log.Fatal("Failed to initialize storage:", er)
	}
	slog.Info("Database initialized successfully.", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetStudentByID(storage))
	router.HandleFunc("GET /api/students", student.GetAllStudents(storage))
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
