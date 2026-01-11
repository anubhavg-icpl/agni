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

package handlers

import (
	"net/http"
	"strings"

	"github.com/anubhavg-icpl/agni/internal/auth"
	"github.com/anubhavg-icpl/agni/internal/vm"
	"github.com/anubhavg-icpl/agni/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from Wails app
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	vmManager   *vm.Manager
	authService *auth.Service
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(vmManager *vm.Manager, authService *auth.Service) *WebSocketHandler {
	return &WebSocketHandler{
		vmManager:   vmManager,
		authService: authService,
	}
}

// StreamLogs streams logs for a VM via WebSocket
func (h *WebSocketHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	// Authenticate via query param or header
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := h.authService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get VM ID
	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		http.Error(w, "VM ID required", http.StatusBadRequest)
		return
	}

	// Check VM exists
	_, err = h.vmManager.Get(vmID)
	if err != nil {
		http.Error(w, "VM not found", http.StatusNotFound)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Subscribe to logs
	logStreamer := h.vmManager.GetLogStreamer()
	filter := vm.LogFilter{}

	// Parse filter from query params
	if level := r.URL.Query().Get("level"); level != "" {
		filter.Level = level
	}
	if contains := r.URL.Query().Get("contains"); contains != "" {
		filter.Contains = contains
	}

	sub := logStreamer.Subscribe(vmID, filter)
	defer logStreamer.Unsubscribe(sub.ID)

	// Send logs to client
	done := make(chan struct{})

	// Handle incoming messages (for filter updates, ping, etc.)
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			// Handle ping/pong or filter updates
			if string(msg) == "ping" {
				_ = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			}
		}
	}()

	// Send logs
	for {
		select {
		case log, ok := <-sub.Channel:
			if !ok {
				return
			}
			msg := models.WebSocketMessage{
				Type:    "log",
				Payload: log,
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}
