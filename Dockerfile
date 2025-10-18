FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/todo-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-service /app/todo-service

COPY .env /app/.env

EXPOSE 8080

CMD ["/app/todo-service"]