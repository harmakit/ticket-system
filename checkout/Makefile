CONTAINER=checkout

DC=docker compose -p ticket-system --env-file ./.env -f ./.docker/docker-compose.yml
DC_DEBUG = docker compose -p ticket-system --env-file ./.env -f ./.docker/docker-compose.yml -f ./.docker/docker-compose.dev.yml

EXEC=$(DC) exec $(CONTAINER)
BINDIR=$(shell pwd)/bin

.DEFAULT_GOAL := help

### docker

up: env ## Run containers
	${DC} up -d

up-debug: env ## Run containers with debugger
	${DC_DEBUG} up -d

build: env generate-proto ## Build Docker image
	${DC} build

build-debug: env generate-proto ## Build Docker image with debugger
	${DC_DEBUG} build

down: ## Stop containers
	${DC} down

migrate: ## Run migrations in running container
	${EXEC} /migrate

### local

build-local: env generate-proto ## Build locally
	go build -o ./bin/server ./cmd/server && \
	go build -o ./bin/migrate ./cmd/migrate

lint: install-lint ## Run linter
	${BINDIR}/golangci-lint run ./...

test: ## Run tests
	go test ./...

pre-commit: ## Use before commit
	lint test

### helpers

help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env:
	if [ ! -f .env ]; then \
		cp .env.dist .env; \
	fi

bindir:
	test -d ${BINDIR} || mkdir -p ${BINDIR}

install-protoc: bindir
	test -f $(BINDIR)/protoc-gen-go || GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	test -f $(BINDIR)/protoc-gen-go-grpc || GOBIN=$(BINDIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

install-vendor-proto:
	mkdir -p vendor-proto
	@if [ ! -d vendor-proto/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
		mkdir -p  vendor-proto/google/ &&\
		mv vendor-proto/googleapis/google/api vendor-proto/google &&\
		rm -rf vendor-proto/googleapis ;\
	fi
	@if [ ! -d vendor-proto/google/protobuf ]; then\
		git clone https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
		mkdir -p  vendor-proto/google/protobuf &&\
		mv vendor-proto/protobuf/src/google/protobuf/*.proto vendor-proto/google/protobuf &&\
		rm -rf vendor-proto/protobuf ;\
	fi

generate-proto: install-protoc install-vendor-proto
	mkdir -p pkg/v1/api
	protoc -I api/v1 -I vendor-proto \
		--go_out=pkg/v1/api --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/v1/api --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		api/v1/*.proto
	mkdir -p pkg/v1/message
	protoc -I message/v1 -I vendor-proto \
		--go_out=pkg/v1/message --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		message/v1/*.proto

install-lint:
	test -f $(BINDIR)/golangci-lint || GOBIN=$(BINDIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3