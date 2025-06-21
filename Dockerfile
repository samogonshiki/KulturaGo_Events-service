FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w" \
    -o /bin/events-service ./cmd/events-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=builder /bin/events-service .

EXPOSE 8090        # HTTP-порт сервиса (REST/SSE/WebSocket)

ENTRYPOINT ["./events-service"]