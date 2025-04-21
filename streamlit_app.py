import streamlit as st
import requests
import json
import io
from PIL import Image
import numpy as np
import base64
from datetime import datetime

# Constants
API_URL = "http://localhost:5000"
SCRIPT_TYPES = ["auto", "hieroglyphic", "cuneiform", "greek", "latin", "runic"]
IMAGE_PROCESSING_ALGORITHMS = ["none", "rotate", "grayscale", "blur", "edge_detection", "sharpen", "threshold", "invert"]

# Page configuration
st.set_page_config(
    page_title="Ancient Script Translator",
    page_icon="📜",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Title and description
st.title("Ancient Script Translator")
st.markdown("""
This tool helps translate and analyze ancient scripts from historical manuscripts.
Upload an image containing ancient text or directly input text to get translations, 
extract metadata, and generate summaries.
""")

# Sidebar for configuration
st.sidebar.title("Settings")
selected_script_type = st.sidebar.selectbox("Script Type", SCRIPT_TYPES)
selected_image_algorithm = st.sidebar.selectbox("Image Processing", IMAGE_PROCESSING_ALGORITHMS)

# Function to encode image to base64
def encode_image(image):
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    return base64.b64encode(buffered.getvalue()).decode()

# Function to send image to API
def translate_image(image, script_type, algorithm):
    # Convert image to base64
    encoded_image = encode_image(image)
    
    # Prepare request data
    data = {
        "image": encoded_image,
        "scriptType": script_type,
        "imageProcessing": algorithm
    }
    
    # Make request to API
    try:
        response = requests.post(f"{API_URL}/api/translate", json=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error communicating with API: {e}")
        return None

# Function to translate text
def translate_text(text, script_type):
    # Prepare request data
    data = {
        "text": text,
        "scriptType": script_type
    }
    
    # Make request to API
    try:
        response = requests.post(f"{API_URL}/api/translate/text", json=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error communicating with API: {e}")
        return None

# Function to extract metadata
def extract_metadata(text, script_type):
    # Prepare request data
    data = {
        "text": text,
        "scriptType": script_type
    }
    
    # Make request to API
    try:
        response = requests.post(f"{API_URL}/api/metadata", json=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error communicating with API: {e}")
        return None

# Function to summarize text
def summarize_text(text, algorithm=""):
    # Prepare request data
    data = {
        "text": text,
        "algorithm": algorithm
    }
    
    # Make request to API
    try:
        response = requests.post(f"{API_URL}/api/summarize", json=data)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        st.error(f"Error communicating with API: {e}")
        return None

# Function to display metadata nicely
def display_metadata(metadata):
    if not metadata:
        return
    
    if "timePeriod" in metadata and metadata["timePeriod"]:
        st.subheader("Time Period")
        time_data = metadata["timePeriod"]
        cols = st.columns(2)
        with cols[0]:
            st.metric("Start Year", f"{time_data.get('startYear', 'Unknown')} {time_data.get('era', 'BCE/CE')}")
        with cols[1]:
            st.metric("End Year", f"{time_data.get('endYear', 'Unknown')} {time_data.get('era', 'BCE/CE')}")
        st.write(f"**Period Name**: {time_data.get('periodName', 'Unknown')}")
    
    if "geographicalOrigin" in metadata and metadata["geographicalOrigin"]:
        st.subheader("Geographical Origin")
        geo_data = metadata["geographicalOrigin"]
        cols = st.columns(3)
        with cols[0]:
            st.metric("Region", geo_data.get("region", "Unknown"))
        with cols[1]:
            st.metric("City", geo_data.get("city", "Unknown"))
        with cols[2]:
            st.metric("Site", geo_data.get("site", "Unknown"))
    
    if "culturalContext" in metadata and metadata["culturalContext"]:
        st.subheader("Cultural Context")
        cultural_data = metadata["culturalContext"]
        st.write(f"**Civilization**: {cultural_data.get('civilization', 'Unknown')}")
        st.write(f"**Cultural Period**: {cultural_data.get('period', 'Unknown')}")
        st.write(f"**Language Family**: {cultural_data.get('languageFamily', 'Unknown')}")
    
    if "materialContext" in metadata and metadata["materialContext"]:
        st.subheader("Material Context")
        material_data = metadata["materialContext"]
        cols = st.columns(2)
        with cols[0]:
            st.metric("Material", material_data.get("material", "Unknown"))
        with cols[1]:
            st.metric("Preservation", material_data.get("preservation", "Unknown"))
        st.write(f"**Creation Technique**: {material_data.get('creationTechnique', 'Unknown')}")
    
    if "historicalEvents" in metadata and metadata["historicalEvents"]:
        st.subheader("Historical Events")
        for event in metadata["historicalEvents"]:
            with st.expander(f"{event.get('name', 'Unknown Event')}"):
                st.write(f"**Date**: {event.get('date', 'Unknown')}")
                st.write(f"**Description**: {event.get('description', 'No description available')}")
                st.write(f"**Significance**: {event.get('significance', 'Unknown')}")

# Main application tabs
tab1, tab2, tab3 = st.tabs(["Image Translation", "Text Translation", "Metadata & Summary"])

# Image Translation Tab
with tab1:
    st.header("Upload an image with ancient script")
    uploaded_file = st.file_uploader("Choose an image...", type=["jpg", "jpeg", "png"])
    
    if uploaded_file is not None:
        # Display the original image
        original_image = Image.open(uploaded_file)
        st.image(original_image, caption="Uploaded Image", use_column_width=True)
        
        # Process and translate button
        if st.button("Process & Translate Image", key="translate_image_button"):
            with st.spinner("Processing image and extracting text..."):
                # Call API for translation
                translation_result = translate_image(original_image, selected_script_type, selected_image_algorithm)
                
                if translation_result:
                    # Display processed image if available
                    if "processedImage" in translation_result and translation_result["processedImage"]:
                        processed_image_data = base64.b64decode(translation_result["processedImage"])
                        processed_image = Image.open(io.BytesIO(processed_image_data))
                        st.image(processed_image, caption=f"Processed Image ({selected_image_algorithm})", use_column_width=True)
                    
                    # Display translation results
                    st.subheader("Translation Results")
                    st.write(f"**Detected Script**: {translation_result.get('detectedScript', 'Unknown')}")
                    st.write(f"**Confidence Score**: {translation_result.get('confidenceScore', 0):.2f}")
                    
                    # Original text and translation
                    col1, col2 = st.columns(2)
                    with col1:
                        st.text_area("Original Text", translation_result.get("originalText", ""), height=200)
                    with col2:
                        st.text_area("Translated Text", translation_result.get("translatedText", ""), height=200)
                    
                    # If metadata is available, display it
                    if "metadata" in translation_result and translation_result["metadata"]:
                        st.subheader("Metadata")
                        display_metadata(translation_result["metadata"])

# Text Translation Tab
with tab2:
    st.header("Translate Ancient Text")
    
    input_text = st.text_area("Enter ancient text for translation", height=200)
    
    if st.button("Translate Text", key="translate_text_button") and input_text:
        with st.spinner("Translating text..."):
            # Call API for translation
            translation_result = translate_text(input_text, selected_script_type)
            
            if translation_result:
                # Display translation results
                st.subheader("Translation Results")
                st.write(f"**Detected Script**: {translation_result.get('detectedScript', 'Unknown')}")
                st.write(f"**Confidence Score**: {translation_result.get('confidenceScore', 0):.2f}")
                
                # Translation
                st.text_area("Translated Text", translation_result.get("translatedText", ""), height=200)
                
                # If metadata is available, display it
                if "metadata" in translation_result and translation_result["metadata"]:
                    st.subheader("Metadata")
                    display_metadata(translation_result["metadata"])

# Metadata & Summary Tab
with tab3:
    st.header("Extract Metadata & Generate Summary")
    
    input_text = st.text_area("Enter ancient text for analysis", height=200)
    
    col1, col2 = st.columns(2)
    with col1:
        if st.button("Extract Metadata", key="extract_metadata_button") and input_text:
            with st.spinner("Extracting metadata..."):
                # Call API for metadata extraction
                metadata_result = extract_metadata(input_text, selected_script_type)
                
                if metadata_result and "metadata" in metadata_result:
                    # Display metadata
                    st.subheader("Metadata")
                    display_metadata(metadata_result["metadata"])
    
    with col2:
        if st.button("Generate Summary", key="generate_summary_button") and input_text:
            with st.spinner("Generating summary..."):
                # Call API for summarization
                summary_result = summarize_text(input_text)
                
                if summary_result and "summary" in summary_result:
                    # Display summary
                    st.subheader("Summary")
                    st.write(summary_result["summary"])
                    st.write(f"**Original Text Length**: {summary_result.get('textLength', 0)} characters")
                    st.write(f"**Processed At**: {summary_result.get('processedAt', datetime.now().isoformat())}")

# Check API health
try:
    response = requests.get(f"{API_URL}/api/health")
    if response.status_code == 200:
        st.sidebar.success("API is healthy ✓")
    else:
        st.sidebar.error("API health check failed")
except requests.exceptions.RequestException:
    st.sidebar.error("Cannot connect to API")

# Information about the application
st.sidebar.markdown("---")
st.sidebar.info("""
**About the Ancient Script Translator**

This tool combines advanced image processing and AI to translate 
and analyze historical manuscripts across multiple ancient languages
and scripts including hieroglyphic, cuneiform, Greek, Latin, and runic.

Features:
- Image-based translation
- Text-based translation
- Comprehensive metadata extraction
- Contextual summarization
- Image processing tools
""")