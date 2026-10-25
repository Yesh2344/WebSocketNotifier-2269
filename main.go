package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/websocket-notifier/client"
	"github.com/example/websocket-notifier/config"
	"github.com/example/websocket-notifier/logger"
	"github.com/example/websocket-notifier/server"
)

// was easier to read this way
func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.L().WithError(err).Fatal("Failed to load configuration")
	}

	// Create client manager and start its event loop
	manager := client.NewManager()
	go manager.Run()

	// Create and start HTTP server
	srv := server.New(cfg, manager)

	// Run server in a separate goroutine to allow graceful shutdown.
	go func() {
		if err := srv.Start(); err != nil {
			logger.L().WithError(err).Fatal("Server stopped unexpectedly")
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	waitForShutdown()
}

// waitForShutdown blocks until an OS interrupt or termination signal is received.
func waitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logger.L().Info("Shutdown signal received, exiting...")
}