#!/bin/bash

# Get the directory of the script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

CGO_CFLAGS="-I${SCRIPT_DIR}/td/tdlib/include" \
CGO_LDFLAGS="-Wl,-rpath,${SCRIPT_DIR}/td/tdlib/lib -L${SCRIPT_DIR}/td/tdlib/lib -ltdjson" \
go run ./main.go