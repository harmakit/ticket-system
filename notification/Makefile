CONTAINER=notification

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

build: env ## Build Docker image
	${DC} build

build-debug: env ## Build Docker image with debugger
	${DC_DEBUG} build

down: ## Stop containers
	${DC} down

### local

build-local: env ## Build locally
	go build -o ./bin/app ./cmd/app

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

install-lint:
	test -f $(BINDIR)/golangci-lint || GOBIN=$(BINDIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3