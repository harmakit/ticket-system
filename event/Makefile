CONTAINER=event

DC=docker compose -p ticket-system -f ./.docker/docker-compose.yml
DC_DEBUG = docker compose -p ticket-system -f ./.docker/docker-compose.yml -f ./.docker/docker-compose.dev.yml

EXEC=$(DC) exec $(CONTAINER)
BINDIR=$(shell pwd)/bin

.PHONY: build

up: env
	${DC} up -d

up-debug: env
	${DC_DEBUG} up -d

build: env generate-proto
	${DC} build

build-debug: env generate-proto
	${DC_DEBUG} build

migrate:
	${EXEC} migrate

### helpers

env:
	if [ ! -f .docker/.env ]; then \
		cp .docker/.env.dist .docker/.env; \
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
	mkdir -p pkg/v1
	protoc -I api/v1 -I vendor-proto \
		--go_out=pkg/v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		api/v1/*.proto