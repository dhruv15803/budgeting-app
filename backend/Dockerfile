# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download \
	&& go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /api /app/api
COPY internal/database/migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

WORKDIR /app

ENV PORT=3000

ENTRYPOINT ["/app/docker-entrypoint.sh"]
