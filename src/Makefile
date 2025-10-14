SHELL := /bin/bash
APP := app
BIN := ./bin/$(APP)


.PHONY: tidy build run sqlc migrate-up migrate-down lint test up down


tidy:
go mod tidy


build:
go build -o $(BIN) ./cmd/server


run:
DATABASE_DSN=$${DATABASE_DSN:-postgres://app:app@localhost:5432/app?sslmode=disable} \
HTTP_ADDR=:8080 $(BIN)


sqlc:
sqlc generate


migrate-up:
migrate -path internal/infra/db/migrations -database $$DATABASE_DSN up


migrate-down:
migrate -path internal/infra/db/migrations -database $$DATABASE_DSN down 1


lint:
golangci-lint run


test:
go test ./...


up:
docker compose up -d --build


down:
docker compose down -v