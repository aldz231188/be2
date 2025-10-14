# syntax=docker/dockerfile:1.7

# --- build stage ---
FROM golang:1.24.0 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/server

# --- run stage ---
FROM gcr.io/distroless/base-debian12
ENV APP_ENV=prod
COPY --from=build /out/app /app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
