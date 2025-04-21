import streamlit as st

st.title("Ancient Script Translator")
st.write("Welcome to the Ancient Script Translator!")

st.write("This application helps translate and analyze ancient scripts from historical manuscripts.")

st.header("Image Translation")
uploaded_file = st.file_uploader("Upload an image with ancient script", type=["jpg", "jpeg", "png"])

st.header("Text Translation")
input_text = st.text_area("Enter ancient text for translation:", height=100)

with st.sidebar:
    st.header("About")
    st.write("Ancient Script Translator v1.0")
    st.write("Supported scripts: Hieroglyphic, Cuneiform, Greek, Latin, Runic")