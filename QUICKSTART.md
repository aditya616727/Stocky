# Quick Start Guide - Stocky

## Prerequisites

1. **Go 1.21+** - [Install Go](https://golang.org/doc/install)
2. **PostgreSQL 13+** - [Install PostgreSQL](https://www.postgresql.org/download/)
   - OR use Docker: `docker run --name stocky-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=assignment -p 5432:5432 -d postgres:15`

## Setup in 5 Steps

### 1. Clone and Navigate
```bash
cd atocky
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Database

**Option A: Using PostgreSQL directly**
```bash
# Create database
psql -U postgres -c "CREATE DATABASE assignment;"

# Or use the setup script
./setup_db.sh
```

**Option B: Using Docker**
```bash
docker run --name stocky-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=assignment \
  -p 5432:5432 -d postgres:15
```

### 4. Configure Environment
```bash
# Copy example env file
cp .env.example .env

# Edit .env if needed (defaults should work)
```

### 5. Run the Application
```bash
go run main.go
```

## Access the Application

- **API Base URL**: `http://localhost:8080/api/v1`
- **Health Check**: The app will log "Starting server on :8080" when ready

## Quick Test

### Using cURL
```bash
# Create a reward
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "rewarded_at": "2025-11-10T10:00:00Z",
    "idempotency_key": "test-reward-001"
  }'

# Get portfolio
curl http://localhost:8080/api/v1/portfolio/user123
```

### Using Postman
1. Import `Stocky.postman_collection.json`
2. Run requests from the collection

## Troubleshooting

### Database Connection Failed
```
Error: Failed to connect to database
```
**Solution**: 
- Check PostgreSQL is running: `pg_isready`
- Verify .env credentials
- Ensure database exists: `psql -U postgres -l | grep assignment`

### Swagger Not Found
```
404 page not found
```
**Solution**: Swagger has been removed from this version. Use Postman collection or cURL for testing.

### Port Already in Use
```
Error: bind: address already in use
```
**Solution**: Change `SERVER_PORT` in .env or kill existing process:
```bash
lsof -ti:8080 | xargs kill -9
```

## Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/reward` | Create stock reward |
| GET | `/api/v1/today-stocks/{userId}` | Today's rewards |
| GET | `/api/v1/historical-inr/{userId}` | Historical INR values |
| GET | `/api/v1/stats/{userId}` | User statistics |
| GET | `/api/v1/portfolio/{userId}` | User portfolio |

## Next Steps

- ✅ Test all endpoints using Postman collection
- ✅ Import Postman collection for testing
- ✅ Check README.md for detailed documentation
- ✅ Review database schema and edge case handling

## Support

For issues or questions, please refer to README.md or check the code comments.
