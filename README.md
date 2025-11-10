# Stocky - Stock Reward System

A comprehensive stock reward system API built with Go, Gin, and PostgreSQL that allows users to earn shares of Indian stocks as incentives.

## Features

- ✅ **Stock Reward Management**: Record and track stock rewards for users
- ✅ **Double-Entry Bookkeeping**: Complete ledger system tracking stock units, cash outflow, and fees
- ✅ **Real-time Stock Prices**: Hypothetical stock price service with hourly updates
- ✅ **Portfolio Management**: Track user holdings and current valuations
- ✅ **Idempotency**: Prevent duplicate reward events
- ✅ **Edge Case Handling**: Duplicate prevention, precise decimal calculations

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Logging**: Logrus

## Project Structure

```
stocky/
├── main.go                  # Application entry point
├── config/                  # Configuration management
│   └── config.go
├── models/                  # Data models and DTOs
│   └── models.go
├── database/                # Database connection and migrations
│   └── database.go
├── handlers/                # HTTP request handlers
│   └── reward_handler.go
├── routes/                  # Route definitions
│   └── routes.go
├── services/                # Business logic services
│   └── price_service.go
├── .env                     # Environment variables
├── .env.example             # Example environment variables
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── README.md                # This file
└── Stocky.postman_collection.json  # Postman API collection
```

## Database Schema

### Tables

#### 1. `stock_rewards`
Records each stock reward event.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL | Primary key |
| user_id | VARCHAR(100) | User identifier |
| stock_symbol | VARCHAR(20) | Stock symbol (e.g., RELIANCE, TCS) |
| quantity | NUMERIC(18,6) | Number of shares (supports fractional) |
| rewarded_at | TIMESTAMP | When the reward was given |
| idempotency_key | VARCHAR(100) | Unique key to prevent duplicates |
| created_at | TIMESTAMP | Record creation time |
| updated_at | TIMESTAMP | Record update time |

**Indexes**: user_id, stock_symbol, rewarded_at, idempotency_key (unique)

#### 2. `ledger_entries`
Double-entry bookkeeping for complete financial tracking.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL | Primary key |
| reward_id | BIGINT | Foreign key to stock_rewards |
| entry_type | VARCHAR(50) | STOCK_CREDIT, CASH_DEBIT, FEE_DEBIT |
| stock_symbol | VARCHAR(20) | Stock symbol |
| quantity | NUMERIC(18,6) | Stock quantity (for STOCK_CREDIT) |
| amount | NUMERIC(18,4) | INR amount (positive/negative) |
| description | TEXT | Entry description |
| created_at | TIMESTAMP | Record creation time |

**Indexes**: reward_id

#### 3. `stock_prices`
Current stock prices updated hourly.

| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL | Primary key |
| stock_symbol | VARCHAR(20) | Stock symbol (unique) |
| price | NUMERIC(18,4) | Current price in INR |
| updated_at | TIMESTAMP | Last price update time |
| created_at | TIMESTAMP | Record creation time |

**Indexes**: stock_symbol (unique), updated_at

## API Endpoints

### 1. Create Reward
**POST** `/api/v1/reward`

Record a stock reward for a user.

**Request Body**:
```json
{
  "user_id": "user123",
  "stock_symbol": "RELIANCE",
  "quantity": 10.5,
  "rewarded_at": "2025-11-10T10:00:00Z",
  "idempotency_key": "reward-123-456"
}
```

**Response** (201 Created):
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

### 2. Get Today's Stocks
**GET** `/api/v1/today-stocks/{userId}`

Get all stock rewards for a user for today.

**Response** (200 OK):
```json
{
  "user_id": "user123",
  "date": "2025-11-10",
  "rewards": [...]
}
```

### 3. Get Historical INR Value
**GET** `/api/v1/historical-inr/{userId}`

Get INR value of user's stock rewards for all past days.

**Response** (200 OK):
```json
{
  "user_id": "user123",
  "daily": [
    {
      "date": "2025-11-09",
      "total_value": 125000.50
    }
  ]
}
```

### 4. Get User Stats
**GET** `/api/v1/stats/{userId}`

Get today's rewards and current portfolio value.

**Response** (200 OK):
```json
{
  "user_id": "user123",
  "today_rewards": [
    {
      "stock_symbol": "RELIANCE",
      "quantity": 15.5
    }
  ],
  "current_portfolio_value": 250000.75
}
```

### 5. Get Portfolio (Bonus)
**GET** `/api/v1/portfolio/{userId}`

Get detailed holdings per stock symbol.

**Response** (200 OK):
```json
{
  "user_id": "user123",
  "holdings": [
    {
      "stock_symbol": "RELIANCE",
      "total_shares": 25.5,
      "current_price": 2450.75,
      "current_value": 62489.13
    }
  ],
  "total_value": 250000.75
}
```

## Setup and Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 13 or higher
- Make (optional)

### Steps

1. **Clone the repository**
```bash
git clone <repository-url>
cd stocky
```

2. **Install Go dependencies**
```bash
go mod download
```

3. **Setup PostgreSQL**
```bash
# Create database
psql -U postgres
CREATE DATABASE assignment;
\q
```

4. **Configure environment variables**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. **Run the application**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Testing with Postman

Import the `Stocky.postman_collection.json` file into Postman to test all endpoints.

## Edge Cases Handled

### 1. Duplicate Reward Events / Replay Attacks
- **Solution**: Idempotency keys (unique constraint)
- If the same idempotency_key is sent twice, the API returns HTTP 409 Conflict with the existing reward

### 2. Stock Splits, Mergers, or Delisting
- **Current**: Not implemented in minimal version
- **Recommendation**: Add an `adjustments` table to record corporate actions and apply them retroactively to holdings

### 3. Rounding Errors in INR Valuation
- **Solution**: Using `NUMERIC(18,4)` for precise decimal storage
- All calculations maintain precision before final rounding

### 4. Price API Downtime or Stale Data
- **Solution**: 
  - Prices are cached in database
  - If price doesn't exist, generate and cache immediately
  - Background service updates prices hourly
  - Last update timestamp tracked for monitoring

### 5. Adjustments/Refunds of Previously Given Rewards
- **Current**: Not implemented in minimal version
- **Recommendation**: Add reversal entries in ledger with negative quantities

## Scaling Considerations

### Database Optimization
- **Indexes**: Added on frequently queried columns (user_id, stock_symbol, rewarded_at)
- **Partitioning**: Consider partitioning `stock_rewards` and `ledger_entries` by date for large datasets
- **Read Replicas**: Use PostgreSQL read replicas for analytics queries

### Application Scaling
- **Stateless Design**: Application is stateless and can be horizontally scaled
- **Connection Pooling**: GORM manages database connection pool
- **Caching**: Add Redis for frequently accessed data (user portfolios, stock prices)

### Background Jobs
- **Price Updates**: Currently runs in-process; move to separate worker service for production
- **Queue System**: Use RabbitMQ/Kafka for async reward processing

### Monitoring
- **Logging**: Structured JSON logs using Logrus
- **Metrics**: Add Prometheus metrics for API latency, database queries
- **Tracing**: Add distributed tracing (Jaeger/OpenTelemetry)

## Double-Entry Bookkeeping System

For each reward, three ledger entries are created:

1. **STOCK_CREDIT**: User receives shares (positive quantity)
2. **CASH_DEBIT**: Company cash outflow for stock purchase (negative amount)
3. **FEE_DEBIT**: Company incurs fees (negative amount)
   - Brokerage: 0.5% of stock value
   - STT (Securities Transaction Tax): 0.1%
   - GST on Brokerage: 18%

This ensures complete financial tracking and auditability.

## Data Types

- **Stock Quantities**: `NUMERIC(18,6)` - Supports fractional shares up to 6 decimal places
- **INR Amounts**: `NUMERIC(18,4)` - Precise to 4 decimal places (paise level)
- **Timestamps**: `TIMESTAMP WITH TIME ZONE` - For accurate time tracking

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o stocky main.go
./stocky
```

### Docker Support (Optional)
```bash
docker build -t stocky .
docker run -p 8080:8080 --env-file .env stocky
```

## License

MIT

## Contact

For questions or support, please contact the development team.
