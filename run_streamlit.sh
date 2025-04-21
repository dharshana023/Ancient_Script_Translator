#!/bin/bash

# Run the minimal Streamlit app directly
streamlit run minimal_streamlit.py --server.port=5000 --server.address=0.0.0.0 --server.enableCORS=false --server.enableXsrfProtection=false --server.headless=true