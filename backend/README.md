# Zorkin Design Store — Deployment & Development

Этот файл описывает **локальную разработку**, **продовую** (production) конфигурации и **миграцию базы** для проекта `Zorkin Design Store`.

---

## 🛠 Локальная разработка

Файл `docker-compose.dev.yaml` это среда для разработки с подключенной базой данных:

### Запуск

```bash
# В корне проекта
docker-compose -f docker-compose.dev.yml up --build
```

Сервис **Postgres** доступен на `localhost:5432`, **pgAdmin** на `http://localhost:8081`.

### Миграции

Перед стартом приложения выполните миграции:

```bash
# Перейдите в папку backend
cd backend
# Запустите миграции up
go run ./cmd/migrate/ -u
```

### Запуск приложения локально
После применения миграций можно запустить сервис без Docker:

```bash
# Перейти в папку backend
cd backend
go mod tidy
# Запуск с локальной конфигурацией
go run ./cmd/store/... --config=./configs/local.yaml
```
---

После этого swagger будет доступен на [localhost:8080/api/swagger](http://localhost:8080/api/swagger)
## 🚀 Продовая конфигурация

Это продовая конфигурация, с готовым билдом `docker-compose.yaml`:


### Переменные окружения

Создайте файл `.env` рядом с `docker-compose.yaml`:

```dotenv
POSTGRES_USER=produser
POSTGRES_PASSWORD=prodpassword
POSTGRES_DB=proddb
STORAGE_HOST=postgres
JWT_SECRET=secret
ADMIN_CODE=supersecret
```

### Запуск

```bash
docker-compose -f docker-compose.yaml up --build
```

---

## 📌 Важное

* **Миграции** в проде выполняются автоматически сервисом `migrate`.
* **NGINX** располагается на портах 80/443 и проксирует запросы к `backend`.
* Все сервисы объединены в сеть `store-net`.
