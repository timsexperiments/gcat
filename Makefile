BINARY_PATH := build/gcat/gcat
GO ?= go
BUILD_FLAGS := -ldflags="-s -w"
COVERAGE_DIR := coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out

.PHONY: all build install clean test coverage fmt lint vet tidy generate help docker

all: build

build:
	mkdir -p build/gcat
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd/gcat

install:
	$(GO) install ./cmd/gcat

clean:
	rm -rf build coverage

test:
	$(GO) test ./...

coverage:
	mkdir -p $(COVERAGE_DIR)
	$(GO) test -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_DIR)/coverage.html

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

generate:
	$(GO) generate ./...

help:
	@echo "Makefile commands:"
	@echo "  all        - Build the application"
	@echo "  build      - Build the application"
	@echo "  install    - Install the application"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  coverage   - Generate a test coverage report"
	@echo "  fmt        - Format the code"
	@echo "  lint       - Run linting (golangci-lint)"
	@echo "  vet        - Run go vet"
	@echo "  tidy       - Tidy up go.mod and go.sum"
	@echo "  generate   - Run code generation (go generate)"
	@echo "  docker     - Build docker image if available"
