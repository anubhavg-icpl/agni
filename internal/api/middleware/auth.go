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

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/anubhavg-icpl/agni/internal/auth"
	"github.com/anubhavg-icpl/agni/pkg/models"
)

// Context keys
type contextKey string

const (
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// JWTAuth returns a middleware that validates JWT tokens
func JWTAuth(authService *auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "No authorization header. Were you raised in a barn?")
				return
			}

			// Parse Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				respondError(w, http.StatusUnauthorized, "It's 'Bearer <token>'. Not that hard")
				return
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Your token is as valid as a three-dollar bill")
				return
			}

			// Get user
			user, err := authService.GetUser(claims.UserID)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Token is valid but you don't exist. Spooky")
				return
			}

			// Add claims and user to context
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			ctx = context.WithValue(ctx, UserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims extracts claims from context
func GetClaims(ctx context.Context) *auth.Claims {
	if claims, ok := ctx.Value(ClaimsContextKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}

// GetUser extracts user from context
func GetUser(ctx context.Context) *models.User {
	if user, ok := ctx.Value(UserContextKey).(*models.User); ok {
		return user
	}
	return nil
}

// RequireAdmin returns a middleware that requires admin role
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil || user.Role != models.UserRoleAdmin {
			respondError(w, http.StatusForbidden, "Nice try, peasant. Admins only")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}
