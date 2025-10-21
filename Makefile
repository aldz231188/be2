SHELL := /bin/bash
APP := app
BIN := ./bin/$(APP)

DB_URL=postgres://app:qwe@localhost:5432/appdb?sslmode=disable
MIGR_DIR=./internal/infra/db/migrations

.PHONY: tidy build run sqlc lint test up down
.PHONY: migrate-create migrate-up migrate-down migrate-version

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=add_table"; exit 1)
	migrate create -ext sql -dir $(MIGR_DIR) -format 20060102150405 $(name)

migrate-up:
	migrate -path $(MIGR_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGR_DIR) -database "$(DB_URL)" down 1

migrate-force:
	@test -n "$(version)" || (echo "Usage: make migrate-force version=v№"; exit 1)
	migrate -path $(MIGR_DIR) -database "$(DB_URL)" force $(version)





tidy:
	go mod tidy


build:
	go build -o $(BIN) ./cmd/server


run:
	DATABASE_DSN=$${DATABASE_DSN:-postgres://app:qwe@localhost:5432/appdb?sslmode=disable} \
	HTTP_ADDR=:8080 $(BIN)


sqlc:
	sqlc generate

lint:
	golangci-lint run


test:
	go test ./...


up:
	docker compose up -d --build


down:
	docker compose down -v
