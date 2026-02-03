#!/bin/bash
<<<<<<< HEAD
# Entry point for the Go-based disk analyzer binary bundled with Mole.
=======
# Mole - Analyze command.
# Runs the Go disk analyzer UI.
# Uses bundled analyze-go binary.
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_BIN="$SCRIPT_DIR/analyze-go"
if [[ -x "$GO_BIN" ]]; then
    exec "$GO_BIN" "$@"
fi

echo "Bundled analyzer binary not found. Please reinstall Mole or run mo update to restore it." >&2
exit 1
