import os
import base64
import json
import io
from flask import Flask, render_template, request, jsonify, redirect, url_for
from PIL import Image
import requests

app = Flask(__name__)

# Constants
SCRIPT_TYPES = ["auto", "hieroglyphic", "cuneiform", "greek", "latin", "runic"]
IMAGE_PROCESSING_ALGORITHMS = ["none", "rotate", "grayscale", "blur", "edge_detection", "sharpen", "threshold", "invert"]
API_URL = "http://localhost:8000"  # Using Go backend API

# Function to encode image to base64
def encode_image(image):
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    return base64.b64encode(buffered.getvalue()).decode()

# Routes
@app.route('/')
def index():
    return render_template('index.html', 
                          script_types=SCRIPT_TYPES,
                          image_algorithms=IMAGE_PROCESSING_ALGORITHMS)

@app.route('/translate-image', methods=['POST'])
def translate_image():
    if 'image' not in request.files:
        return jsonify({'error': 'No image provided'}), 400
    
    file = request.files['image']
    if file.filename == '':
        return jsonify({'error': 'No image selected'}), 400
    
    try:
        # Get form data
        script_type = request.form.get('script_type', 'auto')
        algorithm = request.form.get('algorithm', 'none')
        
        # Process image
        image = Image.open(file)
        encoded_image = encode_image(image)
        
        # Mock API response for now (since Go API might not be fully integrated)
        # In a real scenario, we would call the Go API
        # response = requests.post(f"{API_URL}/api/translate", json={
        #     "image": encoded_image,
        #     "scriptType": script_type,
        #     "imageProcessing": algorithm
        # })
        # return jsonify(response.json())
        
        # Mock response for demonstration
        return jsonify({
            'originalText': "𓀀𓁐𓂓𓃾𓆓 𓇋𓈖𓉐𓊖 𓏏𓐍𓑗𓌙𓍯 𓎛𓏲𓀀",
            'translatedText': "The pharaoh commands the building of a great monument to honor the gods of the Nile.",
            'detectedScript': script_type,
            'confidenceScore': 0.89,
            'processedAt': "2025-04-21T14:30:00Z",
            'metadata': {
                'timePeriod': {
                    'startYear': 1500,
                    'endYear': 1400,
                    'era': 'BCE',
                    'periodName': 'New Kingdom'
                },
                'geographicalOrigin': {
                    'region': 'Ancient Egypt',
                    'city': 'Thebes',
                    'site': 'Valley of the Kings'
                },
                'culturalContext': {
                    'civilization': 'Egyptian',
                    'period': 'New Kingdom',
                    'languageFamily': 'Afro-Asiatic'
                },
                'materialContext': {
                    'material': 'Papyrus',
                    'preservation': 'Well-preserved',
                    'creationTechnique': 'Hieroglyphic inscription'
                },
                'historicalEvents': [
                    {
                        'name': 'Reign of Thutmose III',
                        'date': '1479-1425 BCE',
                        'description': 'Thutmose III was a military leader who created the largest empire Egypt had ever seen.',
                        'significance': 'Major expansion of Egyptian influence'
                    }
                ]
            }
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/translate-text', methods=['POST'])
def translate_text():
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({'error': 'No text provided'}), 400
    
    try:
        text = data['text']
        script_type = data.get('script_type', 'auto')
        
        # Mock response for demonstration
        return jsonify({
            'translatedText': "The inscription describes a royal decree from the time of the Old Kingdom, mentioning tributes to the goddess Hathor.",
            'detectedScript': script_type,
            'confidenceScore': 0.92,
            'processedAt': "2025-04-21T14:35:00Z",
            'metadata': {
                'timePeriod': {
                    'startYear': 2700,
                    'endYear': 2200,
                    'era': 'BCE',
                    'periodName': 'Old Kingdom'
                },
                'geographicalOrigin': {
                    'region': 'Ancient Egypt',
                    'city': 'Memphis',
                    'site': 'Royal Palace'
                },
                'culturalContext': {
                    'civilization': 'Egyptian',
                    'period': 'Old Kingdom',
                    'languageFamily': 'Afro-Asiatic'
                }
            }
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/extract-metadata', methods=['POST'])
def extract_metadata():
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({'error': 'No text provided'}), 400
    
    try:
        text = data['text']
        script_type = data.get('script_type', 'auto')
        
        # Mock response for demonstration
        return jsonify({
            'metadata': {
                'timePeriod': {
                    'startYear': 700,
                    'endYear': 600,
                    'era': 'BCE',
                    'periodName': 'Neo-Babylonian Period'
                },
                'geographicalOrigin': {
                    'region': 'Mesopotamia',
                    'city': 'Babylon',
                    'site': 'Ishtar Gate'
                },
                'culturalContext': {
                    'civilization': 'Babylonian',
                    'period': 'Neo-Babylonian',
                    'languageFamily': 'Semitic'
                },
                'materialContext': {
                    'material': 'Clay Tablet',
                    'preservation': 'Partially damaged',
                    'creationTechnique': 'Cuneiform impression'
                },
                'historicalEvents': [
                    {
                        'name': 'Reign of Nebuchadnezzar II',
                        'date': '605-562 BCE',
                        'description': 'Nebuchadnezzar II rebuilt Babylon into one of the most splendid cities of the ancient world.',
                        'significance': 'Major architectural and cultural achievements'
                    }
                ]
            },
            'processedAt': "2025-04-21T14:40:00Z"
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/summarize', methods=['POST'])
def summarize():
    data = request.get_json()
    if not data or 'text' not in data:
        return jsonify({'error': 'No text provided'}), 400
    
    try:
        text = data['text']
        algorithm = data.get('algorithm', '')
        
        # Mock response for demonstration
        return jsonify({
            'summary': "This text appears to be a record of grain distribution from the royal granaries during the third year of Pharaoh Amenhotep III's reign. It details the amounts allocated to various temples and noble households.",
            'textLength': len(text),
            'processedAt': "2025-04-21T14:45:00Z"
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/health')
def health():
    return jsonify({'status': 'healthy', 'timestamp': "2025-04-21T14:50:00Z"})

# Handler for image preview
@app.route('/process-image-preview', methods=['POST'])
def process_image_preview():
    if 'image' not in request.files:
        return jsonify({'error': 'No image provided'}), 400
    
    file = request.files['image']
    algorithm = request.form.get('algorithm', 'none')
    
    try:
        # Process image (here, we're just returning the original as a placeholder)
        # In a real implementation, we would apply the selected algorithm
        image = Image.open(file)
        
        # Convert image to base64 for display
        encoded_image = encode_image(image)
        
        return jsonify({
            'processedImage': encoded_image
        })
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    # Ensure the templates directory exists
    os.makedirs('templates', exist_ok=True)
    
    # Start the Flask app
    app.run(host='0.0.0.0', port=5000)