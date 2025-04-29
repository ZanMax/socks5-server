# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod init socks5-server || true
RUN go get github.com/armon/go-socks5
COPY *.go ./
COPY *.json ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /socks5-server

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /socks5-server .
COPY --from=builder /app/config.json .

EXPOSE 1080

CMD ["./socks5-server"]