SHELL := /bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

# ENV ?= dev
ENV_FILE := .env
-include $(ENV_FILE)

define EXPORT_BFF_DB_URI
[[ -f "$(DB_PASSWORD_BFF_FILE)" ]] || { echo "Missing password file: $(DB_PASSWORD_BFF_FILE)"; exit 1; }
DB_PASSWORD="$$(tr -d '\n' < "$(DB_PASSWORD_BFF_FILE)")"
export DB_URI="postgres://$(DB_BFF_USER):$${DB_PASSWORD}@$(DB_HOST_OUTSIDE):$(DB_BFF_PORT)/$(DB_BFF_NAME)?sslmode=$(DB_BFF_SSLMODE)"
endef

define EXPORT_AUTH_DB_URI
[[ -f "$(DB_PASSWORD_AUTH_FILE)" ]] || { echo "Missing password file: $(DB_PASSWORD_AUTH_FILE)"; exit 1; }
DB_PASSWORD="$$(tr -d '\n' < "$(DB_PASSWORD_AUTH_FILE)")"
export DB_URI="postgres://$(DB_AUTH_USER):$${DB_PASSWORD}@$(DB_HOST_OUTSIDE):$(DB_AUTH_PORT)/$(DB_AUTH_NAME)?sslmode=$(DB_AUTH_SSLMODE)"
endef


.PHONY: tidy build build-auth run sqlc lint test up down
.PHONY: migrate-create-bff migrate-up-bff migrate-down-bff
.PHONY: migrate-create-auth migrate-up-auth migrate-down-auth

migrate-create-bff:
	@test -n "$(name)" || (echo "Usage: make migrate-create-bff name=add_table"; exit 1)
	migrate create -ext sql -dir $(MIGR_BFF_DIR) -format 20060102150405 $(name)

migrate-up-bff:
	$(EXPORT_BFF_DB_URI)
	migrate -path $(MIGR_BFF_DIR) -database $${DB_URI} up

migrate-down-bff:
	$(EXPORT_BFF_DB_URI)
	migrate -path $(MIGR_BFF_DIR) -database $${DB_URI} down 1

migrate-force-bff:
	$(EXPORT_BFF_DB_URI)
	@test -n "$(version)" || (echo "Usage: make migrate-force-bff version=v№"; exit 1)
	migrate -path $(MIGR_BFF_DIR) -database $${DB_URI} force $(version)

migrate-create-auth:
	@test -n "$(name)" || (echo "Usage: make migrate-create-auth name=add_table"; exit 1)
	migrate create -ext sql -dir $(MIGR_AUTH_DIR) -format 20060102150405 $(name)

migrate-up-auth:
	$(EXPORT_AUTH_DB_URI)
	migrate -path $(MIGR_AUTH_DIR) -database $${DB_URI} up

migrate-down-auth:
	$(EXPORT_AUTH_DB_URI)
	migrate -path $(MIGR_AUTH_DIR) -database $${DB_URI} down 1

migrate-force-auth:
	$(EXPORT_AUTH_DB_URI)
	@test -n "$(version)" || (echo "Usage: make migrate-force-auth version=v№"; exit 1)
	migrate -path $(MIGR_AUTH_DIR) -database $${DB_URI} force $(version)



tidy:
	go work sync
	( cd contracts && go mod tidy )
	( cd services/bff && go mod tidy )
	( cd services/auth && go mod tidy )


build:
	go build -o bff ./services/bff/cmd/server

build-auth:
	go build -o auth ./services/auth/cmd/authserver

# run:
# 	DB_PASSWORD_FILE="./secrets/db_password.txt" \
# 	JWT_SECRET_FILE="./secrets/jwt_secret.txt" $(BIN)



sqlc: sqlc-vet-bff sqlc-vet-auth sqlc-gen-bff sqlc-gen-auth

sqlc-vet-bff:
	$(EXPORT_BFF_DB_URI)
	sqlc vet -f services/bff/sqlc.yml

sqlc-vet-auth:
	$(EXPORT_AUTH_DB_URI)
	sqlc vet -f services/auth/sqlc.yml

sqlc-gen-bff:
	$(EXPORT_BFF_DB_URI)
	sqlc generate -f services/bff/sqlc.yml

sqlc-gen-auth:
	$(EXPORT_AUTH_DB_URI)
	sqlc generate -f services/auth/sqlc.yml




lint:
	golangci-lint run

test:
	( cd services/bff && go test ./... )
	( cd services/auth && go test ./... )

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
	$(COMPOSE) logs bff auth migrator_bff migrator_auth nginx db_bff db_auth

logs-migrator:
	$(COMPOSE) logs  migrator

logs-bff:
	$(COMPOSE) logs  bff

logs-auth:
	$(COMPOSE) logs  auth

logs-nginx:
	$(COMPOSE) logs  nginx
	
logs-db:
	$(COMPOSE) logs  db_bff db_auth

graph_cp:
	$(COMPOSE) cp bff:/tmp/graph.dot ./graph.dot

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
