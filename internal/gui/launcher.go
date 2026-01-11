// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// Copyright 2024 Anubhav Gain. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

// Package gui provides the GUI launcher functionality for firectl.
package gui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/firecracker-microvm/firectl/internal/api"
	"github.com/firecracker-microvm/firectl/internal/auth"
	"github.com/firecracker-microvm/firectl/internal/logging"
	"github.com/firecracker-microvm/firectl/internal/storage"
	"github.com/firecracker-microvm/firectl/internal/vm"
)

// Config holds the GUI launcher configuration
type Config struct {
	Port    string
	DataDir string
	Logger  *logging.Logger
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Port:    "8080",
		DataDir: GetDataDir(),
		Logger:  nil,
	}
}

// Launcher handles starting and stopping the GUI server
type Launcher struct {
	config    Config
	logger    *logging.Logger
	store     *storage.Store
	vmManager *vm.Manager
	apiServer *api.Server
}

// NewLauncher creates a new GUI launcher
func NewLauncher(cfg Config) *Launcher {
	if cfg.Logger == nil {
		cfg.Logger = logging.NewLogger(logging.Config{
			Level:  "info",
			Pretty: true,
		})
	}

	return &Launcher{
		config: cfg,
		logger: cfg.Logger,
	}
}

// Start initializes and starts the GUI server
func (l *Launcher) Start() error {
	l.logger.Info().Msg("Starting Firectl GUI")

	// Ensure data directory exists
	if err := os.MkdirAll(l.config.DataDir, 0750); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize storage
	dbPath := filepath.Join(l.config.DataDir, "firectl.db")
	store, err := storage.NewStore(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	l.store = store

	// Generate or load JWT secret
	jwtSecret := l.getOrCreateJWTSecret()

	// Initialize VM manager
	l.vmManager = vm.NewManager(store)

	// Initialize API server
	l.apiServer = api.NewServer(api.ServerConfig{
		Address:    ":" + l.config.Port,
		JWTSecret:  jwtSecret,
		VMManager:  l.vmManager,
		Store:      store,
		EnableCORS: true,
		RateLimit:  100,
	})

	// Start API server in background
	go func() {
		if err := l.apiServer.Start(":" + l.config.Port); err != nil {
			l.logger.Error().Err(err).Msg("API server error")
		}
	}()

	l.logger.Info().Str("port", l.config.Port).Msg("API server started")

	// Check if setup is required
	setupRequired, _ := l.apiServer.SetupRequired()
	if setupRequired {
		l.logger.Info().Msg("Initial setup required. Create admin user at POST /api/auth/setup")
	}

	return nil
}

// Stop gracefully shuts down the GUI server
func (l *Launcher) Stop() {
	l.logger.Info().Msg("Shutting down...")

	if l.vmManager != nil {
		l.vmManager.StopAll()
	}

	if l.apiServer != nil {
		l.apiServer.Shutdown(context.Background())
	}

	if l.store != nil {
		l.store.Close()
	}

	l.logger.Info().Msg("Goodbye!")
}

// WaitForShutdown blocks until an interrupt signal is received
func (l *Launcher) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	l.logger.Info().Msg("Received shutdown signal")
}

// PrintBanner prints the startup banner
func (l *Launcher) PrintBanner() {
	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("  Firectl GUI - API Server Mode")
	fmt.Println("==============================================")
	fmt.Printf("  API:    http://localhost:%s/api\n", l.config.Port)
	fmt.Printf("  Health: http://localhost:%s/api/health\n", l.config.Port)
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
}

// getOrCreateJWTSecret loads or creates a JWT secret
func (l *Launcher) getOrCreateJWTSecret() string {
	secretPath := filepath.Join(l.config.DataDir, ".jwt_secret")

	// Try to read existing secret
	if data, err := os.ReadFile(secretPath); err == nil && len(data) > 0 {
		return string(data)
	}

	// Generate new secret
	secret, err := auth.GenerateSecret()
	if err != nil {
		l.logger.Fatal().Err(err).Msg("Failed to generate JWT secret")
	}

	// Save secret
	if err := os.WriteFile(secretPath, []byte(secret), 0600); err != nil {
		l.logger.Warn().Err(err).Msg("Failed to save JWT secret, using temporary secret")
	}

	return secret
}

// GetDataDir returns the data directory for firectl
func GetDataDir() string {
	// Try XDG_DATA_HOME first
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, "firectl")
	}

	// Fall back to ~/.local/share/firectl
	home, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/firectl"
	}

	return filepath.Join(home, ".local", "share", "firectl")
}

// PrintHelp prints help for GUI mode
func PrintHelp() {
	fmt.Print(`
Firectl GUI Mode

Usage:
  firectl --gui [options]
  firectl -g [options]

Options:
  --port PORT      API server port (default: 8080)
  --data-dir DIR   Data directory (default: ~/.local/share/firectl)

The GUI provides a web-based interface for managing Firecracker VMs.
Access the interface at http://localhost:8080 after starting.
`)
}
