from flask import Flask, jsonify
import os
import json
app = Flask(__name__)
APP_DIR=os.path.join(os.path.dirname(__file__), '../')
@app.route('/OSList', methods=['GET'])
async def OSList():
    with open(APP_DIR + 'OSList.json', 'r') as f:
        return jsonify(json.load(f))

if __name__ == '__main__':
    with open(APP_DIR + 'config.json', 'r') as f:
        config=json.load(f)
    app.run(debug=True,host=config['host'],port=config['port'])