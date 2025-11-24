# Stocky Backend - Stock Reward Management System

A production-ready Golang backend for managing stock rewards, where users earn shares of Indian stocks (Reliance, TCS, Infosys, etc.) as incentives. The system tracks rewards, maintains a double-entry ledger, and provides portfolio valuation features.

> 📚 **[View Documentation Index](INDEX.md)** - Complete guide to all documentation files

## 🏗 Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
stocky-backend/
├── controllers/       # HTTP request handlers
├── services/          # Business logic layer
├── models/            # Database models (GORM)
├── db/                # Database connection & migrations
├── routes/            # API route definitions
├── utils/             # Helper functions & middleware
├── main.go            # Application entry point
├── .env               # Environment variables
├── .env.example       # Example environment configuration
├── go.mod             # Go module dependencies
└── README.md          # This file
```

## 🚀 Features

### Core Features
- ✅ **Reward System**: Record stock rewards for users
- ✅ **Double-Entry Ledger**: Track stock units, cash outflow, and fees
- ✅ **Price Service**: Hourly mock price updates for Indian stocks
- ✅ **Portfolio Management**: Track user holdings and valuations
- ✅ **Historical Data**: Daily INR valuations up to yesterday

### Edge Case Handling
- ✅ **Deduplication**: Prevents identical reward entries
- ✅ **Stock Splits**: Configuration table with multipliers
- ✅ **Rounding Precision**: INR (4 decimals), Shares (6 decimals)
- ✅ **Price Fallback**: Uses last known price if service is down
- ✅ **Reward Reversal**: Ledger reversal entries support

## 📋 Prerequisites

- **Go** 1.21 or higher
- **PostgreSQL** 12 or higher
- **Git**

## 🔧 Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Assignment
```

### 2. Set Up PostgreSQL Database

Create a database named `assignment`:

```sql
CREATE DATABASE assignment;
```

### 3. Configure Environment Variables

Copy the example environment file and update with your database credentials:

```bash
cp .env.example .env
```

Edit `.env`:

```env
DATABASE_URL=postgres://your_username:your_password@localhost:5432/assignment?sslmode=disable
SERVER_PORT=8080
GIN_MODE=debug
PRICE_API_URL=http://localhost:8080/api/prices
STOCKS=RELIANCE,TCS,INFY,HDFCBANK,ICICIBANK,SBIN,BHARTIARTL,ITC,KOTAKBANK,LT
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Base URL
```
http://localhost:8080/api
```

### 1. **POST /api/reward** - Create Reward

Records that a user has been rewarded X shares of a stock.

**Request:**
```json
{
  "userId": 1,
  "symbol": "RELIANCE",
  "quantity": 2.5,
  "timestamp": "2025-01-23T10:30:00Z"
}
```

**Response:**
```json
{
  "success": true
}
```

**Error (Duplicate):**
```json
{
  "success": false,
  "error": "duplicate reward: identical reward already exists"
}
```

---

### 2. **GET /api/today-stocks/:userId** - Today's Stock Rewards

Returns all reward events for the user for TODAY only.

**Example:** `GET /api/today-stocks/1`

**Response:**
```json
{
  "userId": 1,
  "rewards": [
    {
      "symbol": "RELIANCE",
      "quantity": "2.5",
      "timestamp": "2025-01-23T10:30:00Z"
    },
    {
      "symbol": "INFY",
      "quantity": "1.2",
      "timestamp": "2025-01-23T14:15:00Z"
    }
  ]
}
```

---

### 3. **GET /api/historical-inr/:userId** - Historical INR Valuation

Returns INR valuation per past day (up to yesterday).

**Example:** `GET /api/historical-inr/1`

**Response:**
```json
{
  "userId": 1,
  "days": [
    {
      "date": "2025-01-21",
      "valueINR": "1500.4500"
    },
    {
      "date": "2025-01-20",
      "valueINR": "1320.1000"
    }
  ]
}
```

---

### 4. **GET /api/stats/:userId** - User Statistics

Returns total shares rewarded today (grouped by stock) and current portfolio value.

**Example:** `GET /api/stats/1`

**Response:**
```json
{
  "userId": 1,
  "todayRewards": {
    "RELIANCE": "3.0",
    "TCS": "1.5"
  },
  "portfolioValueINR": "48221.6500"
}
```

---

### 5. **GET /api/portfolio/:userId** - User Portfolio (BONUS)

Shows full holdings grouped by stock with current INR value.

**Example:** `GET /api/portfolio/1`

**Response:**
```json
{
  "userId": 1,
  "holdings": [
    {
      "symbol": "RELIANCE",
      "quantity": "5.5",
      "currentPrice": "2450.5000",
      "currentValue": "13477.7500"
    },
    {
      "symbol": "TCS",
      "quantity": "3.2",
      "currentPrice": "3680.7500",
      "currentValue": "11778.4000"
    }
  ],
  "totalValue": "25256.1500"
}
```

---

### 6. **GET /api/health** - Health Check

**Response:**
```json
{
  "status": "healthy",
  "time": "2025-01-23T10:30:00Z"
}
```

## 🗄️ Database Schema

### 1. **reward_events**
Stores all reward events.

| Column       | Type            | Description                    |
|--------------|-----------------|--------------------------------|
| id           | SERIAL          | Primary key                    |
| user_id      | INTEGER         | User identifier                |
| stock_symbol | VARCHAR(20)     | Stock ticker symbol            |
| quantity     | NUMERIC(18,6)   | Number of shares (fractional)  |
| timestamp    | TIMESTAMPTZ     | Reward timestamp               |
| created_at   | TIMESTAMPTZ     | Record creation time           |
| updated_at   | TIMESTAMPTZ     | Record update time             |
| deleted_at   | TIMESTAMPTZ     | Soft delete timestamp          |

**Indexes:**
- `idx_user_rewards` on `(user_id)`
- `idx_stock_symbol` on `(stock_symbol)`
- `idx_timestamp` on `(timestamp)`
- `idx_reward_dedup` unique on `(user_id, stock_symbol, quantity, timestamp)`

---

### 2. **ledger_entries**
Double-entry ledger for accounting.

| Column          | Type            | Description                      |
|-----------------|-----------------|----------------------------------|
| id              | SERIAL          | Primary key                      |
| reward_event_id | INTEGER         | Foreign key to reward_events     |
| entry_type      | VARCHAR(10)     | STOCK, CASH, or FEE              |
| stock_symbol    | VARCHAR(20)     | Stock ticker (nullable)          |
| quantity        | NUMERIC(18,6)   | Share quantity (for STOCK type)  |
| amount_inr      | NUMERIC(18,4)   | INR amount                       |
| timestamp       | TIMESTAMPTZ     | Entry timestamp                  |
| created_at      | TIMESTAMPTZ     | Record creation time             |

**Entry Types:**
- `STOCK`: User receives shares (+quantity)
- `CASH`: Company pays for shares (-amount_inr)
- `FEE`: Company pays brokerage/taxes (-amount_inr)

---

### 3. **price_history**
Historical stock prices.

| Column       | Type            | Description           |
|--------------|-----------------|-----------------------|
| id           | SERIAL          | Primary key           |
| stock_symbol | VARCHAR(20)     | Stock ticker symbol   |
| price_inr    | NUMERIC(18,4)   | Price in INR          |
| timestamp    | TIMESTAMPTZ     | Price timestamp       |
| created_at   | TIMESTAMPTZ     | Record creation time  |

**Indexes:**
- `idx_symbol_time` on `(stock_symbol, timestamp DESC)`

---

### 4. **stock_config**
Stock configuration (splits, multipliers).

| Column       | Type            | Description                    |
|--------------|-----------------|--------------------------------|
| id           | SERIAL          | Primary key                    |
| stock_symbol | VARCHAR(20)     | Stock ticker (unique)          |
| multiplier   | NUMERIC(18,6)   | Split multiplier (default: 1)  |
| is_active    | BOOLEAN         | Active status                  |
| notes        | TEXT            | Configuration notes            |
| created_at   | TIMESTAMPTZ     | Record creation time           |
| updated_at   | TIMESTAMPTZ     | Record update time             |

## 🔄 Business Logic

### Reward Flow

When `POST /api/reward` is called:

1. **Validation**: Check userId, symbol, quantity, and timestamp
2. **Deduplication**: Reject if identical reward exists
3. **Price Lookup**: Fetch current/historical stock price
4. **Create Reward Event**: Insert into `reward_events`
5. **Create Ledger Entries**:
   - **STOCK**: Credit shares to user
   - **CASH**: Debit company account (stock value)
   - **FEE**: Debit company account (brokerage + STT + GST)
6. **Transaction Commit**: All-or-nothing database transaction

### Fee Calculation

```
Brokerage = min(0.03% of transaction value, ₹20)
STT       = 0.1% of transaction value
GST       = 18% of brokerage
Total Fee = Brokerage + STT + GST
```

### Price Service

- **Initial Prices**: Seeded on startup for 10 Indian stocks
- **Hourly Updates**: Scheduled task generates new prices (±5% variation)
- **Fallback**: If price unavailable, uses last known price from database
- **Storage**: All prices stored in `price_history`

## 🛠️ Tech Stack

| Technology           | Purpose                          |
|----------------------|----------------------------------|
| **Go 1.21**          | Programming language             |
| **Gin**              | HTTP web framework               |
| **GORM**             | ORM for database operations      |
| **PostgreSQL**       | Relational database              |
| **Logrus**           | Structured logging               |
| **godotenv**         | Environment variable management  |
| **shopspring/decimal** | Precise decimal calculations   |

## 📦 Dependencies

```go
github.com/gin-gonic/gin          // Web framework
github.com/sirupsen/logrus        // Logging
github.com/joho/godotenv          // .env file support
github.com/shopspring/decimal     // Decimal precision
gorm.io/gorm                      // ORM
gorm.io/driver/postgres           // PostgreSQL driver
github.com/gin-contrib/cors       // CORS middleware
```

## 🧪 Testing

### Manual Testing with cURL

**Create a reward:**
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

**Get today's stocks:**
```bash
curl http://localhost:8080/api/today-stocks/1
```

**Get portfolio:**
```bash
curl http://localhost:8080/api/portfolio/1
```

### Using Postman

Import the included `Stocky_Postman_Collection.json` file into Postman for pre-configured API requests.

## 🔒 Security Considerations

- **SQL Injection**: GORM parameterized queries prevent SQL injection
- **Input Validation**: All inputs validated before processing
- **Error Handling**: Sensitive errors not exposed to clients
- **CORS**: Configured to allow cross-origin requests

## 📈 Scalability

### Current Optimizations
- Database connection pooling (max 100 connections)
- Composite indexes for fast queries
- Batch price updates
- Efficient date-range queries

### Future Enhancements
- Redis caching for prices and portfolio values
- Message queue for async reward processing
- Read replicas for historical data queries
- Rate limiting on API endpoints

## 🐛 Troubleshooting

### Database Connection Error
```
Error: failed to connect to database
```
**Solution**: Check DATABASE_URL in .env file and ensure PostgreSQL is running.

### Port Already in Use
```
Error: listen tcp :8080: bind: address already in use
```
**Solution**: Change SERVER_PORT in .env or kill the process using port 8080.

### Price Service Error
```
Error: failed to fetch price
```
**Solution**: Prices are auto-generated. If error persists, check database connectivity.

## 📝 License

This project is created for educational and assessment purposes.

## 👨‍💻 Author

Developed as part of the Stocky Backend Assignment.

## 🙏 Acknowledgments

- Gin Web Framework
- GORM ORM
- PostgreSQL Database
- Go Community

---

**Happy Coding! 🚀**
