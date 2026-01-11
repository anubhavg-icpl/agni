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

import "time"

// HealthStatus represents the health status of the system
type HealthStatus struct {
	Status     string                     `json:"status"` // healthy, degraded, unhealthy
	Version    string                     `json:"version"`
	Uptime     string                     `json:"uptime"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
}

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// SystemInfo represents system information
type SystemInfo struct {
	Version            string `json:"version"`
	FirecrackerVersion string `json:"firecracker_version"`
	GoVersion          string `json:"go_version"`
	OS                 string `json:"os"`
	Arch               string `json:"arch"`
	NumCPU             int    `json:"num_cpu"`
	TotalMemory        uint64 `json:"total_memory"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// CreateVMRequest represents a request to create a new VM
type CreateVMRequest struct {
	Name   string   `json:"name"`
	Config VMConfig `json:"config"`
}

// UpdateVMRequest represents a request to update a VM
type UpdateVMRequest struct {
	Name   string   `json:"name,omitempty"`
	Config VMConfig `json:"config,omitempty"`
}

// VMActionResponse represents a response to a VM action
type VMActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	VMID    string `json:"vm_id"`
}

// CreateConfigRequest represents a request to save a config template
type CreateConfigRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Config      VMConfig `json:"config"`
}

// LogEntry represents a log entry from a VM
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source,omitempty"`
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
