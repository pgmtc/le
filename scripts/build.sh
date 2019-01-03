#!/usr/bin/env bash
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-mac-amd64 ./cmd/orchard
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-linux-amd64 ./cmd/orchard
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/orchard-windows-amd64 ./cmd/orchard