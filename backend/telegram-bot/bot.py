from flask import Flask, request, jsonify
import requests
from flask_cors import CORS
import os
from dotenv import load_dotenv
import logging

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

load_dotenv()

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})

BOT_TOKEN = os.getenv('TELEGRAM_TOKEN')
CHAT_ID = os.getenv('TELEGRAM_CHAT_ID')

logger.info(f"Bot starting with CHAT_ID: {CHAT_ID}")


@app.route('/send_order', methods=['POST', 'OPTIONS'])
def send_order():
    if request.method == 'OPTIONS':
        return '', 200

    try:
        data = request.get_json(silent=True) or {}
        message = data.get('message', 'Пустое сообщение')

        if not message:
            logger.warning("Empty message received")
            return jsonify({'error': 'No message provided'}), 400

        if not BOT_TOKEN or not CHAT_ID:
            logger.error("Missing TELEGRAM_TOKEN or TELEGRAM_CHAT_ID")
            return jsonify({'error': 'Bot configuration error'}), 500

        url = f'https://api.telegram.org/bot{BOT_TOKEN}/sendMessage'
        payload = {
            'chat_id': CHAT_ID,
            'text': message,
            'parse_mode': 'HTML'
        }

        logger.info(f"Sending message to Telegram: {message[:50]}...")
        response = requests.post(url, json=payload, timeout=10)

        if response.status_code == 200:
            logger.info("Message sent successfully")
            return jsonify({'status': 'Message sent to Telegram'}), 200
        else:
            logger.error(
                f"Telegram API error: {response.status_code} - {response.text}")
            return jsonify({'error': 'Failed to send message', 'details': response.text}), 500

    except Exception as e:
        logger.error(f"Exception in send_order: {str(e)}", exc_info=True)
        return jsonify({'error': 'Internal server error', 'details': str(e)}), 500


@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({
        'status': 'healthy',
        'service': 'telegram-bot',
        'bot_configured': bool(BOT_TOKEN and CHAT_ID)
    }), 200


@app.route('/', methods=['GET'])
def index():
    return jsonify({
        'service': 'Telegram Bot Service',
        'status': 'running',
        'endpoints': {
            '/send_order': 'POST - Send order to Telegram',
            '/health': 'GET - Health check'
        }
    }), 200


if __name__ == '__main__':
    # В production используется gunicorn, но для разработки можно запустить так
    app.run(host='0.0.0.0', port=5000, debug=False)
