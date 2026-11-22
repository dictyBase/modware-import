#!/usr/bin/env bash
# Run from any directory: bash scripts/check-gopls.sh.
# Lex Go comments and literals, not a fixed header window, to identify generated sources.
set -euo pipefail
cd "$(dirname "$0")/.."
python3 scripts/check_gopls.py
