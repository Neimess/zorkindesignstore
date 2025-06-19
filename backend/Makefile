PWD := $(shell pwd)

# Генерация OpenAPI-документации с помощью swag
doc:
	swag init \
		--generalInfo cmd/store/main.go \
		--output docs \
		--parseInternal \
		--parseDependency

# Выполнение миграций (в режиме разработки)
# Возможные варианты с ключами --up --down --force
# Выполнение миграций (можно указать действие: ACTION=--migrate или ACTION=--force=N)
devMigrate:
	go run ./cmd/migrate/. $(ACTION)

migrate-up:
	make devMigrate ACTION=--up

migrate-down:
	make devMigrate ACTION=--down