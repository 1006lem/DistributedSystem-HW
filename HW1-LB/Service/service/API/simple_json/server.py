# server.py
# : Simple server code for API communication
# -----------------------------------------------

from flask import Flask, jsonify

app = Flask(__name__)

@app.route('/')
def home():
    data = {"message": "hello"}
    return jsonify(data)

if __name__ == "__main__":
    app.run('0.0.0.0', port=8089, debug=True)