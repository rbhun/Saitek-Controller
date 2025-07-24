#!/bin/bash

# Test script for Saitek Controller GUI API
# This script demonstrates the API endpoints

BASE_URL="http://localhost:8080"

echo "Testing Saitek Controller GUI API..."
echo "====================================="

# Test 1: Get status
echo "1. Getting panel status..."
curl -s "$BASE_URL/api/status" | jq '.' 2>/dev/null || curl -s "$BASE_URL/api/status"
echo -e "\n"

# Test 2: Set radio display
echo "2. Setting radio display..."
curl -s -X POST "$BASE_URL/api/radio/set" \
  -H "Content-Type: application/json" \
  -d '{
    "com1Active": "118.25",
    "com1Standby": "118.75", 
    "com2Active": "121.50",
    "com2Standby": "121.90"
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL/api/radio/set" \
  -H "Content-Type: application/json" \
  -d '{"com1Active": "118.25", "com1Standby": "118.75", "com2Active": "121.50", "com2Standby": "121.90"}'
echo -e "\n"

# Test 3: Set multi panel display
echo "3. Setting multi panel display..."
curl -s -X POST "$BASE_URL/api/multi/set" \
  -H "Content-Type: application/json" \
  -d '{
    "topRow": "2500",
    "bottomRow": "3000",
    "leds": 3
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL/api/multi/set" \
  -H "Content-Type: application/json" \
  -d '{"topRow": "2500", "bottomRow": "3000", "leds": 3}'
echo -e "\n"

# Test 4: Set switch panel lights (gear down - green)
echo "4. Setting switch panel lights (gear down)..."
curl -s -X POST "$BASE_URL/api/switch/set" \
  -H "Content-Type: application/json" \
  -d '{
    "greenN": true,
    "greenL": true,
    "greenR": true,
    "redN": false,
    "redL": false,
    "redR": false
  }' | jq '.' 2>/dev/null || curl -s -X POST "$BASE_URL/api/switch/set" \
  -H "Content-Type: application/json" \
  -d '{"greenN": true, "greenL": true, "greenR": true, "redN": false, "redL": false, "redR": false}'
echo -e "\n"

# Test 5: Get updated status
echo "5. Getting updated status..."
curl -s "$BASE_URL/api/status" | jq '.' 2>/dev/null || curl -s "$BASE_URL/api/status"
echo -e "\n"

echo "API tests completed!"
echo "Open http://localhost:8080 in your browser to see the web interface." 