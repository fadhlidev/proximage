FROM golang:1.26.2 AS build

RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/* \

WORKDIR /app

COPY . .

RUN go mod download

RUN go test ./...

RUN mkdir -p /dist \
    && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /dist/app .


# FINAL STAGE
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /dist/app .

RUN apt-get update && apt-get install -y --no-install-recommends \
 ca-certificates \
 dumb-init \
 && update-ca-certificates \
 && rm -rf /var/lib/apt/lists/*

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--", "./app"]