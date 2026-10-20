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
	"sync"
	"time"

	"github.com/anubhavg-icpl/agni/internal/auth"
	"github.com/anubhavg-icpl/agni/internal/vm"
	"github.com/anubhavg-icpl/agni/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	zlog "github.com/rs/zerolog/log"
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

	// gorilla/websocket connections are not safe for concurrent writers
	var writeMu sync.Mutex

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
				writeMu.Lock()
				_ = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				writeMu.Unlock()
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
			writeMu.Lock()
			err := conn.WriteJSON(msg)
			writeMu.Unlock()
			if err != nil {
				return
			}
		case <-done:
			return
		}
	}
}

// Console streams the VM serial console bidirectionally: guest output is
// forwarded as text frames, client frames are written to the guest stdin.
func (h *WebSocketHandler) Console(w http.ResponseWriter, r *http.Request) {
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
	if _, err := h.authService.ValidateToken(token); err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		http.Error(w, "VM ID required", http.StatusBadRequest)
		return
	}
	if _, err := h.vmManager.Get(vmID); err != nil {
		http.Error(w, "VM not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Serial output goes to the client one line per text frame.
	var writeMu sync.Mutex
	writeFrame := func(text string) {
		writeMu.Lock()
		defer writeMu.Unlock()
		_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_ = conn.WriteMessage(websocket.TextMessage, []byte(text))
	}
	writeLine := func(text string) { writeFrame(text + "\n") }

	// Reader: terminal input -> guest stdin. Closing signals the writer loop.
	// Input problems never tear down the console: the user keeps the session
	// and sees an inline notice instead of a reconnect loop.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				h.logClose(vmID, "read error", err)
				return
			}
			if len(data) > 0 {
				if err := h.vmManager.WriteConsole(vmID, data); err != nil {
					writeLine("[agni] " + err.Error())
				}
			}
		}
	}()

	sub := h.vmManager.GetLogStreamer().Subscribe(vmID, vm.LogFilter{})
	defer h.vmManager.GetLogStreamer().Unsubscribe(sub.ID)

	for {
		select {
		case entry, ok := <-sub.Channel:
			if !ok {
				h.logClose(vmID, "log channel closed", nil)
				return
			}
			writeLine(entry.Message)
		case <-done:
			h.logClose(vmID, "client input loop ended", nil)
			return
		case <-r.Context().Done():
			h.logClose(vmID, "request context done", r.Context().Err())
			return
		}
	}
}

// logClose records why a console session ended, making reconnect loops
// diagnosable in the server log.
func (h *WebSocketHandler) logClose(vmID, reason string, err error) {
	msg := reason
	if err != nil {
		msg += ": " + err.Error()
	}
	zlog.Warn().Str("vm_id", vmID).Str("component", "console").Msg("Console session ended: " + msg)
}
