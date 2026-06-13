SHELL    := /bin/bash
GOLANGCI := $(shell go env GOPATH)/bin/golangci-lint

.PHONY: run run-consumer build build-consumer clean \
        lint lint-fix fmt fmt-check \
        test test-v \
        docker-up docker-down docker-logs docker-build

# ── Development ───────────────────────────────────────────────────

run:
	go run ./cmd/server/...

run-consumer:
	go run ./cmd/consumer/...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server/...

build-consumer:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/consumer ./cmd/consumer/...

clean:
	rm -rf bin/

# ── Linting ───────────────────────────────────────────────────────

lint:
	$(GOLANGCI) run ./...

lint-fix:
	$(GOLANGCI) run --fix ./...

fmt:
	$(GOLANGCI) fmt ./...

fmt-check:
	$(GOLANGCI) fmt --diff ./...

# ── Testing ───────────────────────────────────────────────────────

test:
	go test ./...

test-v:
	go test ./... -v

# ── Docker ────────────────────────────────────────────────────────

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-down-v:
	docker compose down -v

docker-logs:
	docker compose logs -f

docker-logs-server:
	docker compose logs -f server

docker-logs-consumer:
	docker compose logs -f consumer