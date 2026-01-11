// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/firecracker-microvm/firectl/internal/storage"
	"github.com/firecracker-microvm/firectl/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ConfigHandler handles configuration template requests
type ConfigHandler struct {
	store *storage.ConfigStore
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(store *storage.ConfigStore) *ConfigHandler {
	return &ConfigHandler{store: store}
}

// List returns all configuration templates
func (h *ConfigHandler) List(w http.ResponseWriter, r *http.Request) {
	configs, err := h.store.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, configs)
}

// Create creates a new configuration template
func (h *ConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	config := &models.ConfigTemplate{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
	}

	if err := h.store.Create(config); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, config)
}

// Get returns a single configuration template
func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Config ID is required")
		return
	}

	config, err := h.store.Get(id)
	if err != nil {
		if err == models.ErrConfigNotFound {
			respondError(w, http.StatusNotFound, "Configuration not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, config)
}

// Update updates a configuration template
func (h *ConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Config ID is required")
		return
	}

	var req models.CreateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	config, err := h.store.Get(id)
	if err != nil {
		if err == models.ErrConfigNotFound {
			respondError(w, http.StatusNotFound, "Configuration not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update fields
	if req.Name != "" {
		config.Name = req.Name
	}
	if req.Description != "" {
		config.Description = req.Description
	}
	config.Config = req.Config

	if err := h.store.Update(config); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, config)
}

// Delete removes a configuration template
func (h *ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Config ID is required")
		return
	}

	if err := h.store.Delete(id); err != nil {
		if err == models.ErrConfigNotFound {
			respondError(w, http.StatusNotFound, "Configuration not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Configuration deleted",
	})
}
