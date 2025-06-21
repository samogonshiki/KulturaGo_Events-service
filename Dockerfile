FROM golang:1.23 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/events ./cmd/events

FROM gcr.io/distroless/static
COPY --from=builder /out/events /events
ENTRYPOINT ["/events"]