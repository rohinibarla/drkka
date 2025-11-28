#!/usr/bin/env python3
"""
Simple Bottle server for dṛkka exam system
- Serves static files (HTML, JS, JSON, CSS)
- POST /submit endpoint that prints JSON to console
"""

from bottle import Bottle, request, response, static_file
import json
import os

app = Bottle()

# Get absolute path of current directory
STATIC_DIR = os.path.abspath(os.path.dirname(__file__))

# Submit endpoint - receives JSON and prints to console
@app.post('/submit')
def submit():
    try:
        # Get JSON data from request
        data = request.json

        if not data:
            response.status = 400
            return {'error': 'No JSON data received'}

        # Print to console with formatting
        print('\n' + '='*80)
        print('SUBMISSION RECEIVED')
        print('='*80)
        print(json.dumps(data, indent=2))
        print('='*80 + '\n')

        # Return success response
        return {'success': True, 'message': 'Submission received'}

    except Exception as e:
        print(f'ERROR: {str(e)}')
        response.status = 500
        return {'error': str(e)}

# Serve HTML files
@app.route('/')
def index():
    return static_file('index.html', root=STATIC_DIR)

@app.route('/index.html')
def index_html():
    return static_file('index.html', root=STATIC_DIR)

@app.route('/review.html')
def review_html():
    return static_file('review.html', root=STATIC_DIR)

# Serve JavaScript files
@app.route('/main.js')
def main_js():
    return static_file('main.js', root=STATIC_DIR)

@app.route('/review.js')
def review_js():
    return static_file('review.js', root=STATIC_DIR)

@app.route('/process_and_pack.js')
def process_and_pack_js():
    return static_file('process_and_pack.js', root=STATIC_DIR)

# Serve JSON files
@app.route('/questions.json')
def questions_json():
    return static_file('questions.json', root=STATIC_DIR)

@app.route('/sample_submission.json')
def sample_submission_json():
    return static_file('sample_submission.json', root=STATIC_DIR)

if __name__ == '__main__':
    print('Starting dṛkka server...')
    print('Access the exam at: http://localhost:5000')
    print('Access the review at: http://localhost:5000/review.html')
    print('Press Ctrl+C to stop the server\n')
    app.run(host='localhost', port=5000, debug=False)
