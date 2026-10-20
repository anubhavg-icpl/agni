package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anubhavg-icpl/agni/internal/images"
	"github.com/anubhavg-icpl/agni/internal/templates"
	"github.com/anubhavg-icpl/agni/internal/vm"
	"github.com/anubhavg-icpl/agni/pkg/models"
	"github.com/go-chi/chi/v5"
)

// TemplateHandler serves the deployable template catalog and one-click deploys.
type TemplateHandler struct {
	images  *images.Manager
	manager *vm.Manager
}

// NewTemplateHandler creates a TemplateHandler.
func NewTemplateHandler(imgs *images.Manager, manager *vm.Manager) *TemplateHandler {
	return &TemplateHandler{images: imgs, manager: manager}
}

type templateListItem struct {
	templates.Template
	Artifacts []images.ArtifactState `json:"artifacts"`
	Ready     bool                   `json:"ready"`
	LatestJob *images.Job            `json:"latest_job,omitempty"`
}

type runtimeState struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path,omitempty"`
	Version   string `json:"version"`
}

type templatesResponse struct {
	Templates    []templateListItem `json:"templates"`
	Runtime      runtimeState       `json:"runtime"`
	KVMAvailable bool               `json:"kvm_available"`
	ImagesDir    string             `json:"images_dir"`
}

type prepareResponse struct {
	Job *images.Job `json:"job"`
}

type deployFromTemplateRequest struct {
	TemplateID string `json:"template_id"`
	Name       string `json:"name"`
	CPUs       int64  `json:"cpus"`
	MemoryMB   int64  `json:"memory_mb"`
	AutoStart  bool   `json:"auto_start"`
}

type deployResponse struct {
	VM         *models.VM `json:"vm"`
	Started    bool       `json:"started"`
	StartError string     `json:"start_error,omitempty"`
}

// List returns the template catalog with local disk state.
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	catalog := templates.Catalog()
	items := make([]templateListItem, 0, len(catalog))
	for _, t := range catalog {
		item := templateListItem{
			Template:  t,
			Artifacts: h.images.State(t),
			Ready:     h.images.Ready(t),
			LatestJob: h.images.LatestJobFor(t.ID),
		}
		items = append(items, item)
	}

	runtime := runtimeState{Installed: false, Version: "v1.12.1"}
	if path, ok := h.images.RuntimePath(); ok {
		runtime.Installed = true
		runtime.Path = path
	}

	respondJSON(w, http.StatusOK, templatesResponse{
		Templates:    items,
		Runtime:      runtime,
		KVMAvailable: h.images.KVMAvailable(),
		ImagesDir:    h.images.Dir(),
	})
}

// Prepare starts a download job for a template.
func (h *TemplateHandler) Prepare(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	t, ok := templates.Find(id)
	if !ok {
		respondError(w, http.StatusNotFound, "Template not found.")
		return
	}

	job, err := h.images.Prepare(t)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusAccepted, prepareResponse{Job: job})
}

// Jobs returns recent download jobs.
func (h *TemplateHandler) Jobs(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string][]images.Job{"jobs": h.images.Jobs()})
}

// Job returns a single download job.
func (h *TemplateHandler) Job(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, ok := h.images.Job(id)
	if !ok {
		respondError(w, http.StatusNotFound, "Download job not found.")
		return
	}
	respondJSON(w, http.StatusOK, map[string]images.Job{"job": job})
}

// Deploy creates a VM from a prepared template and optionally starts it.
func (h *TemplateHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	var req deployFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	t, ok := templates.Find(req.TemplateID)
	if !ok {
		respondError(w, http.StatusNotFound, "Template not found.")
		return
	}
	if t.Category == "runtime" {
		respondError(w, http.StatusBadRequest, "The runtime template cannot be deployed as a VM.")
		return
	}

	var kernelPath, rootfsPath string
	for _, state := range h.images.State(t) {
		if !state.Ready {
			respondError(w, http.StatusConflict, "The template images are not on disk yet. Download them from the Deploy page first.")
			return
		}
		switch state.Kind {
		case templates.KindKernel:
			kernelPath = state.Path
		case templates.KindRootfs:
			rootfsPath = state.Path
		}
	}

	name := req.Name
	if name == "" {
		name = t.ID + "-" + time.Now().Format("150405")
	}
	cpus := req.CPUs
	if cpus <= 0 {
		cpus = t.DefaultCPUs
	}
	mem := req.MemoryMB
	if mem <= 0 {
		mem = t.DefaultMemMB
	}

	config := models.VMConfig{
		Name:       name,
		KernelPath: kernelPath,
		KernelOpts: t.KernelOpts,
		RootDrive:  models.Drive{Path: rootfsPath, ReadOnly: false},
		CPUs:       cpus,
		MemoryMB:   mem,
		LogLevel:   "Debug",
	}

	created, err := h.manager.Create(config)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create the VM from the template.")
		return
	}

	resp := deployResponse{VM: created, Started: false}
	if req.AutoStart {
		if err := h.manager.Start(created.ID); err != nil {
			resp.StartError = err.Error()
		} else {
			resp.Started = true
		}
		if fresh, err := h.manager.Get(created.ID); err == nil {
			resp.VM = fresh
		}
	}

	respondJSON(w, http.StatusCreated, resp)
}
