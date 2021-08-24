GO111MODULE=on
GOBIN=$(shell pwd)/bin
INSTALL_FLAG=-v -ldflags "-s -w"
OS=linux
NAME=checkah

all: build

build: env deps
	@GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) go install $(INSTALL_FLAG) ./...

build-all:
	GO111MODULE=on GOOS=$(OS) GOARCH=arm   cd cmd/checkah; go build -v -o ../../bin/$(NAME)-$(OS)-arm
	GO111MODULE=on GOOS=$(OS) GOARCH=arm64 cd cmd/checkah; go build -v -o ../../bin/$(NAME)-$(OS)-arm64
	GO111MODULE=on GOOS=$(OS) GOARCH=386   cd cmd/checkah; go build -v -o ../../bin/$(NAME)-$(OS)-386
	GO111MODULE=on GOOS=$(OS) GOARCH=amd64 cd cmd/checkah; go build -v -o ../../bin/$(NAME)-$(OS)-amd64

clean:
	@rm -rf $(GOBIN)

deps:
	go mod tidy

env:
	@go env GO111MODULE=on

.PHONY: build env clean all deps
