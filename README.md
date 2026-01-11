Agni
===

[![Build status](https://badge.buildkite.com/92fe02b4bd9564be0f7ea21d1ee782f6a6fe55cbd5465e3480.svg?branch=master)](https://buildkite.com/firecracker-microvm/firectl)

Agni is a powerful tool for managing Firecracker MicroVMs. It provides both a command-line interface (CLI) and a modern web-based graphical user interface (GUI) for creating, managing, and monitoring virtual machines.

**Agni** (Sanskrit: अग्नि) means "fire" - a fitting name for a Firecracker management tool.

## Features

### CLI Mode
- Run Firecracker MicroVMs from the command line
- Full console access to VMs
- Read/write access to filesystems
- Network connectivity via TAP devices
- Vsock support

### GUI Mode (Web Interface)
- **VM Management**: Create, start, stop, and delete VMs through a web interface
- **Real-time Monitoring**: View VM status, metrics, and logs in real-time
- **Configuration Profiles**: Save and reuse VM configurations
- **User Authentication**: JWT-based authentication with role-based access
- **REST API**: Full API access for automation and integration
- **Responsive Design**: Works on desktop and mobile browsers

## Building

The default Makefile rule executes `go build` and relies on the Go toolchain
installed on your computer.
_We use [go modules](https://github.com/golang/go/wiki/Modules), and building
requires Go 1.23 or newer._

### CLI Binary

```bash
make
```

This creates the `agni` binary.

### GUI Binary (with embedded frontend)

```bash
make gui
```

This creates the `agni-gui` binary with the embedded web interface.

### Docker Build

If you do not have a new-enough Go toolchain installed:

```bash
make build-in-docker
```

## Usage

### CLI Mode

You'll need a [firecracker](https://github.com/firecracker-microvm/firecracker) binary,
an uncompressed Linux kernel image (`vmlinux`), and a root filesystem image.

```
Usage:
  agni [OPTIONS]

Application Options:
      --firecracker-binary=     Path to firecracker binary
      --kernel=                 Path to the kernel image (default: ./vmlinux)
      --kernel-opts=            Kernel commandline (default: ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules)
      --root-drive=             Path to root disk image, optionally suffixed with :ro or :rw
      --root-partition=         Root partition UUID
      --add-drive=              Path to additional drive, suffixed with :ro or :rw, can be specified multiple times
      --tap-device=             NIC info, specified as DEVICE/MAC
      --vsock-device=           Vsock interface, specified as PATH:CID. Multiple OK
      --vmm-log-fifo=           FIFO for firecracker logs
      --log-level=              vmm log level (default: Debug)
      --metrics-fifo=           FIFO for firecracker metrics
  -t, --disable-smt             Disable CPU Simultaneous Multithreading
  -c, --ncpus=                  Number of CPUs (default: 1)
      --cpu-template=           Firecracker CPU Template (C3 or T2)
  -m, --memory=                 VM memory, in MiB (default: 512)
      --metadata=               Firecracker Metadata for MMDS (json)
  -l, --firecracker-log=        pipes the fifo contents to the specified file
  -s, --socket-path=            path to use for firecracker socket
  -d, --debug                   Enable debug output

GUI Options:
  -g, --gui                     Start in GUI mode (web interface)

Help Options:
  -h, --help                    Show this help message
```

### CLI Example

```bash
agni \
  --kernel=~/bin/vmlinux \
  --root-drive=/images/image-debootstrap.img -t \
  --cpu-template=T2 \
  --firecracker-log=~/firecracker-vmm.log \
  --kernel-opts="console=ttyS0 noapic reboot=k panic=1 pci=off nomodules rw" \
  --vsock-device=root:3 \
  --metadata='{"foo":"bar"}'
```

### GUI Mode

Start the web interface:

```bash
agni --gui
```

Or use the GUI binary:

```bash
agni-gui
```

The web interface will be available at `http://localhost:8080`.

#### First-time Setup

1. Open `http://localhost:8080` in your browser
2. Create an admin account at the setup page
3. Log in and start managing VMs!

#### API Endpoints

The GUI mode exposes a REST API:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/auth/login` | POST | User login |
| `/api/auth/setup` | POST | Initial admin setup |
| `/api/vms` | GET | List all VMs |
| `/api/vms` | POST | Create a new VM |
| `/api/vms/:id` | GET | Get VM details |
| `/api/vms/:id/start` | POST | Start a VM |
| `/api/vms/:id/stop` | POST | Stop a VM |
| `/api/vms/:id/logs` | GET | Stream VM logs |
| `/api/configs` | GET/POST | Manage configurations |

## Development

### Run API server with frontend dev server

```bash
make dev
```

This starts the API server and the frontend development server with hot reload.

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

## Getting Started on AWS

- Create an `m5d.metal` instance using Amazon Linux 2
- Build or download agni binary
- Get Firecracker binary:

  ```bash
  curl -Lo firecracker https://github.com/firecracker-microvm/firecracker/releases/download/v0.16.0/firecracker-v0.16.0
  chmod +x firecracker
  sudo mv firecracker /usr/local/bin/firecracker
  ```

- Give read/write access to KVM:

  ```bash
  sudo setfacl -m u:${USER}:rw /dev/kvm
  ```

- Download kernel and root filesystem:

  ```bash
  curl -fsSL -o hello-vmlinux.bin https://s3.amazonaws.com/spec.ccfc.min/img/hello/kernel/hello-vmlinux.bin
  curl -fsSL -o hello-rootfs.ext4 https://s3.amazonaws.com/spec.ccfc.min/img/hello/fsfiles/hello-rootfs.ext4
  ```

- Create microVM:

  ```bash
  ./agni \
    --kernel=hello-vmlinux.bin \
    --root-drive=hello-rootfs.ext4
  ```

## Testing

By default the tests require the agni binary to be built and a kernel image
to be present. The integration tests look for the binary and kernel image in
the root directory. By default it will look for vmlinux kernel image. This can
be overwritten by setting the environment variable `KERNELIMAGE` to the desired
path. To disable these tests simply set the environment variable
`SKIP_INTEG_TEST=1`.

```bash
make test                    # Run all tests
SKIP_INTEG_TEST=1 make test  # Skip integration tests
```

## Questions?

Please use
[GitHub issues](https://github.com/firecracker-microvm/firectl/issues)
to report problems, discuss roadmap items, or make feature requests.

If you've discovered an issue that may have security implications to
users or developers of this software, please do not report it using
GitHub issues, but instead follow
[Firecracker's security reporting guidelines](https://github.com/firecracker-microvm/firecracker/blob/main/SECURITY.md).

Other discussion: For general discussion, please get in touch with the [Firecracker community.](https://github.com/firecracker-microvm/firecracker?tab=readme-ov-file#faq--contact).
