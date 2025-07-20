from flask import Flask, request, jsonify
import requests
from flask_cors import CORS

app = Flask(__name__)
CORS(app, resources={r"/send_order": {"origins": "*"}})

BOT_TOKEN = '7674430759:AAGwZzc184Jwazd7ZpxFEVnCOkfjOCy9BOM'
CHAT_ID = '436092326'


@app.route('/send_order', methods=['POST'])
def send_order():
    if request.method == 'OPTIONS':

        return '', 200
    data = request.get_json(silent=True) or {}
    message = data.get('message', 'Пустое сообщение')

    if not message:
        return jsonify({'error': 'No message provided'}), 400

    url = f'https://api.telegram.org/bot{BOT_TOKEN}/sendMessage'
    payload = {
        'chat_id': CHAT_ID,
        'text': message
    }

    response = requests.post(url, json=payload)

    if response.status_code == 200:
        return jsonify({'status': 'Message sent to Telegram'})
    else:
        return jsonify({'error': 'Failed to send message'}), 500


@app.route('/')
def index():
    return 'Python Telegram Bot is running!', 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
