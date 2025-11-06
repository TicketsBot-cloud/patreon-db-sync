# Build container
FROM golang:1.25.3 AS builder

RUN go version

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates git zlib1g-dev

COPY . /go/src/github.com/TicketsBot/patreon-db-sync
WORKDIR /go/src/github.com/TicketsBot/patreon-db-sync

RUN git submodule update --init --recursive --remote

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -o main cmd/patreon-db-sync/main.go

# Prod container
FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates curl

COPY --from=builder /go/src/github.com/TicketsBot/patreon-db-sync/main /srv/patreon-db-sync/main

RUN chmod +x /srv/patreon-db-sync/main

RUN useradd -m container
USER container
WORKDIR /srv/patreon-db-sync

CMD ["/srv/patreon-db-sync/main"]
