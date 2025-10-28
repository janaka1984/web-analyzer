SHELL := /bin/bash

.PHONY: run build test up down logs fmt tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

up:
	docker compose up --build -d

down:
	docker compose down -v

logs:
	docker compose logs -f api

test:
	go test ./... -cover
