#!/bin/sh
set -e
systemctl daemon-reload >/dev/null 2>&1 || true
if systemctl is-enabled --quiet agni 2>/dev/null; then
	systemctl restart agni || true
fi
