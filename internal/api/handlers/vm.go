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

	"github.com/anubhavg-icpl/agni/internal/vm"
	"github.com/anubhavg-icpl/agni/pkg/models"
	"github.com/go-chi/chi/v5"
)

// VMHandler handles VM-related requests
type VMHandler struct {
	manager *vm.Manager
}

// NewVMHandler creates a new VMHandler
func NewVMHandler(manager *vm.Manager) *VMHandler {
	return &VMHandler{manager: manager}
}

// List returns all VMs
func (h *VMHandler) List(w http.ResponseWriter, r *http.Request) {
	vms, err := h.manager.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load VMs.")
		return
	}
	respondJSON(w, http.StatusOK, vms)
}

// Create creates a new VM
func (h *VMHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "A VM name is required.")
		return
	}

	req.Config.Name = req.Name
	vm, err := h.manager.Create(req.Config)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create the VM.")
		return
	}

	respondJSON(w, http.StatusCreated, vm)
}

// Get returns a single VM
func (h *VMHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to load the VM.")
		return
	}

	respondJSON(w, http.StatusOK, vm)
}

// Update updates a VM configuration
func (h *VMHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	// Check if VM is running
	if h.manager.IsRunning(id) {
		respondError(w, http.StatusConflict, "The VM is running. Stop it before updating its configuration.")
		return
	}

	var req models.UpdateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update the VM.")
		return
	}

	// Update fields
	if req.Name != "" {
		vm.Name = req.Name
		vm.Config.Name = req.Name
	}
	if req.Config != nil {
		cfg := *req.Config
		// Keep VM name authoritative when only config is patched without a name
		if req.Name != "" {
			cfg.Name = req.Name
		} else if cfg.Name == "" {
			cfg.Name = vm.Name
		} else {
			vm.Name = cfg.Name
		}
		vm.Config = cfg
	}

	if err := h.manager.Update(vm); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save the VM.")
		return
	}

	respondJSON(w, http.StatusOK, vm)
}

// Delete removes a VM
func (h *VMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	if err := h.manager.Delete(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete the VM.")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM deleted.",
		VMID:    id,
	})
}

// Start starts a VM
func (h *VMHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	if err := h.manager.Start(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		if err == models.ErrVMAlreadyRunning {
			respondError(w, http.StatusConflict, "The VM is already running.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to start the VM. Check that the kernel and root drive paths are valid and that the firecracker binary is installed.")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM started.",
		VMID:    id,
	})
}

// Stop force stops a VM
func (h *VMHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	if err := h.manager.Stop(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "The VM is not running.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to stop the VM.")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM stopped.",
		VMID:    id,
	})
}

// Shutdown gracefully shuts down a VM
func (h *VMHandler) Shutdown(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	if err := h.manager.Shutdown(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found.")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "The VM is not running.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Graceful shutdown failed.")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "Graceful shutdown requested.",
		VMID:    id,
	})
}

// Metrics returns metrics for a VM
func (h *VMHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "A VM ID is required.")
		return
	}

	metrics, err := h.manager.GetMetrics(id)
	if err != nil {
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "Metrics are available only while the VM is running.")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to collect VM metrics.")
		return
	}

	respondJSON(w, http.StatusOK, metrics)
}
