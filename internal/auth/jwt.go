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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/firecracker-microvm/firectl/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// DefaultTokenExpiration is the default JWT token expiration
	DefaultTokenExpiration = 24 * time.Hour
	// SecretKeyLength is the length of generated secret keys
	SecretKeyLength = 32
)

// Claims represents JWT claims
type Claims struct {
	UserID   string          `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

// NewJWTService creates a new JWTService
func NewJWTService(secret string, expiration time.Duration) *JWTService {
	if expiration == 0 {
		expiration = DefaultTokenExpiration
	}
	return &JWTService{
		secret:     []byte(secret),
		expiration: expiration,
		issuer:     "agni",
	}
}

// GenerateSecret generates a random secret key
func GenerateSecret() (string, error) {
	bytes := make([]byte, SecretKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateToken generates a new JWT token for a user
func (j *JWTService) GenerateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.expiration)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, models.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, models.ErrInvalidToken
}

// RefreshToken refreshes a valid token
func (j *JWTService) RefreshToken(tokenString string) (string, time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	// Create a user from claims for token generation
	user := &models.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}

	return j.GenerateToken(user)
}
