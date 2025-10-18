# syntax=docker/dockerfile:1.7

FROM golang:1.24 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/server

FROM gcr.io/distroless/static-debian12
# ENV APP_ENV=prod
COPY --from=build --chown=nonroot:nonroot /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]



