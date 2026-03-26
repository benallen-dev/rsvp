package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"

	"rsvp/handlers"
	"rsvp/storage"
)

// App holds the application dependencies
type App struct {
	Config *Config
	Store  *storage.Store
	Server *http.Server
}

// Start initializes and starts the application with graceful shutdown
func Start(cfg *Config) error {
	log.Info("Starting application", "port", cfg.Port, "db", cfg.DBPath)

	// Initialize database
	store, closeDB, err := storage.Init(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer closeDB()

	// Create app instance
	app := &App{
		Config: cfg,
		Store:  store,
	}

	// Setup routing
	mux := handlers.NewMux(app.Store)

	// Create HTTP server
	app.Server = &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Info("Listening on " + cfg.Port)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "err", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Info("Shutdown signal received, gracefully stopping server")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "err", err)
		return err
	}

	log.Info("Server stopped")
	return nil
}
