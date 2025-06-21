# Stage 1: build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ .

WORKDIR /app/cmd/migrate

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /migrate

# Stage 2: runtime
FROM alpine:3.19 AS migrate

WORKDIR /opt/migrate
COPY --from=builder /migrate .
COPY configs ./configs

ENTRYPOINT ["./migrate"]
