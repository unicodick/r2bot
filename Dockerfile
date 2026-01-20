FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o bot ./cmd/bot

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/bot /bot
ENTRYPOINT ["/bot"]
