#!/usr/bin/env bash

cd "$(dirname "$0")/.." || exit
CGO_ENABLED=0 go build -o bin/volume_injector main.go
