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
	KindKey     = "ssh-key"
	KindInitrd  = "initrd"
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
	Logo         string     `json:"logo"` // static/brand/<logo>.svg
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

func x86Only(arch string) bool { return arch == "x86_64" }

// Catalog returns the templates available for this host architecture.
func Catalog() []Template {
	arch := guestArch()
	ci := s3Base + "/firecracker-ci/v1.12/" + arch
	ci110 := s3Base + "/firecracker-ci/v1.10/" + arch
	img := s3Base + "/img/" + arch

	catalog := []Template{
		{
			ID:          "alpine-hello",
			Name:        "Alpine hello microVM",
			OS:          "Alpine Linux",
			Version:     "hello",
			Logo:        "alpinelinux",
			Description: "Minimal busybox guest from the official Firecracker demo assets. Boots in under a second and opens a shell on the serial console. Best for smoke tests.",
			Category:    "os",
			Arch:        arch,
			Artifacts: []Artifact{
				{
					Name:       "Kernel (hello-vmlinux)",
					Kind:       KindKernel,
					Filename:   "hello-vmlinux.bin",
					URL:        s3Base + "/img/hello/kernel/hello-vmlinux.bin",
					SizeHintMB: 21,
				},
				{
					Name:       "Root filesystem (hello-rootfs)",
					Kind:       KindRootfs,
					Filename:   "hello-rootfs.ext4",
					URL:        s3Base + "/img/hello/fsfiles/hello-rootfs.ext4",
					SizeHintMB: 30,
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
			Version:     "24.04 LTS",
			Logo:        "ubuntu",
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
			Logo:        "firecracker",
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

	// Prebuilt ext4 images from the CI archive; x86_64 only.
	if x86Only(arch) {
		catalog = append(catalog[:len(catalog)-1], append([]Template{
			{
				ID:          "ubuntu-22.04",
				Name:        "Ubuntu 22.04 microVM",
				OS:          "Ubuntu",
				Version:     "22.04 LTS",
				Logo:        "ubuntu",
				Description: "Ubuntu 22.04 with the 5.10 LTS kernel, as a ready-made ext4 image from the Firecracker CI archive. No conversion step, and the CI SSH key is saved next to the rootfs.",
				Category:    "os",
				Arch:        arch,
				Artifacts: []Artifact{
					{
						Name:       "Kernel 5.10 (CI build)",
						Kind:       KindKernel,
						Filename:   "vmlinux-5.10.223",
						URL:        ci110 + "/vmlinux-5.10.223",
						SizeHintMB: 37,
					},
					{
						Name:       "Root filesystem (Ubuntu 22.04, prebuilt)",
						Kind:       KindRootfs,
						Filename:   "ubuntu-22.04.ext4",
						URL:        ci110 + "/ubuntu-22.04.ext4",
						SizeHintMB: 300,
					},
					{
						Name:       "CI SSH key",
						Kind:       KindKey,
						Filename:   "ubuntu-22.04.id_rsa",
						URL:        ci110 + "/ubuntu-22.04.id_rsa",
						SizeHintMB: 1,
					},
				},
				KernelOpts:   "root=/dev/vda console=ttyS0 reboot=k panic=1 pci=off",
				DefaultCPUs:  2,
				DefaultMemMB: 2048,
				DocsURL:      "https://github.com/firecracker-microvm/firecracker/blob/main/docs/getting-started.md",
			},
			{
				ID:          "ubuntu-bionic",
				Name:        "Ubuntu 18.04 classic demo",
				OS:          "Ubuntu",
				Version:     "18.04 LTS",
				Logo:        "ubuntu",
				Description: "The classic bionic demo guest from the original Firecracker launch assets. Boots to a serial login (root/root). The matching SSH key is saved next to the rootfs.",
				Category:    "os",
				Arch:        arch,
				Artifacts: []Artifact{
					{
						Name:       "Kernel (demo build)",
						Kind:       KindKernel,
						Filename:   "vmlinux-demo.bin",
						URL:        img + "/ubuntu/kernel/vmlinux.bin",
						SizeHintMB: 21,
					},
					{
						Name:       "Root filesystem (bionic, prebuilt)",
						Kind:       KindRootfs,
						Filename:   "bionic.rootfs.ext4",
						URL:        img + "/ubuntu/fsfiles/bionic.rootfs.ext4",
						SizeHintMB: 300,
					},
					{
						Name:       "Demo SSH key",
						Kind:       KindKey,
						Filename:   "bionic.rootfs.id_rsa",
						URL:        img + "/ubuntu/fsfiles/bionic.rootfs.id_rsa",
						SizeHintMB: 1,
					},
				},
				KernelOpts:   "root=/dev/vda console=ttyS0 reboot=k panic=1 pci=off",
				DefaultCPUs:  2,
				DefaultMemMB: 1024,
				DocsURL:      "https://github.com/firecracker-microvm/firecracker/tree/main/docs",
			},
			{
				ID:          "debian-ssh",
				Name:        "Debian microVM with SSH",
				OS:          "Debian",
				Version:     "buster",
				Logo:        "debian",
				Description: "Debian guest with an SSH server and the virtio balloon device, from the official Firecracker demo assets. Large 1 GB prebuilt ext4; the matching SSH key is saved next to the rootfs for use once a network interface is attached.",
				Category:    "os",
				Arch:        arch,
				Artifacts: []Artifact{
					{
						Name:       "Kernel (demo build)",
						Kind:       KindKernel,
						Filename:   "vmlinux-debian.bin",
						URL:        img + "/debian_with_ssh_and_balloon/kernel/vmlinux.bin",
						SizeHintMB: 23,
					},
					{
						Name:       "Root filesystem (Debian, prebuilt)",
						Kind:       KindRootfs,
						Filename:   "debian.rootfs.ext4",
						URL:        img + "/debian_with_ssh_and_balloon/fsfiles/debian.rootfs.ext4",
						SizeHintMB: 1000,
					},
					{
						Name:       "SSH key",
						Kind:       KindKey,
						Filename:   "debian.rootfs.id_rsa",
						URL:        img + "/debian_with_ssh_and_balloon/fsfiles/debian.rootfs.id_rsa",
						SizeHintMB: 1,
					},
				},
				KernelOpts:   "root=/dev/vda console=ttyS0 reboot=k panic=1 pci=off",
				DefaultCPUs:  2,
				DefaultMemMB: 2048,
				DocsURL:      "https://github.com/firecracker-microvm/firecracker/tree/main/docs",
			},
			{
				ID:          "boottime-minimal",
				Name:        "Boot-time benchmark microVM",
				OS:          "Linux",
				Version:     "minimal",
				Logo:        "linux",
				Description: "The tiny guest Firecracker CI uses for boot-time benchmarking. Minimal kernel and root filesystem, boots to a shell in milliseconds. Best for measuring cold-start time.",
				Category:    "os",
				Arch:        arch,
				Artifacts: []Artifact{
					{
						Name:       "Kernel (benchmark build)",
						Kind:       KindKernel,
						Filename:   "vmlinux-bench.bin",
						URL:        img + "/minimal/kernel/vmlinux.bin",
						SizeHintMB: 21,
					},
					{
						Name:       "Root filesystem (benchmark)",
						Kind:       KindRootfs,
						Filename:   "boottime-rootfs.ext4",
						URL:        img + "/minimal/fsfiles/boottime-rootfs.ext4",
						SizeHintMB: 30,
					},
				},
				KernelOpts:   "console=ttyS0 reboot=k panic=1 pci=off",
				DefaultCPUs:  1,
				DefaultMemMB: 128,
				DocsURL:      "https://github.com/firecracker-microvm/firecracker/tree/main/docs",
			},
			{
				ID:          "initrd-micro",
				Name:        "Initramfs micro boot",
				OS:          "Linux",
				Version:     "initrd",
				Logo:        "linux",
				Description: "No disk at all: the guest lives entirely in an initramfs loaded next to the kernel. The smallest, fastest boot Agni can deploy. The matching SSH key is saved alongside.",
				Category:    "os",
				Arch:        arch,
				Artifacts: []Artifact{
					{
						Name:       "Kernel (initrd build)",
						Kind:       KindKernel,
						Filename:   "vmlinux-initrd.bin",
						URL:        img + "/minimal_with_initrd/kernel/vmlinux.bin",
						SizeHintMB: 16,
					},
					{
						Name:       "Initramfs",
						Kind:       KindInitrd,
						Filename:   "initrd.img",
						URL:        img + "/minimal_with_initrd/fsfiles/initrd.img",
						SizeHintMB: 25,
					},
					{
						Name:       "SSH key",
						Kind:       KindKey,
						Filename:   "initrd.img.id_rsa",
						URL:        img + "/minimal_with_initrd/fsfiles/initrd.img.id_rsa",
						SizeHintMB: 1,
					},
				},
				KernelOpts:   "console=ttyS0 reboot=k panic=1 pci=off",
				DefaultCPUs:  1,
				DefaultMemMB: 128,
				DocsURL:      "https://github.com/firecracker-microvm/firecracker/tree/main/docs",
			},
		}, catalog[len(catalog)-1])...)
	}

	return catalog
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
