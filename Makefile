GOBIN=$(shell pwd)/bin
INSTALL_FLAG=-v -ldflags "-s -w"
NAME=checkah
MAIN=cmd/checkah/main.go

all: build

build: env deps
	@CGO_ENABLED=0 GOBIN=$(GOBIN) GO111MODULE=on go install $(INSTALL_FLAG) ./...

build-all:
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=arm   go build -v -o ./bin/$(NAME)-$(OS)-arm $(MAIN)
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=arm64 go build -v -o ./bin/$(NAME)-$(OS)-arm64 $(MAIN)
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=386   go build -v -o ./bin/$(NAME)-$(OS)-386 $(MAIN)
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)-$(OS)-amd64 $(MAIN)

clean:
	@rm -rf $(GOBIN)

deps:
	go mod tidy

env:
	@go env -w GO111MODULE=on

.PHONY: build env clean all deps build-all
