FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chartsystem main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chartsystem /app/
COPY .env.example /app/.env

EXPOSE 8080

CMD ["./chartsystem"]
