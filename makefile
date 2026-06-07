SHELL    := /bin/bash
MIGRATE  := $(shell go env GOPATH)/bin/migrate
GOLANGCI := $(shell go env GOPATH)/bin/golangci-lint
BINARY   := bin/server

ifneq (,$(wildcard .env))
    include .env
    export
endif

DB_URL=mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: run build clean \
        lint lint-fix fmt fmt-check \
        test test-v \

# ── Development ───────────────────────────────────────────────────

run:
	go run ./cmd/server/...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" \
		-o $(BINARY) ./cmd/server/...
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
