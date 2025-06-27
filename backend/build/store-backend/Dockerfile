FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o store-backend ./cmd/store

FROM alpine:3.19 AS runner
WORKDIR /app

COPY --from=builder /app/store-backend ./store-backend
COPY configs ./configs
RUN apk add curl
 
ENTRYPOINT ["./store-backend", "--config", "/app/config.yaml"]
