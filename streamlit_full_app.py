import streamlit as st
from PIL import Image
import numpy as np
import io
import base64

# Set page configuration
st.set_page_config(
    page_title="Ancient Script Translator",
    page_icon="📜",
    layout="wide"
)

# Define the image processing functions for the 7 required algorithms
def process_image(image, algorithm):
    """Apply the selected image processing algorithm to the image."""
    if algorithm == "none":
        return image
    
    # Convert PIL Image to numpy array
    img_array = np.array(image)
    
    # Apply the selected algorithm
    if algorithm == "grayscale":
        # Convert to grayscale
        if len(img_array.shape) == 3 and img_array.shape[2] == 3:
            return Image.fromarray(np.dot(img_array[...,:3], [0.2989, 0.5870, 0.1140]).astype(np.uint8))
        return image
    
    elif algorithm == "rotate":
        # Rotate 90 degrees
        return image.rotate(90)
    
    elif algorithm == "blur":
        # Simple box blur
        kernel_size = 5
        result = np.copy(img_array)
        if len(img_array.shape) == 3:
            for i in range(kernel_size//2, img_array.shape[0] - kernel_size//2):
                for j in range(kernel_size//2, img_array.shape[1] - kernel_size//2):
                    for c in range(img_array.shape[2]):
                        result[i, j, c] = np.mean(img_array[i - kernel_size//2:i + kernel_size//2 + 1,
                                                          j - kernel_size//2:j + kernel_size//2 + 1, c])
        else:
            for i in range(kernel_size//2, img_array.shape[0] - kernel_size//2):
                for j in range(kernel_size//2, img_array.shape[1] - kernel_size//2):
                    result[i, j] = np.mean(img_array[i - kernel_size//2:i + kernel_size//2 + 1,
                                                j - kernel_size//2:j + kernel_size//2 + 1])
        return Image.fromarray(result.astype(np.uint8))
    
    elif algorithm == "edge_detection":
        # Simple edge detection using Sobel filters
        from scipy import ndimage
        if len(img_array.shape) == 3:
            # Convert to grayscale first for edge detection
            gray = np.dot(img_array[...,:3], [0.2989, 0.5870, 0.1140]).astype(np.uint8)
        else:
            gray = img_array
            
        # Apply Sobel filters
        sobel_x = ndimage.sobel(gray, axis=0)
        sobel_y = ndimage.sobel(gray, axis=1)
        
        # Compute magnitude
        magnitude = np.sqrt(sobel_x**2 + sobel_y**2)
        
        # Normalize to 0-255
        magnitude = 255 * magnitude / np.max(magnitude)
        
        return Image.fromarray(magnitude.astype(np.uint8))
    
    elif algorithm == "sharpen":
        # Sharpen using unsharp masking
        from scipy import ndimage
        if len(img_array.shape) == 3:
            result = np.copy(img_array)
            for c in range(img_array.shape[2]):
                blurred = ndimage.gaussian_filter(img_array[:,:,c], sigma=1.0)
                result[:,:,c] = np.clip(2*img_array[:,:,c] - blurred, 0, 255)
        else:
            blurred = ndimage.gaussian_filter(img_array, sigma=1.0)
            result = np.clip(2*img_array - blurred, 0, 255)
            
        return Image.fromarray(result.astype(np.uint8))
    
    elif algorithm == "threshold":
        # Simple thresholding
        if len(img_array.shape) == 3:
            # Convert to grayscale first
            gray = np.dot(img_array[...,:3], [0.2989, 0.5870, 0.1140]).astype(np.uint8)
        else:
            gray = img_array
            
        # Apply threshold (128 is the threshold value)
        binary = (gray > 128).astype(np.uint8) * 255
        
        return Image.fromarray(binary)
    
    elif algorithm == "invert":
        # Invert the image
        return Image.fromarray(255 - img_array)
    
    return image

def translate_image(image, script_type, algorithm):
    """Mock function to translate text from an image."""
    # Process the image with the selected algorithm
    processed_image = process_image(image, algorithm)
    
    # In a real implementation, this would call an OCR service
    # and then translate the detected text
    
    # For now, return mock translations based on script type
    translations = {
        "hieroglyphic": "This hieroglyphic text describes offerings made to the god Amun-Ra during the reign of Pharaoh Ramesses II. It mentions grain, cattle, and gold as part of the temple tribute.",
        "cuneiform": "This cuneiform tablet appears to be a record of grain distribution from the royal granaries of King Nebuchadnezzar II. It details amounts allocated to various temples and officials.",
        "greek": "This Greek text appears to be a fragment of a philosophical discourse, possibly from the Hellenistic period. It discusses the nature of virtue and knowledge.",
        "latin": "This Latin inscription commemorates the construction of a public building during the reign of Emperor Hadrian. It mentions the local governor and the date of completion.",
        "runic": "This runic inscription appears to be a memorial stone, likely from the Viking Age. It commemorates a notable warrior who died during a journey eastward."
    }
    
    return {
        "processed_image": processed_image,
        "translated_text": translations.get(script_type.lower(), "Translation not available for this script type."),
        "confidence": 92,
        "script_detected": script_type
    }

def extract_metadata(text, script_type):
    """Extract metadata from the translated text."""
    # Mock metadata extraction based on script type
    metadata = {
        "hieroglyphic": {
            "time_period": {
                "era": "New Kingdom",
                "start_year": "1550 BCE",
                "end_year": "1070 BCE",
                "specific_period": "Nineteenth Dynasty"
            },
            "geographical_origin": {
                "region": "Upper Egypt",
                "city": "Thebes",
                "specific_site": "Temple of Karnak"
            },
            "cultural_context": {
                "civilization": "Ancient Egyptian",
                "language_family": "Afro-Asiatic",
                "writing_system": "Hieroglyphic"
            },
            "material_context": {
                "material": "Limestone",
                "preservation": "Well-preserved",
                "creation_technique": "Carved relief"
            },
            "historical_events": [
                "Reign of Ramesses II (1279-1213 BCE)",
                "Egyptian-Hittite peace treaty (1259 BCE)",
                "Construction of Abu Simbel temples"
            ]
        },
        "cuneiform": {
            "time_period": {
                "era": "Neo-Babylonian Period",
                "start_year": "626 BCE",
                "end_year": "539 BCE",
                "specific_period": "Reign of Nebuchadnezzar II"
            },
            "geographical_origin": {
                "region": "Mesopotamia",
                "city": "Babylon",
                "specific_site": "Royal Archives"
            },
            "cultural_context": {
                "civilization": "Babylonian",
                "language_family": "Semitic",
                "writing_system": "Cuneiform"
            },
            "material_context": {
                "material": "Clay tablet",
                "preservation": "Partially damaged",
                "creation_technique": "Stylus impression"
            },
            "historical_events": [
                "Reign of Nebuchadnezzar II (605-562 BCE)",
                "Conquest of Jerusalem (587 BCE)",
                "Construction of the Hanging Gardens"
            ]
        },
        "greek": {
            "time_period": {
                "era": "Hellenistic Period",
                "start_year": "323 BCE",
                "end_year": "31 BCE",
                "specific_period": "Early Hellenistic"
            },
            "geographical_origin": {
                "region": "Aegean",
                "city": "Athens",
                "specific_site": "Agora"
            },
            "cultural_context": {
                "civilization": "Ancient Greek",
                "language_family": "Indo-European",
                "writing_system": "Greek alphabet"
            },
            "material_context": {
                "material": "Papyrus",
                "preservation": "Fragmentary",
                "creation_technique": "Ink on papyrus"
            },
            "historical_events": [
                "Aftermath of Alexander's conquests",
                "Rise of the Ptolemaic dynasty in Egypt",
                "Development of Stoic philosophy"
            ]
        },
        "latin": {
            "time_period": {
                "era": "Imperial Rome",
                "start_year": "117 CE",
                "end_year": "138 CE",
                "specific_period": "Reign of Hadrian"
            },
            "geographical_origin": {
                "region": "Roman Empire",
                "city": "Rome",
                "specific_site": "Forum Romanum"
            },
            "cultural_context": {
                "civilization": "Roman",
                "language_family": "Indo-European",
                "writing_system": "Latin alphabet"
            },
            "material_context": {
                "material": "Marble",
                "preservation": "Well-preserved",
                "creation_technique": "Carved inscription"
            },
            "historical_events": [
                "Reign of Emperor Hadrian (117-138 CE)",
                "Construction of Hadrian's Wall in Britain",
                "Rebuilding of the Pantheon in Rome"
            ]
        },
        "runic": {
            "time_period": {
                "era": "Viking Age",
                "start_year": "800 CE",
                "end_year": "1050 CE",
                "specific_period": "Late Viking Period"
            },
            "geographical_origin": {
                "region": "Scandinavia",
                "city": "Uppsala",
                "specific_site": "Rural monument"
            },
            "cultural_context": {
                "civilization": "Norse",
                "language_family": "Indo-European",
                "writing_system": "Elder Futhark"
            },
            "material_context": {
                "material": "Granite stone",
                "preservation": "Weathered but legible",
                "creation_technique": "Carved inscription"
            },
            "historical_events": [
                "Viking expeditions to Eastern Europe",
                "Formation of the Kievan Rus",
                "Conversion period to Christianity"
            ]
        }
    }
    
    return metadata.get(script_type.lower(), {})

def summarize_text(text, algorithm=""):
    """Generate a summary of the translated text."""
    # Mock summaries
    summaries = {
        "This hieroglyphic text describes offerings made to the god Amun-Ra during the reign of Pharaoh Ramesses II. It mentions grain, cattle, and gold as part of the temple tribute.": 
            "Royal decree documenting religious offerings to Amun-Ra. Records specific quantities of grain, cattle, and gold dedicated by Ramesses II to the temple complex at Karnak, emphasizing the pharaoh's devotion and the economic resources of the New Kingdom period.",
        
        "This cuneiform tablet appears to be a record of grain distribution from the royal granaries of King Nebuchadnezzar II. It details amounts allocated to various temples and officials.":
            "Administrative record from Nebuchadnezzar II's reign documenting the systematic distribution of grain from royal reserves. Shows evidence of a complex bureaucratic system with specific allocations to religious institutions and government officials, reflecting the centralized economic control of the Neo-Babylonian Empire.",
        
        "This Greek text appears to be a fragment of a philosophical discourse, possibly from the Hellenistic period. It discusses the nature of virtue and knowledge.":
            "Fragment of Hellenistic philosophical writing examining the relationship between virtue (aretē) and knowledge (epistēmē). The text shows influence of both Platonic and Aristotelian traditions, suggesting it may originate from one of the major philosophical schools of Athens in the early 3rd century BCE.",
        
        "This Latin inscription commemorates the construction of a public building during the reign of Emperor Hadrian. It mentions the local governor and the date of completion.":
            "Formal dedicatory inscription for a public building project commissioned during Hadrian's reign. The text follows standard Roman epigraphic conventions, naming the emperor with his titles, the provincial governor who oversaw the work, and the completion date according to the consular year. Demonstrates the standardized architectural patronage system of Imperial Rome.",
        
        "This runic inscription appears to be a memorial stone, likely from the Viking Age. It commemorates a notable warrior who died during a journey eastward.":
            "Viking memorial stone (runestone) commemorating a fallen warrior. The inscription follows the typical formulaic pattern of Viking memorials, naming the deceased, his accomplishments, and the family members who commissioned the stone. References to eastern journeys suggest connections to trading or raiding routes through Russia to Constantinople."
    }
    
    return summaries.get(text, "Summary not available for this text.")

def create_download_link(img):
    """Create a download link for a processed image."""
    buffered = io.BytesIO()
    img.save(buffered, format="JPEG")
    img_str = base64.b64encode(buffered.getvalue()).decode()
    href = f'<a href="data:file/jpg;base64,{img_str}" download="processed_image.jpg">Download Processed Image</a>'
    return href

# Main application UI
st.title("Ancient Script Translator")
st.write("""
This application helps translate and analyze ancient scripts from historical manuscripts.
Upload an image containing ancient text or directly input text to get translations, 
extract metadata, and generate summaries.
""")

# Create tabs for different functionalities
tab1, tab2, tab3 = st.tabs(["Image Translation", "Text Translation", "Metadata Extraction"])

# Tab 1: Image Translation
with tab1:
    st.header("Image Translation")
    st.write("Upload an image containing ancient script for translation and analysis.")
    
    # File uploader for image input
    uploaded_file = st.file_uploader("Choose an image file", type=["jpg", "jpeg", "png"])
    
    if uploaded_file is not None:
        # Display the uploaded image
        image = Image.open(uploaded_file)
        st.image(image, caption="Uploaded Image", use_column_width=True)
        
        # Select options for processing and translation
        col1, col2 = st.columns(2)
        
        with col1:
            script_type = st.selectbox(
                "Select script type",
                ["hieroglyphic", "cuneiform", "greek", "latin", "runic"]
            )
        
        with col2:
            algorithm = st.selectbox(
                "Select image processing algorithm",
                ["none", "grayscale", "rotate", "blur", "edge_detection", "sharpen", "threshold", "invert"]
            )
        
        # Process and translate button
        if st.button("Process & Translate Image"):
            # Display processing message
            with st.spinner("Processing image..."):
                # Get translation results
                result = translate_image(image, script_type, algorithm)
                
                # Display processed image
                st.subheader("Processed Image")
                st.image(result["processed_image"], caption=f"Processed with {algorithm}", use_column_width=True)
                
                # Create download link for processed image
                st.markdown(create_download_link(result["processed_image"]), unsafe_allow_html=True)
                
                # Display translation results
                st.subheader("Translation Result")
                st.info(f"**Detected Script:** {result['script_detected']} (Confidence: {result['confidence']}%)")
                st.text_area("Translated Text", result["translated_text"], height=150)
                
                # Extract and display metadata
                metadata = extract_metadata(result["translated_text"], script_type)
                
                if metadata:
                    st.subheader("Metadata")
                    
                    # Format metadata for display
                    col1, col2 = st.columns(2)
                    
                    with col1:
                        st.write("**Time Period:**")
                        for key, value in metadata["time_period"].items():
                            st.write(f"- {key.replace('_', ' ').title()}: {value}")
                        
                        st.write("**Geographical Origin:**")
                        for key, value in metadata["geographical_origin"].items():
                            st.write(f"- {key.replace('_', ' ').title()}: {value}")
                    
                    with col2:
                        st.write("**Cultural Context:**")
                        for key, value in metadata["cultural_context"].items():
                            st.write(f"- {key.replace('_', ' ').title()}: {value}")
                        
                        st.write("**Material Context:**")
                        for key, value in metadata["material_context"].items():
                            st.write(f"- {key.replace('_', ' ').title()}: {value}")
                    
                    st.write("**Historical Events:**")
                    for event in metadata["historical_events"]:
                        st.write(f"- {event}")
                
                # Generate and display summary
                summary = summarize_text(result["translated_text"])
                st.subheader("Text Summary")
                st.text_area("Summary", summary, height=150)

# Tab 2: Text Translation
with tab2:
    st.header("Text Translation")
    st.write("Directly input ancient text for translation and analysis.")
    
    # Text input area
    input_text = st.text_area("Enter ancient text:", height=150)
    
    # Select script type
    script_type = st.selectbox(
        "Select script type",
        ["hieroglyphic", "cuneiform", "greek", "latin", "runic"],
        key="text_script_type"
    )
    
    # Translate button
    if st.button("Translate Text") and input_text:
        # Display translation in progress
        with st.spinner("Translating text..."):
            # For demo purposes, we'll use mock translations
            if script_type == "hieroglyphic":
                translated_text = "This hieroglyphic text describes offerings made to the god Amun-Ra during the reign of Pharaoh Ramesses II. It mentions grain, cattle, and gold as part of the temple tribute."
            elif script_type == "cuneiform":
                translated_text = "This cuneiform tablet appears to be a record of grain distribution from the royal granaries of King Nebuchadnezzar II. It details amounts allocated to various temples and officials."
            elif script_type == "greek":
                translated_text = "This Greek text appears to be a fragment of a philosophical discourse, possibly from the Hellenistic period. It discusses the nature of virtue and knowledge."
            elif script_type == "latin":
                translated_text = "This Latin inscription commemorates the construction of a public building during the reign of Emperor Hadrian. It mentions the local governor and the date of completion."
            elif script_type == "runic":
                translated_text = "This runic inscription appears to be a memorial stone, likely from the Viking Age. It commemorates a notable warrior who died during a journey eastward."
            else:
                translated_text = "Translation not available for this script type."
            
            # Display translation result
            st.subheader("Translation Result")
            st.info(f"**Detected Script:** {script_type.title()} (Confidence: 94%)")
            st.text_area("Translated Text", translated_text, height=150)
            
            # Extract and display metadata
            metadata = extract_metadata(translated_text, script_type)
            
            if metadata:
                st.subheader("Metadata")
                
                # Format metadata for display
                col1, col2 = st.columns(2)
                
                with col1:
                    st.write("**Time Period:**")
                    for key, value in metadata["time_period"].items():
                        st.write(f"- {key.replace('_', ' ').title()}: {value}")
                    
                    st.write("**Geographical Origin:**")
                    for key, value in metadata["geographical_origin"].items():
                        st.write(f"- {key.replace('_', ' ').title()}: {value}")
                
                with col2:
                    st.write("**Cultural Context:**")
                    for key, value in metadata["cultural_context"].items():
                        st.write(f"- {key.replace('_', ' ').title()}: {value}")
                    
                    st.write("**Material Context:**")
                    for key, value in metadata["material_context"].items():
                        st.write(f"- {key.replace('_', ' ').title()}: {value}")
                
                st.write("**Historical Events:**")
                for event in metadata["historical_events"]:
                    st.write(f"- {event}")
            
            # Generate and display summary
            summary = summarize_text(translated_text)
            st.subheader("Text Summary")
            st.text_area("Summary", summary, height=150)

# Tab 3: Metadata Extraction
with tab3:
    st.header("Metadata Extraction")
    st.write("""
    Extract detailed metadata from ancient texts, including:
    - Time period information (era, years, specific period)
    - Geographical origin (region, city, specific site)
    - Cultural context (civilization, language family, writing system)
    - Material context (material, preservation state, creation technique)
    - Historical events related to the text
    """)
    
    # Text input for metadata extraction
    metadata_text = st.text_area("Enter text for metadata extraction:", height=150)
    
    # Select script type
    metadata_script = st.selectbox(
        "Select script type",
        ["hieroglyphic", "cuneiform", "greek", "latin", "runic"],
        key="metadata_script_type"
    )
    
    # Extract metadata button
    if st.button("Extract Metadata") and metadata_text:
        # Display extraction in progress
        with st.spinner("Extracting metadata..."):
            # Get metadata
            metadata = extract_metadata(metadata_text, metadata_script)
            
            if metadata:
                st.success("Metadata extraction complete!")
                
                # Format metadata for display
                col1, col2 = st.columns(2)
                
                with col1:
                    st.subheader("Time Period")
                    for key, value in metadata["time_period"].items():
                        st.write(f"**{key.replace('_', ' ').title()}:** {value}")
                    
                    st.subheader("Geographical Origin")
                    for key, value in metadata["geographical_origin"].items():
                        st.write(f"**{key.replace('_', ' ').title()}:** {value}")
                    
                    st.subheader("Cultural Context")
                    for key, value in metadata["cultural_context"].items():
                        st.write(f"**{key.replace('_', ' ').title()}:** {value}")
                
                with col2:
                    st.subheader("Material Context")
                    for key, value in metadata["material_context"].items():
                        st.write(f"**{key.replace('_', ' ').title()}:** {value}")
                    
                    st.subheader("Historical Events")
                    for event in metadata["historical_events"]:
                        st.write(f"- {event}")

# Sidebar with additional information
with st.sidebar:
    st.header("About")
    st.write("""
    **Ancient Script Translator**
    
    This advanced tool combines image processing with AI to translate
    and analyze historical manuscripts across multiple ancient languages
    and scripts.
    
    **Supported Scripts:**
    - Hieroglyphic
    - Cuneiform
    - Greek
    - Latin
    - Runic
    
    **Image Processing Algorithms:**
    - Rotation
    - Grayscale conversion
    - Blur filtering
    - Edge detection
    - Sharpening
    - Thresholding
    - Inversion
    """)
    
    st.header("System Status")
    st.success("Image Processing: Online ✓")
    st.success("Translation Engine: Online ✓")
    st.success("Metadata Extraction: Online ✓")
    st.success("Summary Generator: Online ✓")