// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/firecracker-microvm/firectl/internal/api/handlers"
	"github.com/firecracker-microvm/firectl/internal/api/middleware"
	"github.com/firecracker-microvm/firectl/internal/auth"
	"github.com/firecracker-microvm/firectl/internal/logging"
	"github.com/firecracker-microvm/firectl/internal/storage"
	"github.com/firecracker-microvm/firectl/internal/vm"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// ServerConfig holds configuration for the API server
type ServerConfig struct {
	Address    string
	JWTSecret  string
	VMManager  *vm.Manager
	Store      *storage.Store
	EnableCORS bool
	RateLimit  int
}

// Server represents the HTTP API server
type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	config      ServerConfig
	vmManager   *vm.Manager
	authService *auth.Service
	logger      *logging.Logger
	startTime   time.Time
}

// NewServer creates a new API server
func NewServer(cfg ServerConfig) *Server {
	userStore := storage.NewUserStore(cfg.Store)
	authService := auth.NewService(userStore, cfg.JWTSecret)

	s := &Server{
		router:      chi.NewRouter(),
		config:      cfg,
		vmManager:   cfg.VMManager,
		authService: authService,
		logger:      logging.GetLogger().WithComponent("api-server"),
		startTime:   time.Now(),
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Request ID
	s.router.Use(chimiddleware.RequestID)

	// Real IP
	s.router.Use(chimiddleware.RealIP)

	// Logging
	s.router.Use(middleware.RequestLogger(s.logger))

	// Recovery
	s.router.Use(middleware.Recovery(s.logger))

	// CORS
	if s.config.EnableCORS {
		s.router.Use(middleware.CORS())
	}

	// Rate limiting
	if s.config.RateLimit > 0 {
		s.router.Use(middleware.RateLimit(s.config.RateLimit))
	}

	// Timeout
	s.router.Use(chimiddleware.Timeout(60 * time.Second))
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// Health check (no auth required)
	healthHandler := handlers.NewHealthHandler(s.vmManager, s.startTime)
	s.router.Get("/api/health", healthHandler.Health)
	s.router.Get("/api/system/info", healthHandler.SystemInfo)

	// Auth routes
	authHandler := handlers.NewAuthHandler(s.authService)
	s.router.Post("/api/auth/setup", authHandler.Setup)
	s.router.Post("/api/auth/login", authHandler.Login)
	s.router.Get("/api/auth/status", authHandler.Status)

	// Protected routes
	s.router.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(s.authService))

		// Auth
		r.Get("/api/auth/me", authHandler.Me)
		r.Post("/api/auth/refresh", authHandler.Refresh)
		r.Post("/api/auth/logout", authHandler.Logout)

		// VMs
		vmHandler := handlers.NewVMHandler(s.vmManager)
		r.Get("/api/vms", vmHandler.List)
		r.Post("/api/vms", vmHandler.Create)
		r.Get("/api/vms/{id}", vmHandler.Get)
		r.Put("/api/vms/{id}", vmHandler.Update)
		r.Delete("/api/vms/{id}", vmHandler.Delete)
		r.Post("/api/vms/{id}/start", vmHandler.Start)
		r.Post("/api/vms/{id}/stop", vmHandler.Stop)
		r.Post("/api/vms/{id}/shutdown", vmHandler.Shutdown)
		r.Get("/api/vms/{id}/metrics", vmHandler.Metrics)

		// Configs
		configStore := storage.NewConfigStore(s.config.Store)
		configHandler := handlers.NewConfigHandler(configStore)
		r.Get("/api/configs", configHandler.List)
		r.Post("/api/configs", configHandler.Create)
		r.Get("/api/configs/{id}", configHandler.Get)
		r.Put("/api/configs/{id}", configHandler.Update)
		r.Delete("/api/configs/{id}", configHandler.Delete)
	})

	// WebSocket routes (with auth check in handler)
	wsHandler := handlers.NewWebSocketHandler(s.vmManager, s.authService)
	s.router.Get("/api/vms/{id}/logs", wsHandler.StreamLogs)
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	if addr == "" {
		addr = s.config.Address
	}
	if addr == "" {
		addr = ":8080"
	}

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info().Str("address", addr).Msg("Starting API server")
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down API server")
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// Router returns the chi router (for Wails integration)
func (s *Server) Router() *chi.Mux {
	return s.router
}

// GetAddress returns the server address
func (s *Server) GetAddress() string {
	if s.httpServer != nil {
		return s.httpServer.Addr
	}
	return s.config.Address
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// GetVMManager returns the VM manager
func (s *Server) GetVMManager() *vm.Manager {
	return s.vmManager
}

// GetAuthService returns the auth service
func (s *Server) GetAuthService() *auth.Service {
	return s.authService
}

// SetupRequired checks if initial setup is needed
func (s *Server) SetupRequired() (bool, error) {
	return s.authService.IsSetupRequired()
}
