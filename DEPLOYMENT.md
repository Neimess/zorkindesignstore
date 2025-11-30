# 🚀 Руководство по развертыванию Zorkin Design Store на одном сервере

Это руководство описывает, как развернуть backend (Go API) и frontend (React) на одном сервере с использованием Docker Compose и NGINX.

## 📋 Архитектура

```
Клиент → NGINX (порт 80/443)
           ├─→ /api/*     → Backend (Go API на порту 8080)
           └─→ /*         → Frontend (React SPA на порту 80)
                              ↓
                         PostgreSQL (порт 5432)
```

## 🔧 Предварительные требования

На сервере должны быть установлены:
- Docker (версия 20.10+)
- Docker Compose (версия 2.0+)
- Git (для клонирования репозитория)

### Установка Docker на Ubuntu/Debian:

```bash
# Обновить пакеты
sudo apt update

# Установить Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Добавить пользователя в группу docker
sudo usermod -aG docker $USER

# Установить Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Проверить установку
docker --version
docker-compose --version
```

## 📦 Развертывание

### 1. Клонировать репозиторий

```bash
cd /opt
sudo git clone https://github.com/Neimess/zorkindesignstore.git
cd zorkindesignstore
```

### 2. Создать .env файл с переменными окружения

```bash
cp .env.example .env
nano .env
```

Заполните переменные окружения:

```env
# PostgreSQL настройки
POSTGRES_USER=storeuser
POSTGRES_PASSWORD=STRONG_PASSWORD_HERE_12345
POSTGRES_DB=zorkinstore

# JWT и Admin настройки
JWT_SECRET=your_super_secret_jwt_key_min_32_chars_long
ADMIN_CODE=your_admin_secret_code_12345

# Frontend настройки
REACT_APP_API_URL=/api

# Домен вашего сервера (опционально)
DOMAIN=yourdomain.com
```

**⚠️ ВАЖНО:** 
- Используйте надежные пароли для production!
- JWT_SECRET должен быть минимум 32 символа
- Сохраните эти данные в безопасном месте

### 3. Собрать и запустить все сервисы

```bash
# Собрать образы и запустить
docker-compose -f docker-compose.prod.yaml up -d --build

# Проверить статус контейнеров
docker-compose -f docker-compose.prod.yaml ps

# Посмотреть логи
docker-compose -f docker-compose.prod.yaml logs -f
```

### 4. Проверить работоспособность

```bash
# Проверить health endpoint
curl http://localhost/health

# Проверить API
curl http://localhost/api/health

# Проверить frontend
curl http://localhost/
```

Откройте браузер и перейдите на `http://your-server-ip` или `http://yourdomain.com`

## 🔄 Обновление приложения

### Обновить код из репозитория

```bash
cd /opt/zorkindesignstore

# Получить последние изменения
git pull origin main

# Пересобрать и перезапустить контейнеры
docker-compose -f docker-compose.prod.yaml up -d --build

# Проверить логи
docker-compose -f docker-compose.prod.yaml logs -f backend frontend
```

### Обновить только backend

```bash
docker-compose -f docker-compose.prod.yaml up -d --build backend
```

### Обновить только frontend

```bash
docker-compose -f docker-compose.prod.yaml up -d --build frontend
```

## 🗄️ Управление базой данных

### Выполнить бэкап базы данных

```bash
# Создать бэкап
docker exec store-postgres pg_dump -U storeuser zorkinstore > backup_$(date +%Y%m%d_%H%M%S).sql

# Или с сжатием
docker exec store-postgres pg_dump -U storeuser zorkinstore | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz
```

### Восстановить из бэкапа

```bash
# Восстановить из SQL файла
cat backup.sql | docker exec -i store-postgres psql -U storeuser -d zorkinstore

# Восстановить из сжатого файла
gunzip -c backup.sql.gz | docker exec -i store-postgres psql -U storeuser -d zorkinstore
```

### Подключиться к базе данных

```bash
docker exec -it store-postgres psql -U storeuser -d zorkinstore
```

## 🔐 Настройка SSL/HTTPS (опционально)

### С использованием Let's Encrypt (Certbot)

1. Установить Certbot:

```bash
sudo apt install certbot python3-certbot-nginx
```

2. Остановить NGINX контейнер:

```bash
docker-compose -f docker-compose.prod.yaml stop nginx
```

3. Получить сертификат:

```bash
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com
```

4. Создать директорию для SSL сертификатов:

```bash
mkdir -p nginx/ssl
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/key.pem
```

5. Раскомментировать HTTPS блок в `nginx/nginx.prod.conf`:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # ... остальная конфигурация
}
```

6. Запустить NGINX:

```bash
docker-compose -f docker-compose.prod.yaml up -d nginx
```

### Автоматическое обновление SSL сертификатов

Добавьте в crontab:

```bash
sudo crontab -e

# Добавить строку для обновления каждый день в 3 утра
0 3 * * * certbot renew --quiet --post-hook "cd /opt/zorkindesignstore && docker-compose -f docker-compose.prod.yaml restart nginx"
```

## 📊 Мониторинг и логи

### Просмотреть логи всех сервисов

```bash
docker-compose -f docker-compose.prod.yaml logs -f
```

### Логи конкретного сервиса

```bash
# Backend
docker-compose -f docker-compose.prod.yaml logs -f backend

# Frontend
docker-compose -f docker-compose.prod.yaml logs -f frontend

# PostgreSQL
docker-compose -f docker-compose.prod.yaml logs -f postgres

# NGINX
docker-compose -f docker-compose.prod.yaml logs -f nginx
```

### Проверить использование ресурсов

```bash
docker stats
```

### Проверить здоровье контейнеров

```bash
docker-compose -f docker-compose.prod.yaml ps
```

## 🛠️ Полезные команды

### Остановить все сервисы

```bash
docker-compose -f docker-compose.prod.yaml down
```

### Остановить с удалением volumes (⚠️ удалит данные БД!)

```bash
docker-compose -f docker-compose.prod.yaml down -v
```

### Перезапустить сервис

```bash
docker-compose -f docker-compose.prod.yaml restart backend
docker-compose -f docker-compose.prod.yaml restart frontend
docker-compose -f docker-compose.prod.yaml restart nginx
```

### Выполнить команду внутри контейнера

```bash
# Зайти в shell backend
docker exec -it store-backend sh

# Зайти в shell frontend
docker exec -it store-frontend sh

# Зайти в PostgreSQL
docker exec -it store-postgres psql -U storeuser -d zorkinstore
```

### Очистить неиспользуемые Docker объекты

```bash
# Удалить неиспользуемые образы
docker image prune -a

# Удалить неиспользуемые volumes
docker volume prune

# Полная очистка (осторожно!)
docker system prune -a --volumes
```

## 🔧 Настройка конфигурации

### Backend конфигурация

Файл конфигурации: `backend/configs/prod.yaml`

Основные параметры:
- `http_server.address` - адрес и порт сервера
- `storage.*` - настройки подключения к БД (переопределяются через ENV)
- `jwt_config.*` - настройки JWT токенов
- `swagger.enable` - включить/выключить Swagger UI

### Frontend переменные окружения

При сборке можно переопределить API URL:

```bash
docker-compose -f docker-compose.prod.yaml build --build-arg REACT_APP_API_URL=/api frontend
```

## 📝 Структура переменных окружения

### Backend ENV переменные

| Переменная | Описание | Пример |
|-----------|----------|--------|
| `ENV` | Окружение (production/development) | `production` |
| `STORAGE_USER` | Пользователь PostgreSQL | `storeuser` |
| `STORAGE_PASSWORD` | Пароль PostgreSQL | `securepass123` |
| `STORAGE_DBNAME` | Имя базы данных | `zorkinstore` |
| `STORAGE_HOST` | Хост PostgreSQL | `postgres` |
| `STORAGE_PORT` | Порт PostgreSQL | `5432` |
| `JWT_SECRET` | Секретный ключ для JWT | `secret_key_min_32_chars` |
| `ADMIN_CODE` | Код для получения admin токена | `admin_code_123` |

### Frontend ENV переменные

| Переменная | Описание | Пример |
|-----------|----------|--------|
| `REACT_APP_API_URL` | URL бэкенд API | `/api` |

## 🔍 Troubleshooting

### Backend не запускается

1. Проверить логи:
```bash
docker-compose -f docker-compose.prod.yaml logs backend
```

2. Проверить подключение к БД:
```bash
docker exec -it store-backend nc -zv postgres 5432
```

3. Проверить миграции:
```bash
docker-compose -f docker-compose.prod.yaml logs migrate
```

### Frontend показывает ошибки API

1. Проверить NGINX конфигурацию:
```bash
docker exec -it store-nginx nginx -t
```

2. Проверить, что backend отвечает:
```bash
curl http://localhost/api/health
```

3. Проверить переменную окружения API_URL в браузере:
- Открыть DevTools → Console
- Проверить сетевые запросы (Network)

### База данных не доступна

1. Проверить, что контейнер запущен:
```bash
docker-compose -f docker-compose.prod.yaml ps postgres
```

2. Проверить health check:
```bash
docker inspect store-postgres | grep -A 10 Health
```

3. Подключиться напрямую:
```bash
docker exec -it store-postgres psql -U storeuser -d zorkinstore
```

## 🚦 Порты

| Сервис | Внутренний порт | Внешний порт | Описание |
|--------|----------------|--------------|----------|
| NGINX | 80, 443 | 80, 443 | Веб-сервер |
| Backend | 8080 | - | API (доступен через NGINX) |
| Frontend | 80 | - | React SPA (доступен через NGINX) |
| PostgreSQL | 5432 | - | База данных (внутренняя сеть) |

## 📚 Дополнительные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [NGINX Documentation](https://nginx.org/en/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## ⚠️ Безопасность

- ✅ Всегда используйте сильные пароли в production
- ✅ Регулярно обновляйте Docker образы
- ✅ Настройте firewall (UFW/iptables)
- ✅ Используйте SSL/TLS сертификаты
- ✅ Регулярно делайте бэкапы базы данных
- ✅ Ограничьте доступ к PostgreSQL (только внутренняя сеть)
- ✅ Мониторьте логи на подозрительную активность

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи всех сервисов
2. Убедитесь, что все переменные окружения заданы корректно
3. Проверьте сетевые подключения между контейнерами
4. Обратитесь к документации или создайте issue в репозитории
