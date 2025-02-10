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

	"github.com/Saurab-Negi/Go-CRUD/internal/config"
	student "github.com/Saurab-Negi/Go-CRUD/internal/http/handlers/students"
	"github.com/Saurab-Negi/Go-CRUD/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup
	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version","1.0.0"))

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	// Setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Fails to shutdown the server", slog.String("error", err.Error()))
	}

	// alternative way to write above code
	// if err := server.Shutdown(ctx); err != nil {
	// 	slog.Error("Fails to shutdown the server", slog.String("error", err.Error()))
	// }

	slog.Info("Server shutdown gracefully")
}