LOCAL_BIN:=$(CURDIR)/bin

# Install golang dependencies for the project
.PHONY: install-go-deps
install-go-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Generate the protobuf files
.PHONY: protos
protos:
	rm -r pkg/v1
	mkdir -p pkg/v1
	protoc --proto_path ./ \
	--go_out=pkg/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/sysmon.proto

# Run the sysmon service locally
.PHONY: run
run:
	go run ./cmd/sysmon/

# Run the sysmon service in a docker container
.PHONY: docker-up
run-ubuntu:
	docker-compose up --build ubuntu

# Build the sysmon binary for the host OS
.PHONY: build
build:
	go build -o bin/sysmon ./cmd/sysmon/

# Run all tests for the project
.PHONY: tests
tests:
	go test ./... -race -count 100

# Run the linter for the project
.PHONY: lint
lint:
	golangci-lint run

