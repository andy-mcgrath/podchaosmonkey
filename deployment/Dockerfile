# syntax=docker/dockerfile:1

# Builder stage
FROM golang:1.18 AS builder

RUN useradd -u 1001 -m chaosuser

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chaos ./

# Service stage
FROM alpine AS service

COPY --from=builder /etc/passwd /etc/passwd

ENV ENVIRONMENT=production

WORKDIR /app
COPY --from=builder /app/chaos /app/chaos
COPY production.env ./production.env

USER 1001:1001

ENTRYPOINT ["/app/chaos"]