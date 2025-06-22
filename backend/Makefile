PWD := $(shell pwd)
APP_NAME = zorkin-store
VERSION  = $(shell git describe --tags --always --dirty)
COMMIT   = $(shell git rev-parse --short HEAD)
DATE     = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -s -w \
    -X 'main.version=$(VERSION)' \
    -X 'main.commit=$(COMMIT)' \
    -X 'main.buildDate=$(DATE)'

# Генерация OpenAPI-документации с помощью swag
doc:
	swag init \
		--generalInfo cmd/store/main.go \
		--output docs \
		--parseInternal \
		--parseDependency

build:
	go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME) ./cmd/main.go