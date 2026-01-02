# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS build
WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.work ./
COPY contracts ./contracts
COPY services/auth ./services/auth

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./services/auth/cmd/authserver

FROM gcr.io/distroless/static-debian12 AS prod

ENV GIN_MODE=release

COPY --from=build --chown=nonroot:nonroot /out/app /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]

FROM gcr.io/distroless/static-debian12 AS dev

ENV GIN_MODE=release

COPY --from=build --chown=nonroot:nonroot /out/app /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]
