.DEFAULT_GOAL := build

APP_NAME=v2ray-mu
APP_CMD_DIR=cmd/$(APP_NAME)
APP_BINARY=bin/$(APP_NAME)
APP_BINARY_UNIX=bin/$(APP_NAME)_unix_amd64

all: build

.PHONY: test
test: ## test
	go test -v ./...


.PHONY: build
build: ## build
	CGO_ENABLED=0 go build -o $(APP_BINARY) -v cmd/$(APP_NAME)/main.go

