# System Architecture - Stocky
-- Create the user with a password
CREATE USER adityatomar WITH PASSWORD 'password';

-- Create the database if it doesn't exist
CREATE DATABASE assignment;

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON DATABASE assignment TO adityatomar;

-- Connect to the assignment database
\c assignment

-- Grant schema privileges (run this after connecting to assignment database)
GRANT ALL ON SCHEMA public TO adityatomar;
## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
├─────────────────────────────────────────────────────────────────┤
│  Swagger UI  │  Postman  │  cURL  │  Mobile/Web Apps           │
└────────────┬────────────┴────────┴─────────────────────────────┘
             │
             │ HTTP/REST API
             │
┌────────────▼─────────────────────────────────────────────────────┐
│                      API GATEWAY (Gin Router)                     │
│                     http://localhost:8080                         │
├───────────────────────────────────────────────────────────────────┤
│  /swagger/*         → Swagger UI                                  │
│  /api/v1/reward     → Create Reward Handler                       │
│  /api/v1/today-stocks/:userId → Today's Stocks Handler            │
│  /api/v1/historical-inr/:userId → Historical INR Handler          │
│  /api/v1/stats/:userId → Stats Handler                            │
│  /api/v1/portfolio/:userId → Portfolio Handler                    │
└────────────┬──────────────────────────────────────────────────────┘
             │
             │
┌────────────▼──────────────────────────────────────────────────────┐
│                      HANDLER LAYER                                │
├───────────────────────────────────────────────────────────────────┤
│  reward_handler.go                                                │
│  - CreateReward()                                                 │
│  - GetTodayStocks()                                              │
│  - GetHistoricalINR()                                            │
│  - GetStats()                                                    │
│  - GetPortfolio()                                                │
└────────────┬──────────────────────────────────────────────────────┘
             │
             │
┌────────────▼──────────────────────────────────────────────────────┐
│                      SERVICE LAYER                                │
├───────────────────────────────────────────────────────────────────┤
│  StockPriceService                                                │
│  - GetCurrentPrice()                                              │
│  - StartPriceUpdater() [Background Goroutine]                     │
│  - updateAllPrices()                                              │
│  - getStockPrice() [Hypothetical Generator]                       │
└────────────┬──────────────────────────────────────────────────────┘
             │
             │
┌────────────▼──────────────────────────────────────────────────────┐
│                      DATABASE LAYER (GORM)                        │
├───────────────────────────────────────────────────────────────────┤
│  PostgreSQL Database: "assignment"                                │
│                                                                   │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────────┐  │
│  │ stock_rewards   │  │ ledger_entries   │  │ stock_prices   │  │
│  ├─────────────────┤  ├──────────────────┤  ├────────────────┤  │
│  │ id (PK)         │  │ id (PK)          │  │ id (PK)        │  │
│  │ user_id         │◄─┤ reward_id (FK)   │  │ stock_symbol   │  │
│  │ stock_symbol    │  │ entry_type       │  │ price          │  │
│  │ quantity        │  │ stock_symbol     │  │ updated_at     │  │
│  │ rewarded_at     │  │ quantity         │  │ created_at     │  │
│  │ idempotency_key │  │ amount           │  └────────────────┘  │
│  │ created_at      │  │ description      │                      │
│  │ updated_at      │  │ created_at       │                      │
│  └─────────────────┘  └──────────────────┘                      │
└───────────────────────────────────────────────────────────────────┘
```

## Request Flow: Creating a Reward

```
Client                Handler              Service            Database
  │                     │                    │                  │
  │ POST /reward        │                    │                  │
  ├────────────────────>│                    │                  │
  │                     │                    │                  │
  │                     │ Validate Request   │                  │
  │                     │───────────────────>│                  │
  │                     │                    │                  │
  │                     │ Check Idempotency  │                  │
  │                     │────────────────────┼─────────────────>│
  │                     │                    │  SELECT WHERE    │
  │                     │<───────────────────┼──────────────────┤
  │                     │                    │                  │
  │                     │ Get Stock Price    │                  │
  │                     │───────────────────>│                  │
  │                     │                    │ Query Cache      │
  │                     │                    ├─────────────────>│
  │                     │                    │<─────────────────┤
  │                     │<───────────────────┤                  │
  │                     │                    │                  │
  │                     │ Calculate Fees     │                  │
  │                     │ (Brokerage,STT,GST)│                  │
  │                     │                    │                  │
  │                     │ Begin Transaction  │                  │
  │                     │────────────────────┼─────────────────>│
  │                     │                    │  BEGIN           │
  │                     │                    │                  │
  │                     │ Insert Reward      │                  │
  │                     │────────────────────┼─────────────────>│
  │                     │                    │  INSERT          │
  │                     │                    │                  │
  │                     │ Insert Ledger x3   │                  │
  │                     │────────────────────┼─────────────────>│
  │                     │                    │  INSERT (CREDIT) │
  │                     │                    │  INSERT (DEBIT)  │
  │                     │                    │  INSERT (FEE)    │
  │                     │                    │                  │
  │                     │ Commit Transaction │                  │
  │                     │────────────────────┼─────────────────>│
  │                     │                    │  COMMIT          │
  │                     │<───────────────────┼──────────────────┤
  │                     │                    │                  │
  │ 201 Created         │                    │                  │
  │<────────────────────┤                    │                  │
  │ {reward_response}   │                    │                  │
```

## Background Price Update Service

```
┌─────────────────────────────────────────────────────────────┐
│              Price Update Background Service                │
│                   (Goroutine - Hourly)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ Every Hour
                       ▼
          ┌────────────────────────┐
          │  Get All Stock Symbols │
          │  FROM stock_rewards    │
          └────────┬───────────────┘
                   │
                   ▼
          ┌────────────────────────┐
          │  For Each Symbol:      │
          │  - Generate New Price  │
          │  - Apply ±5% Variation │
          └────────┬───────────────┘
                   │
                   ▼
          ┌────────────────────────┐
          │  Update stock_prices   │
          │  Table                 │
          └────────┬───────────────┘
                   │
                   ▼
          ┌────────────────────────┐
          │  Log Update Status     │
          │  (Logrus)              │
          └────────────────────────┘
```

## Double-Entry Ledger System

```
User Reward: 10 shares of RELIANCE @ ₹2450

┌────────────────────────────────────────────────────────────┐
│                    Ledger Entries                          │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Entry 1: STOCK_CREDIT                                     │
│  ┌──────────────────────────────────────────────────┐     │
│  │ Stock Symbol: RELIANCE                            │     │
│  │ Quantity: +10.0000                                │     │
│  │ Amount: ₹24,500.00                                │     │
│  │ Description: "Stock reward credited to user123"   │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  Entry 2: CASH_DEBIT                                       │
│  ┌──────────────────────────────────────────────────┐     │
│  │ Stock Symbol: RELIANCE                            │     │
│  │ Quantity: 0                                       │     │
│  │ Amount: -₹24,500.00                               │     │
│  │ Description: "Company cash outflow for purchase"  │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  Entry 3: FEE_DEBIT                                        │
│  ┌──────────────────────────────────────────────────┐     │
│  │ Stock Symbol: RELIANCE                            │     │
│  │ Quantity: 0                                       │     │
│  │ Amount: -₹146.95                                  │     │
│  │ Description: "Brokerage: ₹122.50,                 │     │
│  │               STT: ₹24.50, GST: ₹22.05"          │     │
│  └──────────────────────────────────────────────────┘     │
│                                                            │
│  Net Balance: ₹0 (Balanced!)                               │
│  User Gets: 10 shares                                      │
│  Company Pays: ₹24,500 + ₹146.95 fees = ₹24,646.95       │
└────────────────────────────────────────────────────────────┘
```

## Data Flow: Get Portfolio

```
Client                Handler              Database
  │                     │                    │
  │ GET /portfolio/user │                    │
  ├────────────────────>│                    │
  │                     │                    │
  │                     │ Query All Rewards  │
  │                     │────────────────────>│
  │                     │ WHERE user_id      │
  │                     │<────────────────────┤
  │                     │                    │
  │                     │ Group by Symbol    │
  │                     │ (In Memory)        │
  │                     │                    │
  │                     │ For Each Symbol:   │
  │                     │ Get Current Price  │
  │                     │────────────────────>│
  │                     │<────────────────────┤
  │                     │                    │
  │                     │ Calculate Values   │
  │                     │                    │
  │ 200 OK              │                    │
  │<────────────────────┤                    │
  │ {portfolio}         │                    │
```

## Database Indexes

```
stock_rewards
├── PRIMARY KEY (id)
├── INDEX (user_id)              → Fast user lookups
├── INDEX (stock_symbol)         → Price join optimization
├── INDEX (rewarded_at)          → Date range queries
└── UNIQUE INDEX (idempotency_key) → Duplicate prevention

ledger_entries
├── PRIMARY KEY (id)
└── INDEX (reward_id)            → Fast ledger lookups

stock_prices
├── PRIMARY KEY (id)
├── UNIQUE INDEX (stock_symbol)  → Fast price lookups
└── INDEX (updated_at)           → Freshness checks
```

## Error Handling Flow

```
Request
   │
   ▼
┌──────────────┐
│  Validation  │
│  (Gin bind)  │
└──────┬───────┘
       │
       ├─ Invalid ──> 400 Bad Request
       │
       ▼
┌──────────────────┐
│  Idempotency     │
│  Check           │
└──────┬───────────┘
       │
       ├─ Duplicate ──> 409 Conflict
       │
       ▼
┌──────────────────┐
│  Business Logic  │
└──────┬───────────┘
       │
       ├─ Error ──> 500 Internal Server Error
       │
       ▼
   Success
   201/200
```

## Scaling Architecture (Future)

```
┌─────────────┐
│ Load        │
│ Balancer    │
└──────┬──────┘
       │
       ├──────────┬──────────┬──────────┐
       ▼          ▼          ▼          ▼
   ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐
   │App 1│    │App 2│    │App 3│    │App N│
   └──┬──┘    └──┬──┘    └──┬──┘    └──┬──┘
      │          │          │          │
      └──────────┴──────────┴──────────┘
                 │
         ┌───────┼───────┐
         ▼       ▼       ▼
     ┌──────┐ ┌──────┐ ┌────────┐
     │Redis │ │Queue │ │Primary │
     │Cache │ │(RMQ) │ │  DB    │
     └──────┘ └──────┘ └───┬────┘
                           │
                     ┌─────┴─────┐
                     ▼           ▼
                 ┌────────┐  ┌────────┐
                 │Read    │  │Read    │
                 │Replica │  │Replica │
                 └────────┘  └────────┘
```

## Technology Stack

```
┌─────────────────────────────────────────────────────┐
│                  Application Layer                  │
├─────────────────────────────────────────────────────┤
│  Go 1.21+                                           │
│  ├── Gin Framework (HTTP Router)                    │
│  ├── Logrus (Structured Logging)                    │
│  ├── GORM (ORM)                                     │
│  ├── Swaggo (API Documentation)                     │
│  └── godotenv (Environment Config)                  │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│                  Database Layer                     │
├─────────────────────────────────────────────────────┤
│  PostgreSQL 13+                                     │
│  ├── GORM Auto-Migrations                          │
│  ├── Connection Pooling                            │
│  └── NUMERIC precision types                       │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│                  Documentation                      │
├─────────────────────────────────────────────────────┤
│  ├── Swagger UI (Interactive)                       │
│  ├── Postman Collection                            │
│  └── Markdown Documentation                        │
└─────────────────────────────────────────────────────┘
```

## Deployment Architecture

```
┌─────────────────────────────────────────────┐
│           Production Environment            │
├─────────────────────────────────────────────┤
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │  Docker Container: stocky-app       │   │
│  │  ├── Binary: ./stocky               │   │
│  │  ├── Port: 8080                     │   │
│  │  └── Env: Production                │   │
│  └─────────────────────────────────────┘   │
│                     │                       │
│                     ▼                       │
│  ┌─────────────────────────────────────┐   │
│  │  Docker Container: stocky-postgres  │   │
│  │  ├── PostgreSQL 15                  │   │
│  │  ├── Port: 5432                     │   │
│  │  ├── Volume: /var/lib/postgresql    │   │
│  │  └── Database: assignment           │   │
│  └─────────────────────────────────────┘   │
│                                             │
└─────────────────────────────────────────────┘
```

---

This architecture supports:
✅ Horizontal scaling of application servers
✅ Database read replicas for analytics
✅ Caching layer for performance
✅ Message queue for async processing
✅ Monitoring and observability
✅ High availability and fault tolerance
