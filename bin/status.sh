#!/bin/bash
<<<<<<< HEAD
# Entry point for the Go-based system status panel bundled with Mole.
=======
# Mole - Status command.
# Runs the Go system status panel.
# Shows live system metrics.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_BIN="$SCRIPT_DIR/status-go"
if [[ -x "$GO_BIN" ]]; then
    exec "$GO_BIN" "$@"
fi

echo "Bundled status binary not found. Please reinstall Mole or run mo update to restore it." >&2
exit 1
