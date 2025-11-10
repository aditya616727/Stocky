# API Documentation - Stocky

## Base URL
```
http://localhost:8080/api/v1
```

## Interactive Documentation
Visit Swagger UI at: `http://localhost:8080/swagger/index.html`

---

## Endpoints

### 1. Create Stock Reward

**POST** `/reward`

Records a stock reward for a user. Creates ledger entries for double-entry bookkeeping.

#### Request Headers
```
Content-Type: application/json
```

#### Request Body
```json
{
  "user_id": "string (required)",
  "stock_symbol": "string (required)",
  "quantity": "number (required, > 0)",
  "rewarded_at": "timestamp (optional, defaults to now)",
  "idempotency_key": "string (optional, auto-generated if not provided)"
}
```

#### Example Request
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "rewarded_at": "2025-11-10T10:00:00Z",
    "idempotency_key": "reward-123-456"
  }'
```

#### Response (201 Created)
```json
{
  "id": 1,
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": 10.5,
  "rewarded_at": "2025-11-10T10:00:00Z",
  "current_price": 2450.75,
  "current_value": 25732.88
}
```

#### Response (409 Conflict) - Duplicate
```json
{
  "error": "Duplicate reward request",
  "existing_reward": { ... }
}
```

#### Response (400 Bad Request) - Validation Error
```json
{
  "error": "Key: 'RewardRequest.Quantity' Error:Field validation for 'Quantity' failed on the 'gt' tag"
}
```

---

### 2. Get Today's Stock Rewards

**GET** `/today-stocks/{userId}`

Returns all stock rewards for a specific user for today (current date).

#### Path Parameters
- `userId` (string, required): User identifier

#### Example Request
```bash
curl http://localhost:8080/api/v1/today-stocks/user123
```

#### Response (200 OK)
```json
{
  "user_id": "user123",
  "date": "2025-11-10",
  "rewards": [
    {
      "id": 1,
      "user_id": "user123",
      "stock_symbol": "RELIANCE",
      "quantity": 10.5,
      "rewarded_at": "2025-11-10T10:00:00Z",
      "idempotency_key": "reward-123-456",
      "created_at": "2025-11-10T10:00:05Z",
      "updated_at": "2025-11-10T10:00:05Z"
    },
    {
      "id": 2,
      "user_id": "user123",
      "stock_symbol": "TCS",
      "quantity": 5.25,
      "rewarded_at": "2025-11-10T11:30:00Z",
      "idempotency_key": "reward-123-457",
      "created_at": "2025-11-10T11:30:05Z",
      "updated_at": "2025-11-10T11:30:05Z"
    }
  ]
}
```

---

### 3. Get Historical INR Value

**GET** `/historical-inr/{userId}`

Returns the INR value of the user's stock rewards for all past days (up to yesterday). Values are calculated using current stock prices.

#### Path Parameters
- `userId` (string, required): User identifier

#### Example Request
```bash
curl http://localhost:8080/api/v1/historical-inr/user123
```

#### Response (200 OK)
```json
{
  "user_id": "user123",
  "daily": [
    {
      "date": "2025-11-08",
      "total_value": 98234.50
    },
    {
      "date": "2025-11-09",
      "total_value": 125000.50
    }
  ]
}
```

---

### 4. Get User Statistics

**GET** `/stats/{userId}`

Returns:
- Total shares rewarded today (grouped by stock symbol)
- Current INR value of the user's entire portfolio

#### Path Parameters
- `userId` (string, required): User identifier

#### Example Request
```bash
curl http://localhost:8080/api/v1/stats/user123
```

#### Response (200 OK)
```json
{
  "user_id": "user123",
  "today_rewards": [
    {
      "stock_symbol": "RELIANCE",
      "quantity": 10.5
    },
    {
      "stock_symbol": "TCS",
      "quantity": 5.25
    }
  ],
  "current_portfolio_value": 250000.75
}
```

---

### 5. Get Portfolio (Bonus)

**GET** `/portfolio/{userId}`

Returns detailed holdings per stock symbol with current INR value.

#### Path Parameters
- `userId` (string, required): User identifier

#### Example Request
```bash
curl http://localhost:8080/api/v1/portfolio/user123
```

#### Response (200 OK)
```json
{
  "user_id": "user123",
  "holdings": [
    {
      "stock_symbol": "RELIANCE",
      "total_shares": 25.5,
      "current_price": 2450.75,
      "current_value": 62489.13
    },
    {
      "stock_symbol": "TCS",
      "total_shares": 12.75,
      "current_price": 3520.50,
      "current_value": 44886.38
    },
    {
      "stock_symbol": "INFOSYS",
      "total_shares": 8.0,
      "current_price": 1485.25,
      "current_value": 11882.00
    }
  ],
  "total_value": 119257.51
}
```

---

## Supported Stock Symbols

The system generates realistic prices for these Indian stocks:
- **RELIANCE** - Reliance Industries
- **TCS** - Tata Consultancy Services
- **INFOSYS** - Infosys Limited
- **HDFC** - HDFC Bank
- **ICICI** - ICICI Bank
- **SBI** - State Bank of India
- **WIPRO** - Wipro Limited
- **BHARTI** - Bharti Airtel
- **ITC** - ITC Limited
- **HCLTECH** - HCL Technologies

Custom stock symbols can be used, and the system will generate default prices.

---

## Error Responses

### 400 Bad Request
Invalid request format or validation failure.
```json
{
  "error": "Error message describing the issue"
}
```

### 409 Conflict
Duplicate idempotency key detected.
```json
{
  "error": "Duplicate reward request",
  "existing_reward": { ... }
}
```

### 500 Internal Server Error
Server-side error occurred.
```json
{
  "error": "Error message"
}
```

---

## Idempotency

The API uses idempotency keys to prevent duplicate reward creation. 

**How it works:**
1. Include an `idempotency_key` in the POST /reward request
2. If the same key is sent twice, the API returns 409 Conflict
3. If no key is provided, one is auto-generated: `{user_id}-{stock_symbol}-{timestamp}`

**Example:**
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": 10.5,
  "idempotency_key": "onboarding-reward-user123-2025-11-10"
}
```

---

## Data Precision

- **Stock Quantities**: Supports up to 6 decimal places (fractional shares)
- **INR Amounts**: Supports up to 4 decimal places (paise level)
- **Stock Prices**: Updated hourly by background service

---

## Rate Limiting

No rate limiting is currently implemented in this version.

---

## Testing Scenarios

### Scenario 1: Onboarding Reward
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "newuser456",
    "stock_symbol": "TCS",
    "quantity": 1.0,
    "idempotency_key": "onboarding-newuser456"
  }'
```

### Scenario 2: Referral Reward
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "INFOSYS",
    "quantity": 2.5,
    "idempotency_key": "referral-user123-ref789"
  }'
```

### Scenario 3: Trading Milestone
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "trader999",
    "stock_symbol": "RELIANCE",
    "quantity": 5.0,
    "idempotency_key": "milestone-100trades-trader999"
  }'
```

### Scenario 4: Fractional Shares
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "HDFC",
    "quantity": 0.123456,
    "idempotency_key": "microtransaction-user123"
  }'
```

---

## Postman Collection

Import `Stocky.postman_collection.json` for pre-configured requests.

The collection includes:
- All CRUD operations
- Multiple user scenarios
- Duplicate detection tests
- Edge case examples

---

## Support

For technical issues or questions:
1. Check Swagger UI for live API testing
2. Review README.md for system architecture
3. Check logs for detailed error messages
