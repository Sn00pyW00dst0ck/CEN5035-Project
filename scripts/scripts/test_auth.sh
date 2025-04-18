```bash
#!/bin/bash

BASE_URL="http://localhost:3000/v1/api"
USERNAME="testuser"

# Get a challenge
echo "Getting challenge for user: $USERNAME"
CHALLENGE_RESP=$(curl -s "$BASE_URL/challenge?username=$USERNAME")
CHALLENGE=$(echo $CHALLENGE_RESP | jq -r .challenge)
echo "Received challenge: $CHALLENGE"

# Sign the challenge with private key (assuming OpenSSL is installed)
echo "Signing challenge with private key..."
echo -n "$CHALLENGE" > challenge.txt
SIGNATURE=$(openssl dgst -sha256 -sign private_key.pem -out signature.bin challenge.txt && base64 -w 0 signature.bin)
echo "Generated signature: $SIGNATURE"

# Login with the signature
echo "Logging in..."
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d "{\"username\":\"$USERNAME\",\"signature\":\"$SIGNATURE\"}")
TOKEN=$(echo $LOGIN_RESP | jq -r .token)
echo "Received token: $TOKEN"

# Test accessing a protected endpoint
echo "Testing protected endpoint..."
curl -s "$BASE_URL/account/" \
     -H "Authorization: Bearer $TOKEN" | jq .

# Clean up temporary files
rm challenge.txt signature.bin