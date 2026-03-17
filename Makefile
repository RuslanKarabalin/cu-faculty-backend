GOBIN = $(CURDIR)/bin

DBIN = ./bin
LINT = $(DBIN)/golangci-lint
OAPI = $(DBIN)/oapi-codegen
GOOSE = $(DBIN)/goose

export GOBIN

all: ffl brun

install-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOBIN)

install-oapi:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

fmt:
	go fmt ./...

fix:
	go fix ./...

lint:
	$(LINT) run

ffl: fmt fix lint

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(DBIN)/main ./cmd/api

run:
	$(DBIN)/main

brun: build run

up:
	docker compose up -d --build

down:
	docker compose down --volumes

ps:
	docker compose ps -a
