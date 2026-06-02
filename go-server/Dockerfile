# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server /app/server

EXPOSE 8080

# SQLite test.db сохраняется в /app — монтируйте volume для персистентности
VOLUME ["/app"]

ENTRYPOINT ["/app/server"]
