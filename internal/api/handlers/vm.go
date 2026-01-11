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
		respondError(w, http.StatusInternalServerError, "Failed to list VMs. The database is being dramatic")
		return
	}
	respondJSON(w, http.StatusOK, vms)
}

// Create creates a new VM
func (h *VMHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request. Did you even try?")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "A VM needs a name. What were you planning to call it? 'Untitled-47'?")
		return
	}

	req.Config.Name = req.Name
	vm, err := h.manager.Create(req.Config)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "VM creation failed. It's not you, it's... actually, it might be you")
		return
	}

	respondJSON(w, http.StatusCreated, vm)
}

// Get returns a single VM
func (h *VMHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required. We're not mind readers here")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found. Either it never existed or it ghosted you")
			return
		}
		respondError(w, http.StatusInternalServerError, "Something went wrong. Classic")
		return
	}

	respondJSON(w, http.StatusOK, vm)
}

// Update updates a VM configuration
func (h *VMHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "VM ID is required. Did you forget which VM you wanted to update?")
		return
	}

	// Check if VM is running
	if h.manager.IsRunning(id) {
		respondError(w, http.StatusConflict, "Can't update a running VM. Stop it first, genius")
		return
	}

	var req models.UpdateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body. JSON is hard, we know")
		return
	}

	vm, err := h.manager.Get(id)
	if err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "VM not found. Are you sure you created it?")
			return
		}
		respondError(w, http.StatusInternalServerError, "Something went catastrophically wrong")
		return
	}

	// Update fields
	if req.Name != "" {
		vm.Name = req.Name
		vm.Config.Name = req.Name
	}

	// Actually persist the changes (unlike before...)
	if err := h.manager.Update(vm); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save. The database rejected your changes")
		return
	}

	respondJSON(w, http.StatusOK, vm)
}

// Delete removes a VM
func (h *VMHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Which VM? Use your words")
		return
	}

	if err := h.manager.Delete(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "Can't delete what doesn't exist. Philosophy 101")
			return
		}
		respondError(w, http.StatusInternalServerError, "Deletion failed. The VM is fighting back")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM yeeted into the void. Gone. Reduced to atoms",
		VMID:    id,
	})
}

// Start starts a VM
func (h *VMHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Start what exactly? The suspense is killing us")
		return
	}

	if err := h.manager.Start(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "That VM is as real as your productivity today")
			return
		}
		if err == models.ErrVMAlreadyRunning {
			respondError(w, http.StatusConflict, "It's already running. What more do you want from it?")
			return
		}
		respondError(w, http.StatusInternalServerError, "VM refused to start. Can't say we blame it")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM is alive! IT'S ALIVE!",
		VMID:    id,
	})
}

// Stop force stops a VM
func (h *VMHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Stop what? Be specific")
		return
	}

	if err := h.manager.Stop(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "Can't stop a VM that doesn't exist. Bold strategy")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "It's not running. You want us to stop it harder?")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to stop. VM has gone rogue")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM has been forcefully terminated. No mercy",
		VMID:    id,
	})
}

// Shutdown gracefully shuts down a VM
func (h *VMHandler) Shutdown(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Shutdown requires a VM ID. We're not shutting down everything")
		return
	}

	if err := h.manager.Shutdown(id); err != nil {
		if err == models.ErrVMNotFound {
			respondError(w, http.StatusNotFound, "That VM is already in a better place. Or it never existed")
			return
		}
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "Can't shut down something that's already off. Logic, try it")
			return
		}
		respondError(w, http.StatusInternalServerError, "Graceful shutdown failed. Ironic")
		return
	}

	respondJSON(w, http.StatusOK, models.VMActionResponse{
		Success: true,
		Message: "VM is shutting down gracefully. Unlike your code reviews",
		VMID:    id,
	})
}

// Metrics returns metrics for a VM
func (h *VMHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Metrics for which VM? We have several")
		return
	}

	metrics, err := h.manager.GetMetrics(id)
	if err != nil {
		if err == models.ErrVMNotRunning {
			respondError(w, http.StatusConflict, "VM isn't running. No metrics for the lazy")
			return
		}
		respondError(w, http.StatusInternalServerError, "Metrics collection failed. The numbers are shy")
		return
	}

	respondJSON(w, http.StatusOK, metrics)
}
