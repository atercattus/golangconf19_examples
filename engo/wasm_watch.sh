#!/bin/bash
watch -n 1 GOOS=js GOARCH=wasm go build -o assets/main.wasm
