
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
