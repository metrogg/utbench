#!/usr/bin/env bash

set -euo pipefail

PYTHON_BIN="${PYTHON_BIN:-python3}"

echo "[setup] python: $PYTHON_BIN"
"$PYTHON_BIN" --version

echo "[setup] installing dependencies"
"$PYTHON_BIN" -m pip install -U pyyaml pytest coverage mutmut

echo "[setup] creating runtime directories"
mkdir -p results

echo "[setup] done"
