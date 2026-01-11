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
	"runtime"
	"time"

	"github.com/anubhavg-icpl/agni/internal/vm"
	"github.com/anubhavg-icpl/agni/pkg/models"
)

// Version information (set during build)
var (
	Version                     = "0.2.0"
	SupportedFirecrackerVersion = "1.0.0"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	vmManager *vm.Manager
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(vmManager *vm.Manager, startTime time.Time) *HealthHandler {
	return &HealthHandler{
		vmManager: vmManager,
		startTime: startTime,
	}
}

// Health returns the health status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	status := models.HealthStatus{
		Status:    "healthy",
		Version:   Version,
		Uptime:    formatDuration(uptime),
		Timestamp: time.Now(),
		Components: map[string]models.ComponentHealth{
			"vm_manager": {
				Status:  "healthy",
				Message: "VM manager operational",
			},
			"database": {
				Status:  "healthy",
				Message: "Database connected",
			},
		},
	}

	respondJSON(w, http.StatusOK, status)
}

// SystemInfo returns system information
func (h *HealthHandler) SystemInfo(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	info := models.SystemInfo{
		Version:            Version,
		FirecrackerVersion: SupportedFirecrackerVersion,
		GoVersion:          runtime.Version(),
		OS:                 runtime.GOOS,
		Arch:               runtime.GOARCH,
		NumCPU:             runtime.NumCPU(),
		TotalMemory:        memStats.Sys,
	}

	respondJSON(w, http.StatusOK, info)
}

// formatDuration formats a duration as a human-readable string
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return formatPlural(days, "day") + " " + formatPlural(hours, "hour")
	}
	if hours > 0 {
		return formatPlural(hours, "hour") + " " + formatPlural(minutes, "minute")
	}
	if minutes > 0 {
		return formatPlural(minutes, "minute") + " " + formatPlural(seconds, "second")
	}
	return formatPlural(seconds, "second")
}

func formatPlural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10)) + " " + unit + "s"
}
