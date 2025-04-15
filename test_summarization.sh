#!/bin/bash

# Test script for verifying the enhanced summarization features:
# - Unique, content-aware summarization for different texts
# - Support for different summarization algorithms (extractive, hybrid)
# - Domain-specific templates based on content type (scientific, historical, literary)
# - Improved keyword extraction with TF-IDF-like scoring and concept diversity
# - Bigram support for more specific concept identification

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

echo -e "\nSending test request 3..."

# Third test text (different topic - literature)
curl -s -X POST -H "Content-Type: application/json" -d '{
    "text": "Shakespeare wrote 37 plays and 154 sonnets. His works include tragedies such as Hamlet, Macbeth, and Romeo and Juliet, as well as comedies like A Midsummer Nights Dream and Much Ado About Nothing. His plays continue to be performed and adapted worldwide, influencing literature, theater, and popular culture. Shakespeares use of language and his exploration of human emotions set a standard for storytelling.",
    "algorithm": "hybrid"
}' http://localhost:5000/api/summarize

echo -e "\nSending test request 4 (same text as #3 but with extractive algorithm)..."

# Fourth test - same as third but with extractive algorithm
curl -s -X POST -H "Content-Type: application/json" -d '{
    "text": "Shakespeare wrote 37 plays and 154 sonnets. His works include tragedies such as Hamlet, Macbeth, and Romeo and Juliet, as well as comedies like A Midsummer Nights Dream and Much Ado About Nothing. His plays continue to be performed and adapted worldwide, influencing literature, theater, and popular culture. Shakespeares use of language and his exploration of human emotions set a standard for storytelling.",
    "algorithm": "extractive"
}' http://localhost:5000/api/summarize

echo -e "\nThe summaries should be different for the same text with different algorithms."
echo -e "All summaries should be unique and content-specific."