BACKEND = backend

GOBIN = $(CURDIR)/$(BACKEND)/bin

DBIN = ./$(BACKEND)/bin
LINT = $(DBIN)/golangci-lint
GOOSE = $(DBIN)/goose

export GOBIN

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "%-15s %s\n", $$1, $$2}'

all: ffvl brun ## fmt + fix + lint + build + run

install-lint: ## install golangci-lint
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOBIN)

install-goose: ## install goose
	go install github.com/pressly/goose/v3/cmd/goose@latest

fmt: ## format
	cd $(BACKEND) && go fmt ./...

fix: ## fix
	cd $(BACKEND) && go fix ./...

vet: ## vet
	cd $(BACKEND) && go vet ./...

lint: ## lint
	cd $(BACKEND) && ./bin/golangci-lint run

ffvl: fmt fix vet lint ## fmt + fix + vet + lint

build: ## build
	cd $(BACKEND) && CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/main ./cmd/api

run: ## run
	$(DBIN)/main

brun: build run ## build + run

up: ## up compose
	docker compose -f deployments/docker-compose.yaml up -d --build

down: ## down compose
	docker compose -f deployments/docker-compose.yaml down --volumes

ps: ## ps compose
	docker compose -f deployments/docker-compose.yaml ps -a
