# 🔄 Миграция с PythonBot.py на новый Telegram Bot сервис

## Что изменилось

**Старая реализация:**
- Файл: `frontend/src/services/PythonBot.py`
- Запускался отдельно, не в Docker
- Требовал ручной установки зависимостей

**Новая реализация:**
- Путь: `backend/telegram-bot/`
- Запускается в Docker контейнере
- Интегрирован с основной архитектурой
- Автоматический запуск через docker-compose
- Проксируется через NGINX

## Изменения в URL

### Старый endpoint:
```
http://localhost:5000/send_order
```

### Новый endpoint:
```
http://your-domain.com/telegram/send_order
или
http://localhost/telegram/send_order
```

## Обновление кода фронтенда

### Было (старый PythonBot.py):

```javascript
const response = await fetch('http://localhost:5000/send_order', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({ message }),
});
```

### Стало (новый telegram-bot):

```javascript
const response = await fetch('/telegram/send_order', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({ message }),
});
```

**Изменения:**
- URL изменен с `http://localhost:5000/send_order` на `/telegram/send_order`
- Теперь используется относительный путь (через NGINX proxy)
- CORS настроен автоматически через NGINX

## Шаги миграции

### 1. Удалите старый файл (опционально)

```bash
# Старый файл больше не нужен
rm frontend/src/services/PythonBot.py
```

### 2. Обновите код фронтенда

Найдите все места в коде, где используется:
- `http://localhost:5000/send_order`
- `PythonBot.py`
- Импорты Python кода из фронтенда

Замените на новый endpoint: `/telegram/send_order`

### 3. Обновите переменные окружения

Добавьте в `.env`:
```bash
TELEGRAM_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

### 4. Перезапустите сервисы

```bash
docker-compose -f docker-compose.prod.yaml up -d --build
```

## Преимущества новой архитектуры

✅ **Интеграция с Docker** - все в одном стеке
✅ **Автоматический запуск** - не требует ручного управления
✅ **Единая точка входа** - через NGINX
✅ **CORS из коробки** - настроен в NGINX
✅ **Health checks** - мониторинг работоспособности
✅ **Логирование** - централизованные логи Docker
✅ **Масштабируемость** - легко увеличить количество workers
✅ **Production-ready** - gunicorn вместо Flask dev server

## Проверка миграции

1. **Проверьте, что старый сервис не запущен:**
```bash
ps aux | grep PythonBot
# Должно быть пусто
```

2. **Проверьте новый сервис:**
```bash
docker ps | grep telegram-bot
# Должен показать запущенный контейнер

curl http://localhost/telegram/health
# Должен вернуть JSON со статусом
```

3. **Отправьте тестовое сообщение:**
```bash
curl -X POST http://localhost/telegram/send_order \
  -H "Content-Type: application/json" \
  -d '{"message": "Тестовое сообщение после миграции"}'
```

## Откат (если нужно)

Если нужно вернуться к старой версии:

1. Остановите Docker контейнеры
2. Верните старый `PythonBot.py` файл
3. Запустите его вручную:
```bash
cd frontend/src/services
python3 PythonBot.py
```

Но мы **не рекомендуем** откат, так как новая архитектура более надежная и production-ready.

## Дополнительные возможности

Новый telegram-bot поддерживает:

- **HTML форматирование** в сообщениях
- **Health check endpoint** для мониторинга
- **Расширенное логирование**
- **Автоматический перезапуск** при ошибках
- **Масштабирование** через увеличение workers

Пример использования HTML:

```javascript
const message = `
<b>🛒 Новый заказ!</b>

<b>📋 Детали заказа:</b>
- Товар 1: 2 шт.
- Товар 2: 1 шт.

<b>💰 Сумма:</b> 5000 руб.

<b>👤 Клиент:</b>
Имя: Иван Иванов
Телефон: +7 123 456-78-90
`;

await fetch('/telegram/send_order', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ message }),
});
```

---

**Статус:** ✅ Готово к использованию
**Совместимость:** Обратная совместимость с API
