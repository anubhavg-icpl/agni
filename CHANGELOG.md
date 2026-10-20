# Unreleased

* Login redesigned as a cinematic ember scene: a GPU fragment-shader fire (rising fbm flame tongues, drifting embers, smoke wisps, breathing ground glow) paints the whole brand panel and blends into the UI, with a glowing AGNI wordmark, tagline and feature list, on a saturated hand-tuned fire palette (no filmic haze - only highlights above white get a soft knee); the animation loop now actually starts (the scene never received its initial visibility kick)

# 0.3.0

* Added an OCI registry template builder: 13 new bootable distros (Rocky 8/9, AlmaLinux 8/9, CentOS Stream 9/10, Oracle Linux 9, Amazon Linux 2023, Arch, Debian 12/13, Fedora 42, Nix) pulled straight from their registries, 34 templates total
* Added dependency preflight: `GET /api/system/dependencies` and an Environment card on the dashboard flag missing host tools (firecracker, KVM, mkfs.ext4, tar, unsquashfs) before the first deploy instead of at deploy time
* Added storage visibility: `GET /api/system/storage` and an Image storage card with per-template disk usage, cache size, and free disk
* Added destroy from any state: `DELETE /api/vms/{id}` force-stops a running VM first, and the trash action is always visible in the VM table and detail page
* Added default VM configuration profiles (Micro, Standard, Performance, Heavy) seeded on first boot for `/configs`
* Added product screenshots (webp) plus installation and GUI quick-start guides to the README
* Added packaging: `packaging/` with an nfpm config, systemd unit, and per-format script dirs; `.github/workflows/release.yml` publishes the binary plus deb and rpm packages on `v*` tags
* Packaging verified end to end locally: nfpm builds `agni_0.3.0_amd64.deb` and `agni_0.3.0_x86_64.rpm`; deb/rpm inspected for metadata (maintainer Anubhav Gain), maintainer scripts, systemd unit, and the packaged binary boots and answers /api/health; unit ExecStart uses the real `--gui` flags
* Fixed serial console showing nothing while typing: raw byte-level console forwarding (keystroke echoes stream live instead of after Enter) plus boot-log replay on connect, and `setsid cttyhack` gives shell-init guests a controlling tty (echo, backspace, Ctrl-C)
* First-run setup now confirms the created admin account on screen: username and a revealable password panel before continuing to the dashboard
* Fixed icons rendering as blank squares (double `ph-` prefix in Icon.svelte)
* Fixed read-only OCI layer directories breaking extraction, and 0000-mode files (gshadow) breaking mkfs.ext4
* Fixed init-less templates (Nix, Arch) missing /dev, /proc and mount(8): a static busybox is seeded into the image and PID 1 sets up mounts itself
* Fixed ext4 build failures on file-heavy images (Nix): inode count and size budgets with a doubling retry loop

* Rebuilt the web GUI on the Invinsense Interface Standard: OKLCH token set with light and dark themes, Geist and Geist Mono, Phosphor Bold icons, sidebar and topbar shell, professional copy throughout
* Added a Deploy page with one-click microVM templates: Alpine hello, Ubuntu 24.04, Ubuntu 22.04 (prebuilt ext4), Ubuntu 18.04 bionic, Debian with SSH, a boot-time benchmark microVM, and a diskless initramfs micro boot, plus the Firecracker v1.12.1 runtime, all pinned from the official CI asset bucket with streamed download progress, squashfs-to-ext4 conversion where needed, distro logos, and matching SSH keys saved next to the rootfs
* Added image management endpoints (catalog, prepare jobs, deploy-from-template with auto-start) and automatic firecracker runtime pickup from the data directory
* Switched the accent palette from ISOC blue to a Vercel-style monochrome: near-black primary on light, white primary on dark, neutral surfaces and borders
* Added an interactive serial console: every VM detail page embeds an xterm.js terminal over a bidirectional WebSocket on ttyS0, with reconnect backoff, boot-log replay, and lifecycle controls; typing can no longer tear down the session and stale firecracker sockets are cleaned on start
* Expanded the catalog to 21 templates with a rootfs tarball builder: Alpine 3.16-3.22 minirootfs, Ubuntu Base 18.04-24.04, Void Linux, and Gentoo stage3 systemd are converted to ext4 at download time (busybox inittab, systemd autologin, or runit agetty injected per family), with kernels deduplicated into a shared URL-keyed cache
* Reworked template card progress: per-artifact byte counters, a distinct building state with an indeterminate bar, error text with a retry button, and a logo fallback

# 0.2.0

* Upgraded to v1.0.0 of the firecracker-go-sdk
* Added initrd option (#84)
* Added read-only option for rootfs (#74)
* Added jailer support (#57)
* Added option to configure multiple NICs (#44)

*Note* - Although still compatible, this release does not provide full feature parity with firecracker v1.0.0. Majors features that are not included in this release are:

* Creating and loading snapshots
* Setting drive IO engine type
* Pause/Resume VM
* MMDS version 2

# 0.1.1 (unreleased)

* Add an `install` target to Makefile

# 0.1.0

* Initial release
