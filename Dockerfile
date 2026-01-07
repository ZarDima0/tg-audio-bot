FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot ./cmd/bot

FROM alpine:latest

RUN apk --no-cache add ffmpeg

WORKDIR /root/

COPY --from=builder /app/bot .

CMD ["./bot"]