#!/bin/bash

set -eu

GOOS=js GOARCH=wasm go build -o ebiten-memory-usage.wasm
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
