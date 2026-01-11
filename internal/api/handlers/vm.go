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

	"github.com/firecracker-microvm/firectl/internal/vm"
	"github.com/firecracker-microvm/firectl/pkg/models"
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
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, vms)
}

// Create creates a new VM
func (h *VMHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	req.Config.Name = req.Name
	vm, err := h.manager.Create(req.Config)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, vm)
}

// Get returns a single VM
func (h *VMHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, vm)
}

// Update updates a VM configuration
func (h *VMHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	// Check if VM is running
	if h.manager.IsRunning(id) {
		respondError(w, http.StatusConflict, "Cannot update running VM")
		return
	}

	var req models.UpdateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update fields
	if req.Name != "" {
		vm.Name = req.Name
		vm.Config.Name = req.Name
	}

	// For now, we don't support updating the full config of a created VM
	// This would need more careful handling

	respondJSON(w, http.StatusOK, vm)
}

// Delete removes a VM
func (h *VMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	if err := h.manager.Delete(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM deleted",
		VMID:    id,
	})
}

// Start starts a VM
func (h *VMHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	if err := h.manager.Start(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		if err == models.ErrVMAlreadyRunning {
			respondError(w, http.StatusConflict, "VM is already running")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM started",
		VMID:    id,
	})
}

// Stop force stops a VM
func (h *VMHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	if err := h.manager.Stop(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "VM is not running")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM stopped",
		VMID:    id,
	})
}

// Shutdown gracefully shuts down a VM
func (h *VMHandler) Shutdown(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	if err := h.manager.Shutdown(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "VM is not running")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM shutdown initiated",
		VMID:    id,
	})
}

// Metrics returns metrics for a VM
func (h *VMHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required")
		return
	}

	metrics, err := h.manager.GetMetrics(id)
	if err != nil {
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "VM is not running")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, metrics)
}
