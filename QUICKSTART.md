# 🚀 Быстрый старт развертывания

## Минимальные шаги для запуска на сервере

### 1. Предустановки
```bash
# Установить Docker и Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### 2. Клонировать и настроить
```bash
cd /opt
sudo git clone https://github.com/Neimess/zorkindesignstore.git
cd zorkindesignstore

# Создать .env файл
cp .env.example .env
nano .env
```

### 3. Заполнить .env (минимум)
```env
POSTGRES_USER=storeuser
POSTGRES_PASSWORD=YOUR_STRONG_PASSWORD
POSTGRES_DB=zorkinstore
JWT_SECRET=your_jwt_secret_min_32_chars_long
ADMIN_CODE=your_admin_code
REACT_APP_API_URL=/api
TELEGRAM_TOKEN=123456:ABC-DEF...  # От @BotFather
TELEGRAM_CHAT_ID=123456789         # От @userinfobot
```

### 4. Запустить
```bash
docker-compose -f docker-compose.prod.yaml up -d --build
```

### 5. Проверить
```bash
# Статус
docker-compose -f docker-compose.prod.yaml ps

# Логи
docker-compose -f docker-compose.prod.yaml logs -f

# Открыть в браузере
# http://your-server-ip
```

## Полезные команды

```bash
# Остановить
docker-compose -f docker-compose.prod.yaml down

# Перезапустить
docker-compose -f docker-compose.prod.yaml restart

# Обновить из git
git pull origin main
docker-compose -f docker-compose.prod.yaml up -d --build

# Бэкап БД
docker exec store-postgres pg_dump -U storeuser zorkinstore > backup.sql

# Логи конкретного сервиса
docker-compose -f docker-compose.prod.yaml logs -f backend
docker-compose -f docker-compose.prod.yaml logs -f frontend
docker-compose -f docker-compose.prod.yaml logs -f telegram-bot
```

## Структура портов

- **80** - HTTP (NGINX)
- **443** - HTTPS (NGINX, если настроен SSL)
- Все остальные сервисы в приватной сети

## Важно! 🔐

- ✅ Измените пароли в .env на надежные
- ✅ Настройте firewall (открыть только 80, 443, 22)
- ✅ Настройте SSL сертификаты для production
- ✅ Делайте регулярные бэкапы БД

Подробная документация: [DEPLOYMENT.md](./DEPLOYMENT.md)
