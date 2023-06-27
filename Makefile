SHELL := /bin/bash
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

migrations:
	pushd config/migrations && go generate && popd

build: clean migrations
	go build $(BUILD_COMMAND)

rebuild: clean migrations ## Force rebuild of all packages.
	go build -a $(BUILD_COMMAND)

linux: clean migrations ## Compile for Linux/Docker.
	CGO_ENABLED=1 GOOS=linux go build -ldflags "-linkmode external -extldflags -static" $(BUILD_COMMAND)

# Cross compile locally on OS X.
# SQLite3 makes this required: brew install FiloSottile/musl-cross/musl-cross
# See: https://github.com/mattn/go-sqlite3#cross-compile
linux-amd64: clean migrations
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-linkmode external -extldflags -static" $(BUILD_COMMAND)

gzip: ## Compress current compiled binary.
	gzip bin/$(BINARY_NAME)
	mv bin/$(BINARY_NAME).gz bin/$(BINARY_NAME)-$(UNAME)-$(ARCH).gz

release: build gzip ## Full release process.

unit: ## Run unit tests.
	go test -mod=vendor -cover -race -short ./... -v

lint: ## See https://github.com/golangci/golangci-lint#install for install instructions
	golangci-lint run ./...

deploy-binary: linux-amd64
	scp bin/ff root@$(DO_BOX):/root/ff

deploy: clean deploy-binary
	scp views/* root@$(DO_BOX):/root/views/
	scp -rp public/* root@$(DO_BOX):/root/public/
	scp import.csv root@$(DO_BOX):/root/import.csv
	scp deploy/tls_cache/* root@$(DO_BOX):/root/tls_cache/
	scp deploy/cloudflare_tls/* root@$(DO_BOX):/root/cloudflare_tls/
	scp deploy/cloudflare.sh root@$(DO_BOX):/root/cloudflare.sh
	scp deploy/setup.sh root@$(DO_BOX):/root/setup.sh

grab-files:
	scp root@$(DO_BOX):/root/tls_cache/* deploy/tls_cache/
	scp root@$(DO_BOX):/root/setup.sh deploy/setup.sh

.PHONY: help all deps clean build gzip release unit lint docker docker-curl