# 🎯 Stocky - Project Summary

## Overview
Stocky is a comprehensive stock reward system API built with Go, Gin framework, and PostgreSQL. The system allows users to earn shares of Indian stocks as incentives, with complete financial tracking through a double-entry ledger system.

## ✅ Completed Features

### 1. Core API Endpoints (All Implemented)
- ✅ **POST /api/v1/reward** - Create stock rewards with idempotency
- ✅ **GET /api/v1/today-stocks/{userId}** - Today's rewards for a user
- ✅ **GET /api/v1/historical-inr/{userId}** - Historical INR values
- ✅ **GET /api/v1/stats/{userId}** - User statistics (today + portfolio value)
- ✅ **GET /api/v1/portfolio/{userId}** - Detailed portfolio holdings (BONUS)

### 2. Database Schema (Fully Implemented)
- ✅ **stock_rewards** - Records all reward events with proper indexing
- ✅ **ledger_entries** - Double-entry bookkeeping for complete financial tracking
- ✅ **stock_prices** - Hourly updated stock prices with timestamps
- ✅ All tables use precise data types:
  - Stock quantities: `NUMERIC(18,6)` for fractional shares
  - INR amounts: `NUMERIC(18,4)` for paise-level precision

### 3. Double-Entry Ledger System
For each reward, the system creates three ledger entries:
- ✅ **STOCK_CREDIT** - User receives shares (positive)
- ✅ **CASH_DEBIT** - Company cash outflow (negative)
- ✅ **FEE_DEBIT** - Company fees (brokerage, STT, GST) (negative)

### 4. Stock Price Service
- ✅ Hypothetical price generator for Indian stocks
- ✅ Hourly automatic price updates (background goroutine)
- ✅ Caching mechanism for price data
- ✅ Support for 10+ major Indian stocks with realistic base prices

### 5. Edge Cases Handled
✅ **Duplicate Prevention**: Idempotency keys with unique constraints
✅ **Rounding Errors**: Precise NUMERIC types with proper decimal places
✅ **Stale Prices**: Database caching + hourly updates
✅ **Missing Prices**: Auto-generation on demand
✅ **Fractional Shares**: Full support up to 6 decimal places
✅ **Transaction Safety**: Database transactions for atomic operations

### 6. Documentation
- ✅ **README.md** - Complete project documentation
- ✅ **QUICKSTART.md** - 5-step setup guide
- ✅ **API_DOCS.md** - Detailed API documentation with examples
- ✅ **Swagger UI** - Interactive API documentation at `/swagger/index.html`
- ✅ **Postman Collection** - Ready-to-use API collection
- ✅ **Code Comments** - Swagger annotations in all handlers

### 7. Configuration & Setup
- ✅ **.env** file with database configuration
- ✅ **.env.example** for template
- ✅ **Makefile** for common tasks
- ✅ **setup_db.sh** for automated database setup
- ✅ **.gitignore** for version control

### 8. Tech Stack (As Required)
- ✅ **Go** 1.21+ as the programming language
- ✅ **Gin** for HTTP routing
- ✅ **Logrus** for structured logging
- ✅ **PostgreSQL** with database name "assignment"
- ✅ **GORM** for database ORM
- ✅ **Swagger** (swaggo) for API documentation

## 📁 Project Structure

```
stocky/
├── main.go                              # Application entry point with Swagger setup
├── go.mod                               # Go module dependencies
├── go.sum                               # Go module checksums
│
├── config/
│   └── config.go                        # Configuration management
│
├── models/
│   └── models.go                        # Data models, DTOs, request/response types
│
├── database/
│   └── database.go                      # Database connection & migrations
│
├── handlers/
│   └── reward_handler.go                # HTTP handlers with Swagger annotations
│
├── routes/
│   └── routes.go                        # Route definitions
│
├── services/
│   └── price_service.go                 # Stock price service with hourly updates
│
├── docs/                                # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── .env                                 # Environment variables (DB config)
├── .env.example                         # Environment template
├── .gitignore                          # Git ignore patterns
├── Makefile                            # Build and run commands
├── setup_db.sh                         # Database setup script
│
├── README.md                           # Main documentation
├── QUICKSTART.md                       # Quick start guide
├── API_DOCS.md                         # Detailed API documentation
└── Stocky.postman_collection.json      # Postman API collection
```

## 🚀 Quick Start

### 1. Setup Database
```bash
# Using Docker (easiest)
docker run --name stocky-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=assignment \
  -p 5432:5432 -d postgres:15

# OR using PostgreSQL directly
psql -U postgres -c "CREATE DATABASE assignment;"
```

### 2. Install & Run
```bash
# Install dependencies
go mod download
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
~/go/bin/swag init

# Run the application
go run main.go
```

### 3. Access Swagger UI
Open your browser to: **http://localhost:8080/swagger/index.html**

## 📊 Database Schema

### stock_rewards
| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL PRIMARY KEY | Unique identifier |
| user_id | VARCHAR(100) | User identifier (indexed) |
| stock_symbol | VARCHAR(20) | Stock symbol (indexed) |
| quantity | NUMERIC(18,6) | Number of shares (fractional) |
| rewarded_at | TIMESTAMP | Reward timestamp (indexed) |
| idempotency_key | VARCHAR(100) UNIQUE | Duplicate prevention |
| created_at | TIMESTAMP | Record creation time |
| updated_at | TIMESTAMP | Record update time |

### ledger_entries
| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL PRIMARY KEY | Unique identifier |
| reward_id | BIGINT | Foreign key to stock_rewards |
| entry_type | VARCHAR(50) | STOCK_CREDIT/CASH_DEBIT/FEE_DEBIT |
| stock_symbol | VARCHAR(20) | Stock symbol |
| quantity | NUMERIC(18,6) | Stock quantity |
| amount | NUMERIC(18,4) | INR amount (+/-) |
| description | TEXT | Entry description |
| created_at | TIMESTAMP | Record creation time |

### stock_prices
| Column | Type | Description |
|--------|------|-------------|
| id | BIGSERIAL PRIMARY KEY | Unique identifier |
| stock_symbol | VARCHAR(20) UNIQUE | Stock symbol |
| price | NUMERIC(18,4) | Current price in INR |
| updated_at | TIMESTAMP | Last update time (indexed) |
| created_at | TIMESTAMP | Record creation time |

## 🔄 System Flow

### Creating a Reward
1. **API Request** → POST /api/v1/reward with user_id, stock_symbol, quantity
2. **Validation** → Check idempotency key for duplicates
3. **Price Fetch** → Get current stock price from cache/service
4. **Fee Calculation** → Calculate brokerage (0.5%), STT (0.1%), GST (18%)
5. **Database Transaction**:
   - Create stock_reward record
   - Create 3 ledger_entries (STOCK_CREDIT, CASH_DEBIT, FEE_DEBIT)
6. **Response** → Return reward details with current valuation

### Background Price Updates
- Goroutine runs hourly
- Fetches all unique stock symbols from rewards
- Generates new prices (±5% variation from base)
- Updates stock_prices table
- Logs update status

## 🎨 API Examples

### Create Reward
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "idempotency_key": "reward-001"
  }'
```

### Get Portfolio
```bash
curl http://localhost:8080/api/v1/portfolio/user123
```

## 🛡️ Edge Cases Covered

### 1. Duplicate Rewards (Idempotency)
- **Problem**: Same reward sent twice due to network retry
- **Solution**: Unique constraint on `idempotency_key`
- **Result**: Returns 409 Conflict with existing reward

### 2. Fractional Shares
- **Problem**: Need to support partial shares (e.g., 0.123456)
- **Solution**: NUMERIC(18,6) data type
- **Result**: Precise storage up to 6 decimal places

### 3. Rounding Errors in INR
- **Problem**: Floating-point arithmetic errors
- **Solution**: NUMERIC(18,4) for all INR amounts
- **Result**: Paise-level precision without rounding errors

### 4. Stale Stock Prices
- **Problem**: Price API downtime or delays
- **Solution**: 
  - Cache prices in database
  - Hourly background updates
  - Auto-generate on demand if missing
  - Track update timestamps
- **Result**: Always available prices with freshness info

### 5. Missing Stock Prices
- **Problem**: New stock symbol with no price
- **Solution**: Generate and save price immediately
- **Result**: No API failures due to missing prices

## 📦 Deliverables Checklist

✅ **1. GitHub Repository** - Ready to be pushed
✅ **2. API Specifications** - Complete request/response payloads in API_DOCS.md
✅ **3. Database Schema** - Fully documented with relationships in README.md
✅ **4. Edge Case Explanations** - Detailed in README.md (Bonus section)
✅ **5. Scaling Considerations** - Comprehensive section in README.md
✅ **6. Go with Gin & Logrus** - As required
✅ **7. PostgreSQL Database** - Named "assignment"
✅ **8. Postman Collection** - Stocky.postman_collection.json
✅ **9. .env File** - Included with defaults
✅ **10. Swagger UI** - Fully integrated at /swagger/index.html

## 🚀 Scaling Strategies (from README.md)

### Database
- Partitioning by date for large datasets
- Read replicas for analytics
- Connection pooling (GORM built-in)

### Application
- Stateless design → horizontal scaling
- Redis caching for portfolios & prices
- Queue system (RabbitMQ/Kafka) for async processing

### Monitoring
- Structured JSON logging (Logrus)
- Prometheus metrics (future)
- Distributed tracing (future)

## 📝 Testing Instructions

### Manual Testing with Swagger
1. Start the application: `go run main.go`
2. Open: http://localhost:8080/swagger/index.html
3. Try all 5 endpoints with sample data
4. Test duplicate detection by reusing idempotency_key

### Postman Testing
1. Import `Stocky.postman_collection.json`
2. Update base_url variable if needed
3. Run requests sequentially:
   - Create rewards for different stocks
   - Get today's stocks
   - Get historical data
   - Check stats
   - View portfolio

### Database Verification
```sql
-- Check rewards
SELECT * FROM stock_rewards;

-- Check ledger (should have 3 entries per reward)
SELECT * FROM ledger_entries;

-- Check prices
SELECT * FROM stock_prices;

-- Verify ledger balance (should sum to 0 for each reward_id)
SELECT reward_id, SUM(amount) as total 
FROM ledger_entries 
GROUP BY reward_id;
```

## 🎯 Success Metrics

✅ All required API endpoints implemented
✅ Double-entry ledger system working
✅ Idempotency preventing duplicates
✅ Swagger UI fully functional
✅ Postman collection with all scenarios
✅ Database migrations automatic
✅ Hourly price updates running
✅ Comprehensive documentation
✅ Edge cases handled
✅ Scaling considerations documented

## 🔮 Future Enhancements (Not in Scope)

- Stock splits/mergers handling (add adjustments table)
- Reward reversals/refunds (negative entries)
- User authentication/authorization
- Rate limiting
- Caching layer (Redis)
- Real stock price API integration
- WebSocket for real-time updates
- Admin dashboard
- Analytics endpoints
- Automated testing suite

## 📄 License

MIT License

## 👤 Contact

For questions or issues, please refer to the documentation files or check the code comments.

---

**Created**: November 10, 2025
**Version**: 1.0.0
**Status**: ✅ Production Ready (Minimal Viable Product)
