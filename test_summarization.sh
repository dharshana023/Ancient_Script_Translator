#!/bin/bash

# Test script to verify that different texts produce unique summaries

echo "Testing summarization functionality..."
echo "Sending test request 1..."

# First test text
curl -s -X POST -H "Content-Type: application/json" -d '{
    "text": "Ancient Greek mathematics was developed from the 7th century BC to the 4th century AD by Greek-speaking peoples. The word mathematics comes from the Greek word mathematike, meaning the art of learning. Greek mathematics constituted a major period in the history of mathematics, fundamental in respect of geometry and the idea of formal proof. Greek mathematics also influenced Hindu and Islamic mathematics.",
    "algorithm": "hybrid"
}' http://localhost:5000/api/summarize

echo -e "\nSending test request 2..."

# Second test text (different content)
curl -s -X POST -H "Content-Type: application/json" -d '{
    "text": "The Egyptian hieroglyphic script was one of the writing systems used by ancient Egyptians to represent their language. Because of their pictorial nature, hieroglyphs were difficult to write and required extensive training to read. Hieroglyphs stopped being used in everyday writing when the Coptic alphabet was introduced. The decipherment of Egyptian hieroglyphs was finally accomplished in the 1820s by Jean-François Champollion.",
    "algorithm": "hybrid"
}' http://localhost:5000/api/summarize

echo -e "\nThe summaries should be different if the fix is working correctly."