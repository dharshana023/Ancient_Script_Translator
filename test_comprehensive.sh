#!/bin/bash

# Comprehensive test script for Ancient Script Decoder
# Tests multiple script types and content for metadata extraction and summarization

echo "===== ANCIENT SCRIPT DECODER - COMPREHENSIVE TEST ====="
echo "Testing metadata extraction and summarization for multiple historical texts"

# Test 1: Ancient Egyptian Hieroglyphic Text
echo -e "\n\n===== TEST 1: EGYPTIAN HIEROGLYPHIC TEXT ====="
EGYPTIAN_TEXT="The ancient hieroglyphic inscriptions found in the pyramids of Egypt reveal the complex religious beliefs of the pharaohs during the Old Kingdom (2700-2200 BCE). These texts, known as the Pyramid Texts, are among the oldest religious writings in the world, predating the more famous Book of the Dead. The texts were written during the reign of Unas, the last pharaoh of the Fifth Dynasty. They contain spells and instructions to help the deceased pharaoh navigate the afterlife and join the gods in the heavenly realm."

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
  http://localhost:5000/api/translate/text | jq '.metadata'

# Test 2: Mesopotamian Cuneiform Text
echo -e "\n\n===== TEST 2: MESOPOTAMIAN CUNEIFORM TEXT ====="
MESOPOTAMIAN_TEXT="The cuneiform tablet from ancient Sumer describes the legendary king Gilgamesh and his quest for immortality. Written during the Third Dynasty of Ur (circa 2100 BCE), this text forms part of the Epic of Gilgamesh, one of the earliest known works of literature. The narrative recounts Gilgamesh's journey to meet Utnapishtim, the survivor of the Great Flood, who was granted eternal life by the gods. The clay tablet was discovered near the ancient city of Uruk along the Euphrates river."

# Create JSON file for request
cat > /tmp/mesopotamian_request.json << EOL
{
  "originalText": "$MESOPOTAMIAN_TEXT",
  "scriptType": "cuneiform"
}
EOL

# Send request for Mesopotamian text
echo "Sending Mesopotamian manuscript test request..."
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d @/tmp/mesopotamian_request.json \
  http://localhost:5000/api/translate/text | jq '.metadata'

# Test 3: Ancient Greek Text
echo -e "\n\n===== TEST 3: ANCIENT GREEK TEXT ====="
GREEK_TEXT="This fragment of ancient Greek text recovered from a papyrus in Alexandria contains part of Aristotle's discourse on politics from the 4th century BCE. The philosopher discusses the ideal forms of government and the concept of citizenship in the Greek polis. He references the Athenian democracy and compares it to the political systems of Sparta and other city-states. The text also mentions the influence of Plato's teachings on Aristotle's political philosophy. This document was preserved in the library of Alexandria until its partial destruction."

# Create JSON file for request
cat > /tmp/greek_request.json << EOL
{
  "originalText": "$GREEK_TEXT",
  "scriptType": "greek"
}
EOL

# Send request for Greek text
echo "Sending Greek manuscript test request..."
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d @/tmp/greek_request.json \
  http://localhost:5000/api/translate/text | jq '.metadata'

# Test 4: Roman Latin Text
echo -e "\n\n===== TEST 4: ROMAN LATIN TEXT ====="
ROMAN_TEXT="The Latin inscription on this stone tablet commemorates the victory of Emperor Trajan over the Dacians in the early 2nd century CE. It was placed at the base of Trajan's Column in Rome, which was completed around 113 CE. The text describes the major battles of the Dacian Wars, including the siege of Sarmizegetusa and the final defeat of King Decebalus. The Roman Senate commissioned this monument to celebrate the conquest of Dacia and the expansion of the Empire to its greatest territorial extent under Trajan's rule."

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
  http://localhost:5000/api/translate/text | jq '.metadata'

# Test 5: Norse Runic Text
echo -e "\n\n===== TEST 5: NORSE RUNIC TEXT ====="
NORSE_TEXT="The runic inscription carved on this runestone tells of Viking expeditions to distant lands during the reign of Harald Bluetooth in the 10th century. Found in Denmark near the ancient settlement of Jelling, the stone bears witness to the Norse sailors who traveled to Vinland (North America) and Miklagard (Constantinople). The text mentions the worship of Thor and Odin, as well as descriptions of battles against Christian kingdoms to the south. The runes were carved by a master stoneworker named Thorvald, who signed his work at the bottom of the inscription."

# Create JSON file for request
cat > /tmp/norse_request.json << EOL
{
  "originalText": "$NORSE_TEXT",
  "scriptType": "runic"
}
EOL

# Send request for Norse text
echo "Sending Norse manuscript test request..."
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d @/tmp/norse_request.json \
  http://localhost:5000/api/translate/text | jq '.metadata'

# Cleanup
rm /tmp/egyptian_request.json /tmp/mesopotamian_request.json /tmp/greek_request.json /tmp/roman_request.json /tmp/norse_request.json

echo -e "\n\n===== COMPREHENSIVE TEST COMPLETE ====="