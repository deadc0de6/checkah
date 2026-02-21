#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2021, deadc0de6

set -e

# deps
echo "install deps..."
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/mgechev/revive@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser/v2@latest

go mod tidy

# linting
echo "go fmt..."
go fmt ./...

#echo "golint..."
#golint -set_exit_status ./...
echo "golangci-lint..."
golangci-lint run

echo "revive..."
revive -set_exit_status -config ./revive.toml ./...

echo "staticcheck..."
staticcheck ./...

echo "go vet..."
go vet ./...

# tests
make clean
make
