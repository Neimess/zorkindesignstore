FROM python:3.11-slim

WORKDIR /app

# Установка зависимостей
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копирование кода бота
COPY bot.py .

# Переменные окружения (будут переопределены через docker-compose)
ENV TELEGRAM_TOKEN=""
ENV TELEGRAM_CHAT_ID=""

# Открываем порт
EXPOSE 5000

# Запуск приложения
CMD ["python", "bot.py"]
