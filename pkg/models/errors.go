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

package models

import "errors"

// Configuration errors
var (
	ErrInvalidNicConfig                  = errors.New("NIC config must be specified as DEVICE/MAC")
	ErrInvalidDriveSpecificationNoSuffix = errors.New("drive specification must have :ro or :rw suffix")
	ErrInvalidDriveSpecificationNoPath   = errors.New("drive specification is missing path")
	ErrUnableToParseVsockDevices         = errors.New("unable to parse vsock devices")
	ErrUnableToParseVsockCID             = errors.New("unable to parse vsock CID as uint")
	ErrConflictingLogOpts                = errors.New("vmm-log-fifo and firecracker-log are mutually exclusive")
	ErrUnableToCreateFifoLogFile         = errors.New("failed to create fifo log file")
	ErrInvalidMetadata                   = errors.New("metadata is not valid JSON")
)

// VM errors
var (
	ErrVMNotFound       = errors.New("VM not found")
	ErrVMAlreadyExists  = errors.New("VM already exists")
	ErrVMNotRunning     = errors.New("VM is not running")
	ErrVMAlreadyRunning = errors.New("VM is already running")
	ErrVMStartFailed    = errors.New("failed to start VM")
	ErrVMStopFailed     = errors.New("failed to stop VM")
)

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrSetupRequired      = errors.New("initial setup required")
	ErrSetupAlreadyDone   = errors.New("setup already completed")
)

// Storage errors
var (
	ErrDatabaseNotInitialized = errors.New("database not initialized")
	ErrConfigNotFound         = errors.New("configuration not found")
)

// APIError represents an API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code int, message string, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
