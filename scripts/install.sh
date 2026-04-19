#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

LABEL="com.github.rin2yh.tiny-deck"
BINARY_PATH="${HOME}/.local/bin/tiny-deck-host"
LOG_PATH="${HOME}/Library/Logs/tiny-deck.log"
PLIST_SRC="${SCRIPT_DIR}/${LABEL}.plist"
PLIST_DST="${HOME}/Library/LaunchAgents/${LABEL}.plist"

echo "==> build binary"
mkdir -p "$(dirname "${BINARY_PATH}")"
(cd "${PROJECT_ROOT}" && go build -o "${BINARY_PATH}" ./cmd/host_daemon)

echo "==> prepare log directory"
mkdir -p "$(dirname "${LOG_PATH}")"

echo "==> render plist -> ${PLIST_DST}"
mkdir -p "$(dirname "${PLIST_DST}")"
sed \
    -e "s|{{BINARY_PATH}}|${BINARY_PATH}|g" \
    -e "s|{{LOG_PATH}}|${LOG_PATH}|g" \
    "${PLIST_SRC}" > "${PLIST_DST}"

echo "==> launchctl bootstrap"
launchctl bootout "gui/${UID}/${LABEL}" 2>/dev/null || true
launchctl bootstrap "gui/${UID}" "${PLIST_DST}"

echo "done. plug the device in; logs: ${LOG_PATH}"
