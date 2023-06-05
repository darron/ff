BINARY_NAME ?= ff
CONTAINER_NAME ?= darron/ff
DO_BOX ?= 127.0.0.1

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

deploy-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_COMMAND)
	scp bin/ff root@$(DO_BOX):/root/ff

deploy: clean deploy-binary
	scp views/* root@$(DO_BOX):/root/views/
	scp -rp public/* root@$(DO_BOX):/root/public/
	scp import.csv root@$(DO_BOX):/root/import.csv
	scp deploy/tls_cache/* root@$(DO_BOX):/root/tls_cache/
	scp deploy/setup.sh root@$(DO_BOX):/root/setup.sh

grab-files:
	scp root@$(DO_BOX):/root/tls_cache/* deploy/tls_cache/
	scp root@$(DO_BOX):/root/setup.sh deploy/setup.sh

.PHONY: help all deps clean build gzip release unit lint docker docker-curl