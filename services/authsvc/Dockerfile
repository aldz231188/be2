# syntax=docker/dockerfile:1.7

# --- Stage 1: Build ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# Устанавливаем git и сертификаты (иногда нужны для go mod)
RUN apk add --no-cache git ca-certificates

# Копируем go.mod обоих модулей (пути от корня проекта)
COPY contracts/go.mod ./contracts/
COPY internal/go.mod ./internal/

# Переходим в папку сервиса для загрузки зависимостей
WORKDIR /src/internal
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Копируем исходники (снова пути от корня)
WORKDIR /src
COPY contracts/ ./contracts/
COPY internal/ ./internal/

# Собираем
WORKDIR /src/internal

# COPY . .

# Сборка бинарника
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/server

# --- Stage 2: Final (Production) ---
# Используем distroless для безопасности и минимального размера (всего ~2MB + ваш бинарник)
FROM gcr.io/distroless/static-debian12 AS prod

# Настройки среды
ENV GIN_MODE=release

# Копируем бинарник из билдера
COPY --from=build --chown=nonroot:nonroot /out/app /app

# Используем не-root пользователя для безопасности
USER nonroot:nonroot

ENTRYPOINT ["/app"]

FROM alpine:3.20 AS dev
RUN apk add --no-cache ca-certificates bash curl busybox-extras
ENV GIN_MODE=debug
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]