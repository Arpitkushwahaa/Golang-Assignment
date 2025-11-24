# Quick Start Guide

Get the Stocky Backend up and running in 5 minutes!

## Prerequisites

- Go 1.21+ installed
- PostgreSQL 12+ installed and running
- Git installed

## Step 1: Clone the Repository

```bash
git clone <repository-url>
cd Assignment
```

## Step 2: Set Up Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE assignment;

# Exit
\q
```

## Step 3: Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your database credentials
# Update the DATABASE_URL line with your username and password
```

**Example .env:**
```env
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/assignment?sslmode=disable
SERVER_PORT=8080
GIN_MODE=debug
```

## Step 4: Install Dependencies

```bash
go mod download
```

## Step 5: Run the Application

```bash
go run main.go
```

You should see:
```
INFO[...] Starting Stocky Backend...
INFO[...] Connecting to database...
INFO[...] Database connected successfully
INFO[...] Running database migrations...
INFO[...] Database migrations completed successfully
INFO[...] Starting price update scheduler (hourly)
INFO[...] Server is running on http://localhost:8080
```

## Step 6: Test the API

### Using cURL

```bash
# Health check
curl http://localhost:8080/api/health

# Create a reward
curl -X POST http://localhost:8080/api/reward \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "symbol": "RELIANCE",
    "quantity": 2.5,
    "timestamp": "2025-01-23T10:30:00Z"
  }'

# Get user portfolio
curl http://localhost:8080/api/portfolio/1
```

### Using Postman

1. Import `Stocky_Postman_Collection.json`
2. Run the collection to test all endpoints

## Step 7: View Database Tables

```bash
psql -U postgres -d assignment

# List tables
\dt

# View reward events
SELECT * FROM reward_events;

# View ledger entries
SELECT * FROM ledger_entries;

# View price history
SELECT * FROM price_history ORDER BY timestamp DESC LIMIT 5;
```

## Troubleshooting

### Can't connect to database
```
Error: failed to connect to database
```
**Solution:** Check your DATABASE_URL in .env and ensure PostgreSQL is running.

### Port already in use
```
Error: listen tcp :8080: bind: address already in use
```
**Solution:** Change SERVER_PORT in .env or stop the process using port 8080.

### Go dependencies error
```
Error: cannot find package
```
**Solution:** Run `go mod download` and `go mod tidy`

## Next Steps

- Read the [README.md](README.md) for detailed documentation
- Check [API_TESTING.md](API_TESTING.md) for comprehensive API examples
- Review [EDGE_CASES.md](EDGE_CASES.md) for edge case handling
- See [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment

## Quick Commands Reference

```bash
# Run the server
go run main.go

# Build binary
go build -o stocky-backend main.go

# Run binary
./stocky-backend

# Run with custom port
SERVER_PORT=3000 go run main.go

# View logs (if running in background)
tail -f /var/log/stocky/stdout.log
```

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/reward | Create a reward |
| GET | /api/today-stocks/:userId | Get today's rewards |
| GET | /api/historical-inr/:userId | Get historical valuations |
| GET | /api/stats/:userId | Get user statistics |
| GET | /api/portfolio/:userId | Get user portfolio |
| GET | /api/health | Health check |

## Sample Test Data

```bash
# User 1 - Multiple stocks
curl -X POST http://localhost:8080/api/reward -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "RELIANCE", "quantity": 2.5, "timestamp": "2025-01-23T10:30:00Z"}'

curl -X POST http://localhost:8080/api/reward -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "TCS", "quantity": 1.5, "timestamp": "2025-01-23T14:15:00Z"}'

curl -X POST http://localhost:8080/api/reward -H "Content-Type: application/json" \
  -d '{"userId": 1, "symbol": "INFY", "quantity": 3.25, "timestamp": "2025-01-23T16:45:00Z"}'

# User 2
curl -X POST http://localhost:8080/api/reward -H "Content-Type: application/json" \
  -d '{"userId": 2, "symbol": "HDFCBANK", "quantity": 5.0, "timestamp": "2025-01-23T11:00:00Z"}'

# Check results
curl http://localhost:8080/api/stats/1
curl http://localhost:8080/api/portfolio/1
```

---

**Happy Coding! 🚀**

Need help? Check the full documentation in README.md or create an issue on GitHub.
