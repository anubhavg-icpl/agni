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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/firecracker-microvm/firectl/internal/api/middleware"
	"github.com/firecracker-microvm/firectl/internal/auth"
	"github.com/firecracker-microvm/firectl/pkg/models"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Setup handles first-time setup
func (h *AuthHandler) Setup(w http.ResponseWriter, r *http.Request) {
	var req models.SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, err := h.authService.Setup(req.Username, req.Password)
	if err != nil {
		if err == models.ErrSetupAlreadyDone {
			respondError(w, http.StatusConflict, "Setup already completed")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "Admin user created",
		"user":    user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			respondError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// Status returns the authentication status
func (h *AuthHandler) Status(w http.ResponseWriter, r *http.Request) {
	setupRequired, err := h.authService.IsSetupRequired()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"setup_required": setupRequired,
	})
}

// Me returns the current user
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// Refresh refreshes the auth token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get token from header
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		respondError(w, http.StatusBadRequest, "Invalid authorization header")
		return
	}

	token, expiresAt, err := h.authService.RefreshToken(parts[1])
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"expires_at": expiresAt,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT setup, we don't need to do anything server-side
	// The client should discard the token
	// For added security, we could maintain a blacklist of revoked tokens
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Logged out",
	})
}
