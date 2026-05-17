GOBIN = $(CURDIR)/bin

DBIN = ./bin
LINT = $(DBIN)/golangci-lint
GOOSE = $(DBIN)/goose

export GOBIN

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "%-15s %s\n", $$1, $$2}'

all: ffl brun ## fmt + fix + lint + build + run

install-lint: ## install golangci-lint
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOBIN)

install-goose: ## install goose
	go install github.com/pressly/goose/v3/cmd/goose@latest

fmt: ## format
	go fmt ./...

fix: ## fix
	go fix ./...

lint: ## lint
	$(LINT) run

ffl: fmt fix lint ## fmt + fix + lint

build: ## build
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(DBIN)/main ./cmd/api

run: ## run
	$(DBIN)/main

brun: build run ## build + run

up: ## up compose
	docker compose -f deployments/docker-compose.yaml up -d --build

down: ## down compose
	docker compose -f deployments/docker-compose.yaml down --volumes

ps: ## ps compose
	docker compose -f deployments/docker-compose.yaml ps -a
