FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]