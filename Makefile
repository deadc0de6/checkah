GO111MODULE=on
GOBIN=$(shell pwd)/bin
INSTALL_FLAG=-v -ldflags "-s -w"

all: build

build: env deps
	@GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) go install $(INSTALL_FLAG) ./...

clean:
	@rm -rf $(GOBIN)

deps:
	go mod tidy

env:
	@go env GO111MODULE=on

.PHONY: build env clean all deps
