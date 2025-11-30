# 📋 Инструкции по развертыванию Telegram Bot

## ✅ Что было сделано

Telegram бот интегрирован в backend как отдельный микросервис на Python (Flask). Он принимает заказы от фронтенда и отправляет уведомления в Telegram.

### Структура файлов:

```
backend/telegram-bot/
├── Dockerfile              # Docker образ для бота
├── requirements.txt        # Python зависимости
├── bot.py                  # Основной код бота
├── .dockerignore          # Исключения для Docker
├── .env.example           # Пример конфигурации
└── README.md              # Документация бота
```

## 🚀 Быстрый старт

### 1. Настройка Telegram бота

1. **Создайте бота в Telegram:**
   - Откройте [@BotFather](https://t.me/BotFather)
   - Отправьте команду `/newbot`
   - Следуйте инструкциям
   - Сохраните полученный токен

2. **Получите Chat ID:**
   - Откройте [@userinfobot](https://t.me/userinfobot)
   - Отправьте любое сообщение
   - Скопируйте ваш ID

### 2. Настройка переменных окружения

Добавьте в корневой `.env` файл:

```bash
TELEGRAM_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=123456789
```

### 3. Запуск в production

```bash
# Запустить все сервисы включая telegram-bot
docker-compose -f docker-compose.prod.yaml up -d --build

# Проверить статус
docker-compose -f docker-compose.prod.yaml ps

# Проверить логи бота
docker-compose -f docker-compose.prod.yaml logs -f telegram-bot

# Проверить health check
curl http://localhost/telegram/health
```

### 4. Проверка работоспособности

```bash
# Отправить тестовое сообщение
curl -X POST http://localhost/telegram/send_order \
  -H "Content-Type: application/json" \
  -d '{"message": "Тестовое сообщение от сервера"}'
```

Сообщение должно прийти в указанный Telegram чат.

## 🔄 Обновление архитектуры

### Обновленная схема развертывания:

```
Клиент → NGINX (80/443)
           ├─→ /api/*       → Backend (Go API :8080)
           ├─→ /telegram/*  → Telegram Bot (Python :5000)
           └─→ /*           → Frontend (React :80)
                                ↓
                           PostgreSQL (:5432)
```

### Новые endpoints:

- `POST /telegram/send_order` - отправка заказа в Telegram
- `GET /telegram/health` - health check бота

## 📝 Использование на фронтенде

В вашем React приложении обновите код отправки заказа:

```javascript
const sendOrderToTelegram = async (orderData) => {
  const message = `
🛒 Новый заказ!

📋 Детали заказа:
${orderData.items.map(item => `- ${item.name}: ${item.quantity} шт.`).join('\n')}

💰 Сумма: ${orderData.total} руб.

👤 Клиент:
Имя: ${orderData.customerName}
Телефон: ${orderData.customerPhone}
Email: ${orderData.customerEmail}
  `.trim();

  try {
    const response = await fetch('/telegram/send_order', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message }),
    });

    if (response.ok) {
      console.log('✅ Заказ отправлен в Telegram');
      return true;
    } else {
      console.error('❌ Ошибка отправки в Telegram');
      return false;
    }
  } catch (error) {
    console.error('❌ Ошибка:', error);
    return false;
  }
};
```

## 🔧 Локальная разработка

Для разработки бота локально:

```bash
cd backend/telegram-bot

# Создать виртуальное окружение
python3 -m venv venv
source venv/bin/activate

# Установить зависимости
pip install -r requirements.txt

# Создать .env файл
cp .env.example .env
# Отредактировать .env и добавить реальные значения

# Запустить бот
python bot.py
```

Сервис будет доступен на `http://localhost:5000`

## 🐛 Troubleshooting

### Бот не отправляет сообщения

1. Проверьте правильность токена:
```bash
docker-compose -f docker-compose.prod.yaml logs telegram-bot | grep TELEGRAM
```

2. Проверьте, что бот запущен:
```bash
docker ps | grep telegram-bot
```

3. Проверьте health check:
```bash
curl http://localhost/telegram/health
```

### CORS ошибки

CORS настроен в NGINX. Проверьте конфигурацию:
- `/nginx/nginx.prod.conf` - для production
- `/backend/deployment/nginx/nginx.docker.conf` - для development

### Сообщения не приходят в Telegram

1. Убедитесь, что токен валидный:
```bash
curl https://api.telegram.org/bot<YOUR_TOKEN>/getMe
```

2. Проверьте Chat ID:
```bash
curl https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates
```

## 📚 Документация

- [ARCHITECTURE.md](../ARCHITECTURE.md) - Полная архитектура проекта
- [DEPLOYMENT.md](../DEPLOYMENT.md) - Подробное руководство по развертыванию
- [QUICKSTART.md](../QUICKSTART.md) - Быстрый старт
- [backend/telegram-bot/README.md](../backend/telegram-bot/README.md) - Документация бота

## 🔐 Безопасность

⚠️ **Важно:**
- Никогда не комитьте `.env` файл с реальными токенами
- Используйте разные токены для dev/prod окружений
- Регулярно обновляйте зависимости Python
- Храните токены в безопасном месте

## ⚙️ Переменные окружения

| Переменная | Описание | Пример |
|-----------|----------|--------|
| `TELEGRAM_TOKEN` | Токен бота от BotFather | `123456:ABC-DEF...` |
| `TELEGRAM_CHAT_ID` | ID чата для уведомлений | `123456789` |

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose -f docker-compose.prod.yaml logs -f telegram-bot`
2. Проверьте health: `curl http://localhost/telegram/health`
3. Создайте issue в репозитории с описанием проблемы

---

**Статус:** ✅ Готово к развертыванию
**Версия:** 1.0.0
**Дата:** $(date +%Y-%m-%d)
