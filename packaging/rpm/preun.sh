#!/bin/sh
if [ "$1" = "0" ] && systemctl is-active --quiet agni 2>/dev/null; then
	systemctl stop agni || true
fi
