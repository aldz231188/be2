SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

# ENV ?= dev
ENV_FILE := .env
-include $(ENV_FILE)

define EXPORT_DB_URI
[[ -f "$(DB_PASSWORD_FILE)" ]] || { echo "Missing password file: $(DB_PASSWORD_FILE)"; exit 1; }
DB_PASSWORD="$$(tr -d '\n' < "$(DB_PASSWORD_FILE)")"
export DB_URI="postgres://$(DB_USER):$${DB_PASSWORD}@$(DB_HOST_OUTSIDE):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
endef

# define SETUP_DB_URI
# 	@[ -f "$(DB_PASSWORD_FILE)" ] || { echo "❌ Err: Missing $(DB_PASSWORD_FILE)"; exit 1; }
# 	$(eval DB_PASSWORD := $(shell cat $(DB_PASSWORD_FILE)))
# 	$(eval DB_URI := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST_OUTSIDE):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE))
# endef


.PHONY: tidy build run sqlc lint test up down
.PHONY: migrate-create migrate-up migrate-down migrate-version

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=add_table"; exit 1)
	migrate create -ext sql -dir $(MIGR_DIR) -format 20060102150405 $(name)

migrate-up:
	$(EXPORT_DB_URI)
	migrate -path $(MIGR_DIR) -database $${DB_URI} up

migrate-down:
	$(EXPORT_DB_URI)
	migrate -path $(MIGR_DIR) -database $${DB_URI} down 1

migrate-force:
	$(EXPORT_DB_URI)
	@test -n "$(version)" || (echo "Usage: make migrate-force version=v№"; exit 1)
	migrate -path $(MIGR_DIR) -database $${DB_URI} force $(version)



tidy:
	go mod tidy


build:
	go build -o app ./cmd/server

# run:
# 	DB_PASSWORD_FILE="./secrets/db_password.txt" \
# 	JWT_SECRET_FILE="./secrets/jwt_secret.txt" $(BIN)



sqlc: sqlc-vet sqlc-gen

sqlc-vet:
	$(EXPORT_DB_URI)
# 	cd internal && sqlc vet
# 	cd services/clientsvc && sqlc vet
	cd services/authsvc && sqlc vet

sqlc-gen:
	$(EXPORT_DB_URI)
# 	cd internal && sqlc generate
# 	cd services/clientsvc && sqlc generate
	cd services/authsvc && sqlc generate




lint:
	golangci-lint run

test:
	go test ./...

COMPOSE := docker compose -f docker-compose.yml 

up-prod:
	$(COMPOSE) pull app
	$(COMPOSE) up -d --no-build --pull=always
up:
	$(COMPOSE) up -d --build 
down:
	$(COMPOSE) down -v
stop:
	$(COMPOSE) stop
start:
	$(COMPOSE) start

logs-all:
	$(COMPOSE) logs app nginx clientsvc clientsvc_db clientsvc_migrator authsvc_db authsvc_migrator

logs-migrator:
	$(COMPOSE) logs  migrator clientsvc_migrator

logs-app:
	$(COMPOSE) logs  app clientsvc authsvc

logs-nginx:
	$(COMPOSE) logs  nginx
	
logs-db:
	$(COMPOSE) logs  db clientsvc_db

graph_cp:
	$(COMPOSE) cp app:/tmp/graph.dot ./graph.dot

ci:
	act -P ubuntu-latest=catthehacker/ubuntu:act-latest


cert:
	set -a
	source $(ENV_FILE)
	set +a
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 \
	  -sha256 -days 365 -nodes \
	  -keyout "nginx/ssl/server.key" \
	  -out    "nginx/ssl/server.crt" \
	  -subj "/CN=$${DOMAIN}" \
	  -addext "subjectAltName=DNS:$${DOMAIN}"


# 	  pg_dump "postgres://postgres:Qwaszx_1@localhost:5432/shopdb" --schema-only --no-owner > internal/infra/db/schema/000_schema_dump.sql
