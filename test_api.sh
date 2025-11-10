#!/bin/bash

# Quick Test Script for Stocky API
# Run this after starting the application: go run main.go

echo "🧪 Stocky API Test Script"
echo "=========================="
echo ""

BASE_URL="http://localhost:8080/api/v1"

# Check if server is running
echo "📡 Checking if server is running..."
if ! curl -s "http://localhost:8080/api/v1/portfolio/test" > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    echo "Please start the server first: go run main.go"
    exit 1
fi
echo "✅ Server is running"
echo ""

# Test 1: Create a reward
echo "Test 1: Creating stock reward for user123..."
REWARD_RESPONSE=$(curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "idempotency_key": "test-reward-001"
  }')

echo "Response:"
echo "$REWARD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$REWARD_RESPONSE"
echo ""

# Test 2: Create another reward
echo "Test 2: Creating another reward (TCS)..."
curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "TCS",
    "quantity": 5.25,
    "idempotency_key": "test-reward-002"
  }' | python3 -m json.tool 2>/dev/null
echo ""

# Test 3: Test duplicate detection
echo "Test 3: Testing duplicate detection (should return 409)..."
DUPLICATE_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "idempotency_key": "test-reward-001"
  }')

echo "$DUPLICATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DUPLICATE_RESPONSE"
echo ""

# Test 4: Get today's stocks
echo "Test 4: Getting today's stocks for user123..."
curl -s "$BASE_URL/today-stocks/user123" | python3 -m json.tool 2>/dev/null
echo ""

# Test 5: Get portfolio
echo "Test 5: Getting portfolio for user123..."
curl -s "$BASE_URL/portfolio/user123" | python3 -m json.tool 2>/dev/null
echo ""

# Test 6: Get stats
echo "Test 6: Getting stats for user123..."
curl -s "$BASE_URL/stats/user123" | python3 -m json.tool 2>/dev/null
echo ""

echo "✨ Test Complete!"
echo ""
echo "Next steps:"
echo "  - Import Postman collection: Stocky.postman_collection.json"
echo "  - Check database: psql -U postgres -d assignment"
