FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o telegram-bot .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/telegram-bot .

RUN adduser -D -s /bin/sh appuser
USER appuser

CMD ["./telegram-bot"]