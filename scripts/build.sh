#!/usr/bin/env bash

cd "$(dirname "$0")/.." || exit
CGO_ENABLED=0 go build -o bin/filepath_adapter main.go
