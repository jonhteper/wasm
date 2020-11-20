#!/bin/bash

# This script create a server where wasm tests can executed
file=$(go env GOROOT)
cp $file/misc/wasm/wasm_exec.js wasm_exec.js
echo "Copy file!"

echo "Compiling..."
GOOS=js GOARCH=wasm go build -o app.wasm tests.go
echo "Compiled!"

echo "Initializing server ..."
go run server.go

