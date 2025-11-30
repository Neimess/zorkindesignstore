# 🏗️ Архитектура развертывания

```
┌─────────────────────────────────────────────────────────────┐
│                         Интернет                             │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Сервер (VPS)  │
                    │  Port 80/443   │
                    └────────┬───────┘
                             │
        ┌────────────────────┼────────────────────────────┐
        │       Docker Network: store-net                 │
        │                                                  │
        │  ┌─────────────────────────────────────────┐    │
        │  │         NGINX Proxy                     │    │
        │  │      (store-nginx:80)                   │    │
        │  └──┬──────────────┬────────────┬──────────┘    │
        │     │              │            │                │
        │     │ /api/*       │ /telegram/* │ /*            │
        │     │              │            │                │
        │     ▼              ▼            ▼                │
        │  ┌─────────┐  ┌──────────┐ ┌──────────┐        │
        │  │ Backend │  │Telegram  │ │ Frontend │        │
        │  │  Go API │  │   Bot    │ │  React   │        │
        │  │  :8080  │  │ Python   │ │   :80    │        │
        │  └────┬────┘  │  :5000   │ └──────────┘        │
        │       │       └──────────┘                      │
        │       │ SQL queries                             │
        │       ▼                                         │
        │  ┌──────────┐                                  │
        │  │PostgreSQL│                                  │
        │  │   :5432  │                                  │
        │  └──────────┘                                  │
        │       ▲                                        │
        │       │                                        │
        │  ┌────┴─────┐                                 │
        │  │ Migrate  │ (выполняется при старте)        │
        │  └──────────┘                                 │
        └────────────────────────────────────────────────┘
```

## Потоки данных

### 1. Пользователь открывает сайт (/)
```
Браузер → NGINX:80 → Frontend Container → React SPA
```

### 2. Frontend делает API запрос (/api/*)
```
React App → NGINX:80 (/api/*) → Backend:8080 → PostgreSQL:5432
                                      ↓
                                  JWT Auth
                                      ↓
                                   Response
                                      ↓
                               Backend → NGINX → React
```

### 3. Frontend отправляет заказ в Telegram (/telegram/*)
```
React App → NGINX:80 (/telegram/send_order) → Telegram Bot:5000
                                                      ↓
                                              Telegram API
                                                      ↓
                                               Notification sent
                                                      ↓
                                            Response → NGINX → React
```

### 4. Первый запуск системы
```
1. Docker Compose запускает PostgreSQL
2. PostgreSQL health check проходит успешно
3. Migrate контейнер применяет миграции
4. После успешной миграции запускается Backend
5. Backend health check проходит
6. Запускается Frontend
7. Запускается Telegram Bot
8. Telegram Bot health check проходит
9. Запускается NGINX и начинает принимать запросы
```

## Переменные окружения

### Backend контейнер получает:
- `ENV=production` - режим работы
- `STORAGE_*` - параметры подключения к БД
- `JWT_SECRET` - для генерации токенов
- `ADMIN_CODE` - для админ-доступа

### Frontend контейнер собирается с:
- `REACT_APP_API_URL=/api` - относительный путь к API

### Telegram Bot контейнер получает:
- `TELEGRAM_TOKEN` - токен бота от BotFather
- `TELEGRAM_CHAT_ID` - ID чата для отправки уведомлений

### PostgreSQL получает:
- `POSTGRES_USER` - пользователь БД
- `POSTGRES_PASSWORD` - пароль
- `POSTGRES_DB` - имя базы данных

## Volumes (постоянное хранилище)

```
postgres_data:
  ├─ /var/lib/postgresql/data (внутри контейнера)
  └─ Хранит все данные PostgreSQL
```

## Сети

```
store-net (bridge):
  ├─ postgres:5432
  ├─ backend:8080
  ├─ frontend:80
  ├─ telegram-bot:5000
  └─ nginx:80,443
```

Только NGINX имеет маппинг портов на хост-систему.
Все остальные сервисы доступны только внутри Docker сети.

## Безопасность

✅ PostgreSQL не доступен извне (нет port mapping)
✅ Backend не доступен напрямую (только через NGINX)
✅ Frontend не доступен напрямую (только через NGINX)
✅ Telegram Bot не доступен напрямую (только через NGINX)
✅ Единственная точка входа - NGINX на портах 80/443
✅ Все секреты в .env файле (не в репозитории)
✅ TELEGRAM_TOKEN защищен переменными окружения
