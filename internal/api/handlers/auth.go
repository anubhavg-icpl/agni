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

	"github.com/anubhavg-icpl/agni/internal/api/middleware"
	"github.com/anubhavg-icpl/agni/internal/auth"
	"github.com/anubhavg-icpl/agni/pkg/models"
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
		respondError(w, http.StatusBadRequest, "Can't even send valid JSON? We're off to a great start")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password. Both of them. It's not optional")
		return
	}

	user, err := h.authService.Setup(req.Username, req.Password)
	if err != nil {
		if err == models.ErrSetupAlreadyDone {
			respondError(w, http.StatusConflict, "Someone already claimed the throne. You're too slow")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	safeUser := user.SafeUser()
	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "Congratulations, you're the admin now. Try not to break everything",
		"user":    safeUser,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Your request is as malformed as your life choices")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required. This isn't rocket science")
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			respondError(w, http.StatusUnauthorized, "Wrong credentials. Did you forget already? Impressive")
			return
		}
		respondError(w, http.StatusInternalServerError, "Something broke. Probably your fault somehow")
		return
	}

	// Return safe response without password hash
	safeResp := resp.SafeLoginResponse()
	respondJSON(w, http.StatusOK, safeResp)
}

// Status returns the authentication status
func (h *AuthHandler) Status(w http.ResponseWriter, r *http.Request) {
	setupRequired, err := h.authService.IsSetupRequired()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database having an existential crisis")
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
		respondError(w, http.StatusUnauthorized, "Who are you? No seriously, we have no idea")
		return
	}
	safeUser := user.SafeUser()
	respondJSON(w, http.StatusOK, safeUser)
}

// Refresh refreshes the auth token
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get token from header
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		respondError(w, http.StatusBadRequest, "Your authorization header is a disaster")
		return
	}

	token, expiresAt, err := h.authService.RefreshToken(parts[1])
	if err != nil {
		respondError(w, http.StatusUnauthorized, "This token is as expired as your excuses")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"expires_at": expiresAt,
		"message":    "Here's more time. Don't waste it",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT setup, we don't need to do anything server-side
	// The client should discard the token
	// For added security, we could maintain a blacklist of revoked tokens
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Fine, leave. See if we care",
	})
}
