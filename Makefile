OS=$(shell uname -s)

# Setup and run app
all: setup install run
.PHONY: all

# Install required tools
setup:
	go get -u gopkg.in/alecthomas/gometalinter.v2
ifeq ($(OS), Darwin)
	brew install dep
else
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	gometalinter.v2 --install
.PHONY: setup

# Install dependencies
install:
	dep ensure -v
.PHONY: install

# Run app
run:
	go run main.go
.PHONY: run

# Run tests
test:
	go test -race -v ./...
.PHONY: test

# Run tests and show code coverage
cover:
	go test -race -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt
.PHONY: cover

# Run linters
lint:
	gometalinter.v2 --vendor --deadline 60s ./...
.PHONY: lint

# Build staticailly linked binary
build:
	CGO_ENABLED=0 go build -ldflags="-s -w"
.PHONY: build

# Test and lint code
ci: test lint

.DEFAULT_GOAL := build

.SILENT: # this has no purpose but to prevent echoing of commands for all targets
