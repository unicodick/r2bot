FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/bot /bot

ENTRYPOINT ["/bot"]
