# Unreleased

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
