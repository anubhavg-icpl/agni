# Agni packaging

One static binary, two package formats. The Go binary embeds the web UI
(`go build -tags gui`), so a package needs nothing but the binary itself,
a systemd unit, and a KVM-capable host.

## Layout

- `nfpm.yaml` - shared package definition for deb and rpm
- `agni.service` - systemd unit (serves on :8080, data in /var/lib/agni)
- `deb/` - Debian maintainer scripts (postinst, prerm)
- `rpm/` - RPM scriptlets (post.sh, preun.sh)

## Build a package locally

    make package VERSION=0.1.0

which is equivalent to:

    (cd frontend && npm ci && npm run build)
    go build -tags gui -o agni .
    VERSION=0.1.0 nfpm package -f packaging/nfpm.yaml -p deb -t agni_0.1.0_amd64.deb
    VERSION=0.1.0 nfpm package -f packaging/nfpm.yaml -p rpm -t agni_0.1.0_x86_64.rpm

## Install

    # Debian/Ubuntu
    sudo apt install ./agni_0.1.0_amd64.deb
    # Fedora/RHEL/Arch (via dnf)
    sudo dnf install ./agni_0.1.0_x86_64.rpm

    sudo systemctl enable --now agni
    # UI at http://localhost:8080

Releases are attached automatically when a `v*` tag is pushed; see
`.github/workflows/release.yml`.
