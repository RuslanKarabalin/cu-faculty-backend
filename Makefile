GOBIN = $(CURDIR)/bin

DBIN = ./bin
LINT = $(DBIN)/golangci-lint
OAPI = $(DBIN)/oapi-codegen
GOOSE = $(DBIN)/goose

export GOBIN

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

up:
	docker compose -f infra/docker-compose.yaml up -d --build

down:
	docker compose -f infra/docker-compose.yaml down --volumes

ps:
	docker compose -f infra/docker-compose.yaml ps -a
