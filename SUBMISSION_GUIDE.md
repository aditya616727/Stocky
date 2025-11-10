# 📋 Assignment Submission Guide

## Deliverables Checklist

Before submitting, verify you have:

### ✅ 1. GitHub Repository
- [ ] Repository is public or shared with reviewers
- [ ] All code is committed and pushed
- [ ] .gitignore properly excludes binary files and .env

### ✅ 2. Code Requirements
- [x] Written in **Golang** ✓
- [x] Uses **Gin** framework for routing ✓
- [x] Uses **github.com/sirupsen/logrus** for logging ✓
- [x] Database named **"assignment"** ✓
- [x] PostgreSQL database ✓

### ✅ 3. API Implementation
- [x] POST /api/v1/reward ✓
- [x] GET /api/v1/today-stocks/{userId} ✓
- [x] GET /api/v1/historical-inr/{userId} ✓
- [x] GET /api/v1/stats/{userId} ✓
- [x] GET /api/v1/portfolio/{userId} (Bonus) ✓

### ✅ 4. Database Schema
- [x] stock_rewards table with proper data types ✓
- [x] ledger_entries table (double-entry bookkeeping) ✓
- [x] stock_prices table ✓
- [x] NUMERIC(18,6) for fractional shares ✓
- [x] NUMERIC(18,4) for INR amounts ✓
- [x] Proper indexes and relationships ✓

### ✅ 5. Documentation
- [x] README.md with project explanation ✓
- [x] API specifications ✓
- [x] Database schema documentation ✓
- [x] Edge case handling explanations ✓
- [x] Scaling considerations ✓

### ✅ 6. Additional Files
- [x] Postman collection (Stocky.postman_collection.json) ✓
- [x] .env.example file ✓
- [x] .env file (if needed) ✓

### ✅ 7. Edge Cases Handled
- [x] Duplicate reward events / replay attacks ✓
- [x] Stock splits, mergers, or delisting (documented) ✓
- [x] Rounding errors in INR valuation ✓
- [x] Price API downtime or stale data ✓
- [x] Adjustments/refunds (documented) ✓

---

## Repository Structure

Your repository should look like this:

```
stocky/
├── main.go
├── go.mod
├── go.sum
├── README.md                              ⭐ Required
├── QUICKSTART.md
├── API_DOCS.md
├── PROJECT_SUMMARY.md
├── ARCHITECTURE.md
├── .env                                   ⭐ Include if needed
├── .env.example                          ⭐ Required
├── .gitignore
├── Makefile
├── setup_db.sh
├── test_api.sh
├── Stocky.postman_collection.json        ⭐ Required
│
├── config/
│   └── config.go
├── models/
│   └── models.go
├── database/
│   └── database.go
├── handlers/
│   └── reward_handler.go
├── routes/
│   └── routes.go
├── services/
│   └── price_service.go
└── (docs/ removed - no longer using Swagger)
```

---

## Pre-Submission Testing

### 1. Clean Build Test
```bash
# Remove all build artifacts
make clean

# Fresh build
go build -o bin/stocky main.go

# Should succeed without errors
```

### 2. Database Setup Test
```bash
# Ensure database exists
psql -U postgres -l | grep assignment

# If not, create it
psql -U postgres -c "CREATE DATABASE assignment;"
```

### 3. Application Start Test
```bash
# Start the application
go run main.go

# You should see:
# - "Connected to database successfully"
# - "Database migrations completed successfully"
# - "Starting stock price updater"
# - "Starting server on :8080"
```

### 4. API Functional Test
```bash
# Run the test script
./test_api.sh

# Or manually test
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test123",
    "stock_symbol": "RELIANCE",
    "quantity": 5.0,
    "idempotency_key": "test-001"
  }'
```

### 5. Postman Collection Test
1. Import `Stocky.postman_collection.json`
2. Run all requests in sequence
3. Verify all return successful responses

### 6. Database Verification
```sql
-- Connect to database
psql -U postgres -d assignment

-- Check tables exist
\dt

-- Check data
SELECT COUNT(*) FROM stock_rewards;
SELECT COUNT(*) FROM ledger_entries;
SELECT COUNT(*) FROM stock_prices;

-- Verify ledger balance (should be near 0 for each reward)
SELECT reward_id, SUM(amount) 
FROM ledger_entries 
GROUP BY reward_id;
```

---

## GitHub Repository Setup

### Step 1: Initialize Git (if not already done)
```bash
cd /Users/adityatomar/Desktop/atocky
git init
```

### Step 2: Add .gitignore
Already created in the project. Ensures:
- Binary files not tracked
- .env not committed (sensitive)
- IDE files excluded

### Step 3: Commit All Files
```bash
git add .
git commit -m "Initial commit: Stocky stock reward system

Features:
- 5 REST API endpoints (reward, today-stocks, historical-inr, stats, portfolio)
- Double-entry ledger system
- Hourly stock price updates
- Swagger UI documentation
- PostgreSQL with precise NUMERIC types
- Edge case handling (idempotency, rounding, stale prices)
- Complete documentation and Postman collection"
```

### Step 4: Create GitHub Repository
1. Go to https://github.com/new
2. Repository name: `stocky` or `stock-reward-system`
3. Description: "Stock reward system API built with Go, Gin, and PostgreSQL"
4. Public repository
5. Do NOT initialize with README (we already have one)

### Step 5: Push to GitHub
```bash
# Add remote
git remote add origin https://github.com/YOUR_USERNAME/stocky.git

# Push code
git branch -M main
git push -u origin main
```

---

## Email Template

Subject: **Stocky Assignment Submission - [Your Name]**

```
Dear Hiring Team,

I have completed the Stocky stock reward system assignment. Here are the details:

📌 GitHub Repository: https://github.com/YOUR_USERNAME/stocky

✅ Completed Features:
• All 5 required API endpoints implemented
• Double-entry ledger system with complete financial tracking
• Hourly stock price update service
• Interactive Swagger UI documentation
• PostgreSQL database with precise decimal types
• Idempotency for duplicate prevention
• Comprehensive edge case handling
• Complete documentation (README, API docs, Quick Start)
• Postman collection included

🛠 Tech Stack:
• Go 1.21+ with Gin framework
• Logrus for structured logging
• PostgreSQL (database: "assignment")
• GORM for database operations
• Swagger/OpenAPI documentation

📚 Documentation Files:
• README.md - Complete project overview
• QUICKSTART.md - 5-step setup guide
• API_DOCS.md - Detailed API documentation
• ARCHITECTURE.md - System architecture diagrams
• PROJECT_SUMMARY.md - Comprehensive summary
• Postman collection - Ready-to-use API tests
• .env file - Configuration template

🚀 Quick Start:
1. Clone the repository
2. Create PostgreSQL database: `CREATE DATABASE assignment;`
3. Run: `go mod download && ~/go/bin/swag init && go run main.go`
4. Access Swagger UI: http://localhost:8080/swagger/index.html

📊 Key Highlights:
• Fractional shares support (NUMERIC 18,6)
• Paise-level INR precision (NUMERIC 18,4)
• Automatic database migrations
• Background price updates every hour
• Complete ledger tracking (stock, cash, fees)
• Handles 10+ major Indian stocks

Please let me know if you need any clarification or additional information.

Best regards,
[Your Name]
[Your Email]
[Your Phone]
```

---

## Environment File for Reviewers

Include in email or repository README:

```env
# .env Configuration

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=assignment
SERVER_PORT=8080
```

Note: Reviewers should create the PostgreSQL database before running:
```bash
psql -U postgres -c "CREATE DATABASE assignment;"
```

---

## Common Issues & Solutions

### Issue 1: Database Connection Error
**Error**: `Failed to connect to database`

**Solution**:
```bash
# Check PostgreSQL is running
pg_isready

# Create database
psql -U postgres -c "CREATE DATABASE assignment;"

# Verify .env configuration
cat .env
```

### Issue 2: Swagger Not Loading
**Error**: `404 page not found on /swagger/`

**Solution**:
```bash
# Generate Swagger docs
~/go/bin/swag init

# Restart application
go run main.go
```

### Issue 3: Port Already in Use
**Error**: `bind: address already in use`

**Solution**:
```bash
# Kill existing process
lsof -ti:8080 | xargs kill -9

# Or change port in .env
echo "SERVER_PORT=8081" >> .env
```

### Issue 4: Missing Dependencies
**Error**: `package not found`

**Solution**:
```bash
go mod download
go mod tidy
```

---

## Final Verification Checklist

Before submitting, verify:

- [ ] Repository is public on GitHub
- [ ] All code is pushed to main branch
- [ ] README.md is comprehensive and clear
- [ ] .env.example is included (not .env with passwords)
- [ ] Postman collection is in repository
- [ ] Application builds successfully: `go build main.go`
- [ ] All tests pass (manual or automated)
- [ ] Swagger UI works at /swagger/index.html
- [ ] Database migrations run automatically
- [ ] Edge cases are documented in README.md
- [ ] Scaling considerations explained
- [ ] Code is well-commented
- [ ] No sensitive data committed (.env, passwords)

---

## Submission Timeline

1. **Code Complete** ✅ - Done
2. **Testing** ⏳ - Run test_api.sh
3. **Documentation Review** ⏳ - Read through all .md files
4. **Git Push** ⏳ - Push to GitHub
5. **Email Submission** ⏳ - Send to hiring team

---

## Questions Reviewers Might Ask

**Q1: Why use a double-entry ledger system?**
A: Ensures complete financial tracking and auditability. Every reward creates balanced entries (stock credit, cash debit, fees debit) that sum to zero.

**Q2: How do you handle stock splits?**
A: Currently documented as a future enhancement. Would implement an adjustments table to record corporate actions and apply them retroactively.

**Q3: What happens if the price service is down?**
A: Prices are cached in database. If missing, system generates a price immediately and caches it. Hourly updates keep prices fresh.

**Q4: How do you prevent duplicate rewards?**
A: Idempotency keys with unique database constraint. Duplicate requests return 409 Conflict with existing reward details.

**Q5: Can this system scale?**
A: Yes! Stateless design allows horizontal scaling. Database uses indexes for performance. Can add Redis caching, read replicas, and message queues for higher scale.

---

## Support Resources

- **README.md** - Complete project documentation
- **QUICKSTART.md** - Fast setup guide
- **API_DOCS.md** - Detailed API reference
- **ARCHITECTURE.md** - System design diagrams
- **PROJECT_SUMMARY.md** - Executive summary
- **Swagger UI** - Interactive API testing

---

Good luck with your submission! 🚀
