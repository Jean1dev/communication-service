#FROM golang:1.22.7-alpine3.20 AS builder
FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]