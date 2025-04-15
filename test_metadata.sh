#!/bin/bash

# Test script for metadata extraction functionality
echo "Testing metadata extraction via direct text submission..."

# Test 1: Ancient Egyptian hieroglyphic text
EGYPTIAN_TEXT="The ancient hieroglyphic inscriptions found in the pyramids of Egypt reveal the complex religious beliefs of the pharaohs during the Old Kingdom (2700-2200 BCE). These texts, known as the Pyramid Texts, are among the oldest religious writings in the world, predating the more famous Book of the Dead. The texts were written during the reign of Unas, the last pharaoh of the Fifth Dynasty. They contain spells and instructions to help the deceased pharaoh navigate the afterlife."

# Create JSON file for request
cat > /tmp/egyptian_request.json << EOL
{
  "originalText": "$EGYPTIAN_TEXT",
  "scriptType": "hieroglyphic"
}
EOL

# Send request for Egyptian text
echo "Sending Egyptian manuscript test request..."
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d @/tmp/egyptian_request.json \
  http://localhost:5000/api/translate/text

echo -e "\n\nNow testing a different historical text (Roman)..."

# Test 2: Roman Latin text
ROMAN_TEXT="The Roman Senate during the time of Julius Caesar was a complex political institution that represented the patrician class. After Caesar crossed the Rubicon in 49 BCE, he effectively declared war on the Republic. The Senate initially opposed him, but after his victory in the civil war, many senators pledged loyalty to him. The Battle of Pharsalus marked a turning point where Caesar defeated Pompey's forces."

# Create JSON file for request
cat > /tmp/roman_request.json << EOL
{
  "originalText": "$ROMAN_TEXT",
  "scriptType": "latin"
}
EOL

# Send request for Roman text
echo "Sending Roman manuscript test request..."
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d @/tmp/roman_request.json \
  http://localhost:5000/api/translate/text

# Cleanup
rm /tmp/egyptian_request.json /tmp/roman_request.json

echo -e "\n\nMetadata extraction test complete!"