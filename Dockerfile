FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-rate-limiter ./cmd/main.go

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/go-rate-limiter .

COPY cmd/.env .env

EXPOSE 8080

CMD ["./go-rate-limiter"]