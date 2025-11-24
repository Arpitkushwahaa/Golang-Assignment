# API Testing Guide

This guide provides comprehensive examples for testing all API endpoints.

## Setup

1. Ensure the server is running:
   ```bash
   go run main.go
   ```

2. Base URL: `http://localhost:8080/api`

## Using cURL

### 1. Health Check

```bash
curl http://localhost:8080/api/health
```

**Expected Response:**
```json
{
  "status": "healthy",
  "time": "2025-01-23T10:30:00Z"
}
```

---

### 2. Create Reward

#### Example 1: RELIANCE Stock
```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": 2.5,
    "timestamp": "2025-01-23T10:30:00Z"
  }'
```

#### Example 2: TCS Stock
```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "TCS",
    "quantity": 1.5,
    "timestamp": "2025-01-23T14:15:00Z"
  }'
```

#### Example 3: INFY Stock
```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "INFY",
    "quantity": 3.25,
    "timestamp": "2025-01-23T16:45:00Z"
  }'
```

**Expected Response:**
```json
{
  "success": true
}
```

---

### 3. Get Today's Stocks

```bash
curl http://localhost:8080/api/today-stocks/1
```

**Expected Response:**
```json
{
  "userId": 1,
  "rewards": [
    {
      "symbol": "INFY",
      "quantity": "3.25",
      "timestamp": "2025-01-23T16:45:00Z"
    },
    {
      "symbol": "TCS",
      "quantity": "1.5",
      "timestamp": "2025-01-23T14:15:00Z"
    },
    {
      "symbol": "RELIANCE",
      "quantity": "2.5",
      "timestamp": "2025-01-23T10:30:00Z"
    }
  ]
}
```

---

### 4. Get Historical INR

```bash
curl http://localhost:8080/api/historical-inr/1
```

**Expected Response:**
```json
{
  "userId": 1,
  "days": [
    {
      "date": "2025-01-22",
      "valueINR": "15234.5600"
    },
    {
      "date": "2025-01-21",
      "valueINR": "15102.3400"
    }
  ]
}
```

---

### 5. Get Stats

```bash
curl http://localhost:8080/api/stats/1
```

**Expected Response:**
```json
{
  "userId": 1,
  "todayRewards": {
    "INFY": "3.25",
    "RELIANCE": "2.5",
    "TCS": "1.5"
  },
  "portfolioValueINR": "18456.7800"
}
```

---

### 6. Get Portfolio

```bash
curl http://localhost:8080/api/portfolio/1
```

**Expected Response:**
```json
{
  "userId": 1,
  "holdings": [
    {
      "symbol": "RELIANCE",
      "quantity": "2.5",
      "currentPrice": "2450.5000",
      "currentValue": "6126.2500"
    },
    {
      "symbol": "TCS",
      "quantity": "1.5",
      "currentPrice": "3680.7500",
      "currentValue": "5521.1250"
    },
    {
      "symbol": "INFY",
      "quantity": "3.25",
      "currentPrice": "1520.3000",
      "currentValue": "4940.9750"
    }
  ],
  "totalValue": "16588.3500"
}
```

---

## Edge Case Testing

### Test 1: Duplicate Reward (Should Fail)

```bash
# First reward
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": 2.5,
    "timestamp": "2025-01-23T10:30:00Z"
  }'

# Duplicate (Should return 409 Conflict)
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": 2.5,
    "timestamp": "2025-01-23T10:30:00Z"
  }'
```

**Expected Response:**
```json
{
  "success": false,
  "error": "duplicate reward: identical reward already exists"
}
```

---

### Test 2: Invalid Quantity (Should Fail)

```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": -5,
    "timestamp": "2025-01-23T10:30:00Z"
  }'
```

**Expected Response:**
```json
{
  "success": false,
  "error": "Invalid request payload: ..."
}
```

---

### Test 3: Invalid Timestamp (Should Fail)

```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": 2.5,
    "timestamp": "invalid-timestamp"
  }'
```

**Expected Response:**
```json
{
  "success": false,
  "error": "Invalid timestamp format, use RFC3339"
}
```

---

### Test 4: Missing Required Fields (Should Fail)

```bash
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "quantity": 2.5
  }'
```

**Expected Response:**
```json
{
  "success": false,
  "error": "Invalid request payload: ..."
}
```

---

## Using Postman

### Import Collection

1. Open Postman
2. Click **Import**
3. Select `Stocky_Postman_Collection.json`
4. Collection will be loaded with all endpoints

### Run Collection

1. Select **Stocky Backend API** collection
2. Click **Run** to execute all requests
3. View results in the Collection Runner

### Environment Variables

Set `base_url` in Postman environment:
- Variable: `base_url`
- Value: `http://localhost:8080`

---

## Testing Workflow

### Complete User Flow

```bash
# Step 1: Create first reward
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "RELIANCE", "quantity": 2.5, "timestamp": "2025-01-23T10:30:00Z"}'

# Step 2: Create second reward
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "TCS", "quantity": 1.5, "timestamp": "2025-01-23T14:15:00Z"}'

# Step 3: Check today's rewards
curl http://localhost:8080/api/today-stocks/1

# Step 4: Check stats
curl http://localhost:8080/api/stats/1

# Step 5: Check portfolio
curl http://localhost:8080/api/portfolio/1
```

---

## Multi-User Testing

```bash
# User 1 - RELIANCE
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "RELIANCE", "quantity": 2.5, "timestamp": "2025-01-23T10:30:00Z"}'

# User 2 - HDFCBANK
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 2, "symbol": "HDFCBANK", "quantity": 5.0, "timestamp": "2025-01-23T11:00:00Z"}'

# Check User 1 portfolio
curl http://localhost:8080/api/portfolio/1

# Check User 2 portfolio
curl http://localhost:8080/api/portfolio/2
```

---

## Automated Testing Script

Create `test_api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api"

echo "=== Testing Stocky Backend API ==="

echo "\n1. Health Check"
curl -s $BASE_URL/health | jq

echo "\n2. Create Reward - RELIANCE"
curl -s -X POST $BASE_URL/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "RELIANCE", "quantity": 2.5, "timestamp": "2025-01-23T10:30:00Z"}' | jq

echo "\n3. Create Reward - TCS"
curl -s -X POST $BASE_URL/reward \
  -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "TCS", "quantity": 1.5, "timestamp": "2025-01-23T14:15:00Z"}' | jq

echo "\n4. Get Today's Stocks"
curl -s $BASE_URL/today-stocks/1 | jq

echo "\n5. Get Stats"
curl -s $BASE_URL/stats/1 | jq

echo "\n6. Get Portfolio"
curl -s $BASE_URL/portfolio/1 | jq

echo "\n=== Tests Complete ==="
```

Make it executable and run:
```bash
chmod +x test_api.sh
./test_api.sh
```

---

## PowerShell Testing Script (Windows)

Create `test_api.ps1`:

```powershell
$BASE_URL = "http://localhost:8080/api"

Write-Host "=== Testing Stocky Backend API ===" -ForegroundColor Green

Write-Host "`n1. Health Check" -ForegroundColor Yellow
Invoke-RestMethod -Uri "$BASE_URL/health" -Method Get | ConvertTo-Json

Write-Host "`n2. Create Reward - RELIANCE" -ForegroundColor Yellow
$body1 = @{
    userId = 1
    symbol = "RELIANCE"
    quantity = 2.5
    timestamp = "2025-01-23T10:30:00Z"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$BASE_URL/reward" -Method Post -Body $body1 -ContentType "application/json" | ConvertTo-Json

Write-Host "`n3. Get Today's Stocks" -ForegroundColor Yellow
Invoke-RestMethod -Uri "$BASE_URL/today-stocks/1" -Method Get | ConvertTo-Json

Write-Host "`n4. Get Stats" -ForegroundColor Yellow
Invoke-RestMethod -Uri "$BASE_URL/stats/1" -Method Get | ConvertTo-Json

Write-Host "`n5. Get Portfolio" -ForegroundColor Yellow
Invoke-RestMethod -Uri "$BASE_URL/portfolio/1" -Method Get | ConvertTo-Json

Write-Host "`n=== Tests Complete ===" -ForegroundColor Green
```

Run:
```powershell
.\test_api.ps1
```

---

## Monitoring Logs

Watch application logs in real-time:

```bash
# Run with verbose logging
GIN_MODE=debug go run main.go

# Or tail log file if logging to file
tail -f stocky.log
```

---

## Performance Testing

### Using Apache Bench

```bash
# Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/api/health

# Test portfolio endpoint
ab -n 100 -c 5 http://localhost:8080/api/portfolio/1
```

### Using wrk

```bash
# Install wrk
brew install wrk  # macOS
sudo apt install wrk  # Linux

# Run load test
wrk -t4 -c100 -d30s http://localhost:8080/api/health
```

---

## Troubleshooting

### Connection Refused
```
curl: (7) Failed to connect to localhost port 8080: Connection refused
```
**Solution:** Ensure the server is running with `go run main.go`

### Invalid JSON
```
{"success":false,"error":"Invalid request payload: ..."}
```
**Solution:** Check JSON syntax and required fields

### Database Error
```
{"success":false,"error":"Failed to create reward: ..."}
```
**Solution:** Check database connection and logs

---

**Happy Testing! 🚀**
