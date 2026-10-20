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
	"os"
	"path/filepath"
	"syscall"

	"github.com/anubhavg-icpl/agni/internal/images"
	"github.com/anubhavg-icpl/agni/internal/templates"
)

// StorageHandler reports how much disk the image cache uses.
type StorageHandler struct {
	images  *images.Manager
	catalog func() []templates.Template
}

// NewStorageHandler builds a StorageHandler over the images directory.
func NewStorageHandler(img *images.Manager, catalog func() []templates.Template) *StorageHandler {
	return &StorageHandler{images: img, catalog: catalog}
}

// Usage walks the images directory and reports per-template disk usage.
func (h *StorageHandler) Usage(w http.ResponseWriter, r *http.Request) {
	type templateUsage struct {
		ID      string `json:"id"`
		Ready   bool   `json:"ready"`
		Bytes   int64  `json:"bytes"`
		IsCache bool   `json:"is_cache"`
	}
	root := h.images.Dir()
	entries, err := filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	usage := make([]templateUsage, 0, len(entries))
	var total, cache int64
	ready := map[string]bool{}
	if h.catalog != nil {
		for _, t := range h.catalog() {
			ready[t.ID] = h.images.Ready(t)
		}
	}
	for _, dir := range entries {
		name := filepath.Base(dir)
		size := dirSize(dir)
		isCache := name == "_cache"
		usage = append(usage, templateUsage{
			ID: name, Ready: ready[name], Bytes: size, IsCache: isCache,
		})
		total += size
		if isCache {
			cache = size
		}
	}
	var free, capacity uint64
	var sysSt syscall.Statfs_t
	if err := syscall.Statfs(root, &sysSt); err == nil {
		free = sysSt.Bavail * uint64(sysSt.Bsize)
		capacity = sysSt.Blocks * uint64(sysSt.Bsize)
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data_directory":   root,
		"total_bytes":      total,
		"cache_bytes":      cache,
		"disk_free_bytes":  free,
		"disk_total_bytes": capacity,
		"templates":        usage,
	})
}

func dirSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info.Mode().IsRegular() {
			size += info.Size()
		}
		return nil
	})
	return size
}
