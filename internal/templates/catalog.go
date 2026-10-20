// Package templates provides the built-in, single-click deployable guest image
// catalog. Artifacts are pinned from Firecracker's public CI asset bucket
// (spec.ccfc.min) so a working kernel + rootfs pair can be fetched without
// building anything on the host.
package templates

import "runtime"

// Artifact kinds.
const (
	KindKernel  = "kernel"
	KindRootfs  = "rootfs"
	KindRuntime = "runtime"
)

// Artifact is one downloadable file that makes up a template.
type Artifact struct {
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	Filename    string   `json:"filename"`
	URL         string   `json:"url"`
	PostProcess string   `json:"post_process,omitempty"` // squashfs-to-ext4
	ToolDeps    []string `json:"tool_deps,omitempty"`    // host tools required by PostProcess
	SizeHintMB  int      `json:"size_hint_mb"`
}

// Template is a single-click deployable guest configuration.
type Template struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	OS           string     `json:"os"`
	Version      string     `json:"version"`
	Description  string     `json:"description"`
	Category     string     `json:"category"` // os | runtime
	Arch         string     `json:"arch"`
	Artifacts    []Artifact `json:"artifacts"`
	KernelOpts   string     `json:"kernel_opts"`
	DefaultCPUs  int64      `json:"default_cpus"`
	DefaultMemMB int64      `json:"default_memory_mb"`
	DocsURL      string     `json:"docs_url"`
}

const s3Base = "https://s3.amazonaws.com/spec.ccfc.min"

func guestArch() string {
	if runtime.GOARCH == "arm64" {
		return "aarch64"
	}
	return "x86_64"
}

// Catalog returns the templates available for this host architecture.
func Catalog() []Template {
	arch := guestArch()
	ci := s3Base + "/firecracker-ci/v1.12/" + arch

	return []Template{
		{
			ID:          "alpine-hello",
			Name:        "Alpine hello microVM",
			OS:          "Alpine Linux",
			Version:     "hello",
			Description: "Minimal busybox guest from the official Firecracker demo assets. Boots in under a second and opens a shell on the serial console. Best for smoke tests.",
			Category:    "os",
			Arch:        arch,
			Artifacts: []Artifact{
				{
					Name:       "Kernel (hello-vmlinux)",
					Kind:       KindKernel,
					Filename:   "hello-vmlinux.bin",
					URL:        s3Base + "/img/hello/kernel/hello-vmlinux.bin",
					SizeHintMB: 1,
				},
				{
					Name:       "Root filesystem (hello-rootfs)",
					Kind:       KindRootfs,
					Filename:   "hello-rootfs.ext4",
					URL:        s3Base + "/img/hello/fsfiles/hello-rootfs.ext4",
					SizeHintMB: 5,
				},
			},
			KernelOpts:   "console=ttyS0 reboot=k panic=1 pci=off",
			DefaultCPUs:  1,
			DefaultMemMB: 256,
			DocsURL:      "https://github.com/firecracker-microvm/firecracker/tree/main/docs",
		},
		{
			ID:          "ubuntu-24.04",
			Name:        "Ubuntu 24.04 microVM",
			OS:          "Ubuntu",
			Version:     "24.04",
			Description: "Ubuntu 24.04 root filesystem with the 6.1 LTS kernel, both from Firecracker CI. The published squashfs rootfs is converted to an ext4 image on first download.",
			Category:    "os",
			Arch:        arch,
			Artifacts: []Artifact{
				{
					Name:       "Kernel 6.1 (CI build)",
					Kind:       KindKernel,
					Filename:   "vmlinux-6.1.128",
					URL:        ci + "/vmlinux-6.1.128",
					SizeHintMB: 130,
				},
				{
					Name:        "Root filesystem (Ubuntu 24.04)",
					Kind:        KindRootfs,
					Filename:    "ubuntu-24.04.ext4",
					URL:         ci + "/ubuntu-24.04.squashfs",
					PostProcess: "squashfs-to-ext4",
					ToolDeps:    []string{"unsquashfs", "mkfs.ext4"},
					SizeHintMB:  1200,
				},
			},
			KernelOpts:   "root=/dev/vda console=ttyS0 reboot=k panic=1 pci=off",
			DefaultCPUs:  2,
			DefaultMemMB: 2048,
			DocsURL:      "https://github.com/firecracker-microvm/firecracker/blob/main/docs/getting-started.md",
		},
		{
			ID:          "firecracker-runtime",
			Name:        "Firecracker runtime",
			OS:          "Firecracker",
			Version:     "v1.12.1",
			Description: "The Firecracker VMM binary itself, from the official CI build. Required to start any VM. Installed into the Agni data directory and picked up automatically, no system install needed.",
			Category:    "runtime",
			Arch:        arch,
			Artifacts: []Artifact{
				{
					Name:       "firecracker v1.12.1",
					Kind:       KindRuntime,
					Filename:   "firecracker",
					URL:        ci + "/firecracker/firecracker-v1.12.1",
					SizeHintMB: 3,
				},
			},
			DocsURL: "https://github.com/firecracker-microvm/firecracker/releases",
		},
	}
}

// Find returns the template with the given ID.
func Find(id string) (Template, bool) {
	for _, t := range Catalog() {
		if t.ID == id {
			return t, true
		}
	}
	return Template{}, false
}
