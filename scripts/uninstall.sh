#!/usr/bin/env bash
set -euo pipefail

LABEL="com.github.rin2yh.tiny-deck"
BINARY_PATH="${HOME}/.local/bin/tiny-deck-host"
PLIST_DST="${HOME}/Library/LaunchAgents/${LABEL}.plist"

echo "==> launchctl bootout"
launchctl bootout "gui/${UID}/${LABEL}" 2>/dev/null || true

echo "==> remove plist"
rm -f "${PLIST_DST}"

echo "==> remove binary"
rm -f "${BINARY_PATH}"

echo "done."
