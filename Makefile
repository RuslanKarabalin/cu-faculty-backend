BACKEND = backend

COMPOSE = docker compose -f deployments/docker-compose.yaml
ENV_FILE = deployments/.env

GOBIN = $(CURDIR)/$(BACKEND)/bin

DBIN = ./$(BACKEND)/bin
LINT = $(DBIN)/golangci-lint
GOOSE = $(DBIN)/goose

export GOBIN

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "%-15s %s\n", $$1, $$2}'

all: ffvl brun ## fmt + fix + vet + lint + build + run

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
	$(COMPOSE) up -d --build

down: ## down compose
	$(COMPOSE) down --volumes

ps: ## ps compose
	$(COMPOSE) ps -a

garage-init: ## init garage layout, bucket and key (run once after `up`; re-run after `down` since it wipes volumes)
	@set -a; . $(ENV_FILE); set +a; \
	NODE=$$($(COMPOSE) exec -T garage /garage node id -q | cut -d@ -f1); \
	$(COMPOSE) exec -T garage /garage layout assign -z dc1 -c 1G $$NODE; \
	$(COMPOSE) exec -T garage /garage layout apply --version 1; \
	$(COMPOSE) exec -T garage /garage bucket create $$S3_BUCKET; \
	$(COMPOSE) exec -T garage /garage key import --yes -n cu-faculty-app $$S3_ACCESS_KEY $$S3_SECRET_KEY; \
	$(COMPOSE) exec -T garage /garage bucket allow --read --write $$S3_BUCKET --key $$S3_ACCESS_KEY
