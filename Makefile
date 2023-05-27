BINARY_NAME ?= ff
CONTAINER_NAME ?= darron/ff

BUILD_COMMAND=-mod=vendor -o bin/$(BINARY_NAME) ../$(BINARY_NAME)
UNAME=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m)
GIT_SHA=$(shell git rev-parse HEAD)

all: build

deps: ## Install all dependencies.
	go mod vendor
	go mod tidy -compat=1.20

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Remove compiled binaries.
	rm -f bin/$(BINARY_NAME) || true
	rm -f bin/$(BINARY_NAME)*gz || true

docker: ## Build Docker image
	docker buildx build . --platform linux/amd64,linux/arm64 -t $(CONTAINER_NAME):$(GIT_SHA) --push

build: clean
	go build $(BUILD_COMMAND)

rebuild: clean ## Force rebuild of all packages.
	go build -a $(BUILD_COMMAND)

linux: clean ## Cross compile for linux.
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_COMMAND)

gzip: ## Compress current compiled binary.
	gzip bin/$(BINARY_NAME)
	mv bin/$(BINARY_NAME).gz bin/$(BINARY_NAME)-$(UNAME)-$(ARCH).gz

release: build gzip ## Full release process.

unit: ## Run unit tests.
	go test -mod=vendor -cover -race -short ./... -v

lint: ## See https://github.com/golangci/golangci-lint#install for install instructions
	golangci-lint run ./...

.PHONY: help all deps clean build gzip release unit lint docker docker-curl