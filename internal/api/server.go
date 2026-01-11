// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package api

import (
	"context"
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/anubhavg-icpl/agni/internal/api/handlers"
	"github.com/anubhavg-icpl/agni/internal/api/middleware"
	"github.com/anubhavg-icpl/agni/internal/auth"
	"github.com/anubhavg-icpl/agni/internal/logging"
	"github.com/anubhavg-icpl/agni/internal/storage"
	"github.com/anubhavg-icpl/agni/internal/vm"
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
	Assets     *embed.FS // Embedded frontend assets (optional)
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

	// Serve embedded frontend assets if available
	if s.config.Assets != nil {
		s.serveStaticFiles()
	}
}

// serveStaticFiles serves the embedded frontend assets
func (s *Server) serveStaticFiles() {
	// Get the sub-filesystem for the dist directory
	distFS, err := fs.Sub(*s.config.Assets, "frontend/dist")
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to access embedded frontend assets")
		return
	}

	s.logger.Info().Msg("Static file serving enabled for embedded frontend")

	// Handle all non-API routes
	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Skip API routes
		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Determine the file to serve
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// Try to read the file
		content, err := fs.ReadFile(distFS, filePath)
		if err != nil {
			// Log the error for debugging
			s.logger.Debug().Str("path", filePath).Err(err).Msg("File not found, serving index.html")
			// File not found - serve index.html for SPA routing
			content, err = fs.ReadFile(distFS, "index.html")
			if err != nil {
				s.logger.Error().Err(err).Msg("Failed to read index.html")
				http.NotFound(w, r)
				return
			}
			filePath = "index.html"
		}

		// Set Content-Type using Go's built-in detection
		ext := "." + getFileExt(filePath)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = http.DetectContentType(content)
		}
		s.logger.Debug().Str("path", filePath).Str("ext", ext).Str("contentType", contentType).Msg("Serving file")
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	})
}

// getFileExt returns the file extension without the dot
func getFileExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
		if path[i] == '/' {
			break
		}
	}
	return ""
}

// serveSPAIndex serves the index.html for SPA routing
func (s *Server) serveSPAIndex(w http.ResponseWriter, r *http.Request, distFS fs.FS) {
	content, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// getContentType returns the MIME type based on file extension
func getContentType(path string) string {
	ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
	switch ext {
	case "html":
		return "text/html; charset=utf-8"
	case "css":
		return "text/css; charset=utf-8"
	case "js":
		return "application/javascript; charset=utf-8"
	case "json":
		return "application/json; charset=utf-8"
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "svg":
		return "image/svg+xml"
	case "ico":
		return "image/x-icon"
	case "woff":
		return "font/woff"
	case "woff2":
		return "font/woff2"
	case "ttf":
		return "font/ttf"
	case "webp":
		return "image/webp"
	case "wasm":
		return "application/wasm"
	default:
		return "application/octet-stream"
	}
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
