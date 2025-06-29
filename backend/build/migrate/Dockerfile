# Stage 1: build
FROM golang:1.24-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev

RUN echo "Building for OS=$TARGETOS ARCH=$TARGETARCH VERSION=$VERSION"

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY cmd/migrate ./cmd/migrate
COPY migrations ./migrations

WORKDIR /app/cmd/migrate

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /migrate

# Stage 2: runtime
FROM alpine:3.19 AS migrate

WORKDIR /opt/migrate
COPY --from=builder /migrate .
COPY configs ./configs

ENTRYPOINT ["./migrate"]
