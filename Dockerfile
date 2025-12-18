# syntax=docker/dockerfile:1.7

# --- Stage 1: Build ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# Устанавливаем git и сертификаты (иногда нужны для go mod)
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

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
FROM gcr.io/distroless/static-debian12 AS dev

# Настройки среды
ENV GIN_MODE=release

# Копируем бинарник из билдера
COPY --from=build --chown=nonroot:nonroot /out/app /app

# Используем не-root пользователя для безопасности
USER nonroot:nonroot

ENTRYPOINT ["/app"]