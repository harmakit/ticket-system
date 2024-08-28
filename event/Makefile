CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
LINTVER=v1.60.2
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=ticket-system/events

.PHONY: all

all: build test lint

build: bindir
	env GOOS=linux GOARCH=amd64 go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

precommit: build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

install-go-deps:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BINDIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-go-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

vendor-proto:
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

generate:
	mkdir -p pkg/event_v1
	protoc -I api/event/v1 -I vendor-proto \
	--go_out=pkg/event_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	api/event/v1/service.proto