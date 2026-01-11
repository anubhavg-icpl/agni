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

package auth

import (
	"time"

	"github.com/firecracker-microvm/firectl/internal/storage"
	"github.com/firecracker-microvm/firectl/pkg/models"
	"github.com/google/uuid"
)

// Service provides authentication operations
type Service struct {
	userStore  *storage.UserStore
	jwtService *JWTService
}

// NewService creates a new auth Service
func NewService(userStore *storage.UserStore, jwtSecret string) *Service {
	return &Service{
		userStore:  userStore,
		jwtService: NewJWTService(jwtSecret, DefaultTokenExpiration),
	}
}

// Login authenticates a user and returns a token
func (s *Service) Login(username, password string) (*models.LoginResponse, error) {
	user, err := s.userStore.GetByUsername(username)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	// Update last login (ignore error, non-critical)
	_ = s.userStore.UpdateLastLogin(user.ID)

	token, expiresAt, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

// Setup creates the initial admin user
func (s *Service) Setup(username, password string) (*models.User, error) {
	setupRequired, err := s.userStore.IsSetupRequired()
	if err != nil {
		return nil, err
	}
	if !setupRequired {
		return nil, models.ErrSetupAlreadyDone
	}

	if !ValidatePasswordStrength(password) {
		return nil, models.NewAPIError(400, "Password must be at least 8 characters", "")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: passwordHash,
		Role:         models.UserRoleAdmin,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userStore.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ValidateToken validates a token and returns the claims
func (s *Service) ValidateToken(token string) (*Claims, error) {
	return s.jwtService.ValidateToken(token)
}

// RefreshToken refreshes a valid token
func (s *Service) RefreshToken(token string) (string, time.Time, error) {
	return s.jwtService.RefreshToken(token)
}

// IsSetupRequired checks if initial setup is needed
func (s *Service) IsSetupRequired() (bool, error) {
	return s.userStore.IsSetupRequired()
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id string) (*models.User, error) {
	return s.userStore.Get(id)
}

// CreateUser creates a new user (admin only)
func (s *Service) CreateUser(username, password string, role models.UserRole) (*models.User, error) {
	if !ValidatePasswordStrength(password) {
		return nil, models.NewAPIError(400, "Password must be at least 8 characters", "")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userStore.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(userID, currentPassword, newPassword string) error {
	user, err := s.userStore.Get(userID)
	if err != nil {
		return err
	}

	if !CheckPassword(currentPassword, user.PasswordHash) {
		return models.ErrInvalidCredentials
	}

	if !ValidatePasswordStrength(newPassword) {
		return models.NewAPIError(400, "Password must be at least 8 characters", "")
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	return s.userStore.Update(user)
}
