from flask import Flask

app = Flask(__name__)

@app.route('/')
def index():
    return '''
    <!DOCTYPE html>
    <html>
    <head>
        <title>Ancient Script Translator</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                margin: 0;
                padding: 20px;
                background-color: #f5f5f5;
            }
            .container {
                max-width: 800px;
                margin: 0 auto;
                background-color: white;
                padding: 20px;
                border-radius: 8px;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
            h1 {
                color: #333;
                text-align: center;
            }
            .feature {
                margin-bottom: 15px;
                padding: 10px;
                background-color: #f9f9f9;
                border-left: 4px solid #007bff;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Ancient Script Translator</h1>
            <p>Welcome to the Ancient Script Translator tool. This application provides comprehensive analysis and translation of ancient texts.</p>
            
            <h2>Core Features:</h2>
            <div class="feature">
                <h3>Image Translation</h3>
                <p>Upload images containing ancient scripts for translation with 7 different image processing algorithms.</p>
            </div>
            
            <div class="feature">
                <h3>Text Translation</h3>
                <p>Input ancient text directly for translation with automatic script type detection.</p>
            </div>
            
            <div class="feature">
                <h3>Metadata Extraction</h3>
                <p>Extract comprehensive historical context including time periods, geographical origins, and cultural context.</p>
            </div>
            
            <div class="feature">
                <h3>Text Summarization</h3>
                <p>Generate concise summaries of ancient texts highlighting key information.</p>
            </div>
            
            <div class="feature">
                <h3>Multi-Script Support</h3>
                <p>Support for hieroglyphic, cuneiform, Greek, Latin, and runic scripts.</p>
            </div>
            
            <h2>Implementation Details:</h2>
            <ul>
                <li>Implemented TCP and UDP server components to demonstrate different transport protocols</li>
                <li>Created image processing utilities with 7 different algorithms</li>
                <li>Developed comprehensive metadata extraction service</li>
                <li>Built a user-friendly web interface</li>
            </ul>
        </div>
    </body>
    </html>
    '''

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)