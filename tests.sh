#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2021, deadc0de6

set -e

# deps
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# tests
make clean
make
go fmt ./...
golint -set_exit_status ./...
go vet ./...
staticcheck ./...
