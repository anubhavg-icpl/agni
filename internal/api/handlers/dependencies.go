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
	"os/exec"

	"github.com/anubhavg-icpl/agni/internal/images"
)

type dependencyStatus struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Required  bool   `json:"required"`
	Available bool   `json:"available"`
	Detail    string `json:"detail"`
	Hint      string `json:"hint,omitempty"`
}

// DependenciesHandler reports which host tools Agni needs before the first
// deploy, so problems surface in the UI instead of mid-deployment.
type DependenciesHandler struct {
	images *images.Manager
}

func NewDependenciesHandler(img *images.Manager) *DependenciesHandler {
	return &DependenciesHandler{images: img}
}

func (h *DependenciesHandler) Check(w http.ResponseWriter, r *http.Request) {
	deps := []dependencyStatus{}

	if path, ok := h.images.RuntimePath(); ok {
		deps = append(deps, dependencyStatus{
			ID: "firecracker", Label: "Firecracker VMM", Required: true,
			Available: true, Detail: path,
		})
	} else {
		deps = append(deps, dependencyStatus{
			ID: "firecracker", Label: "Firecracker VMM", Required: true,
			Available: false, Hint: "Open any template once and Agni downloads the runtime automatically.",
		})
	}

	kvm, kvmDetail := checkFile("/dev/kvm")
	deps = append(deps, dependencyStatus{
		ID: "kvm", Label: "KVM (/dev/kvm)", Required: kvm,
		Available: kvm, Detail: kvmDetail,
		Hint: "Without KVM, Firecracker cannot start microVMs on this host.",
	})

	for _, tool := range []struct {
		id, label, hint string
		required        bool
	}{
		{"mkfs.ext4", "mkfs.ext4 (e2fsprogs)", "Required to build root filesystems from template downloads.", true},
		{"tar", "tar", "Required to unpack distribution rootfs tarballs.", true},
		{"unsquashfs", "unsquashfs (squashfs-tools)", "Required only by the Ubuntu 24.04 squashfs template.", false},
		{"docker", "Docker", "Not required: Agni pulls container images straight from their registries.", false},
	} {
		if path, err := exec.LookPath(tool.id); err == nil {
			deps = append(deps, dependencyStatus{
				ID: tool.id, Label: tool.label, Required: tool.required,
				Available: true, Detail: path,
			})
		} else {
			deps = append(deps, dependencyStatus{
				ID: tool.id, Label: tool.label, Required: tool.required,
				Available: false, Hint: tool.hint,
			})
		}
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"dependencies": deps})
}

func checkFile(path string) (bool, string) {
	info, err := os.Stat(path)
	if err != nil {
		return false, ""
	}
	return true, path + " (" + info.Mode().String() + ")"
}
