# Project Structure

Complete overview of the Stocky Backend project organization.

## Directory Tree

```
stocky-backend/
│
├── main.go                           # Application entry point
├── go.mod                            # Go module dependencies
├── go.sum                            # Dependency checksums
├── .env                              # Environment variables (not in git)
├── .env.example                      # Example environment config
├── .gitignore                        # Git ignore rules
│
├── controllers/                      # HTTP request handlers
│   └── reward_controller.go          # Reward API endpoints
│
├── services/                         # Business logic layer
│   ├── reward_service.go             # Reward management logic
│   ├── ledger_service.go             # Ledger operations
│   └── price_service.go              # Price management
│
├── models/                           # Database models (GORM)
│   └── models.go                     # All data models
│
├── db/                               # Database layer
│   └── database.go                   # DB connection & migrations
│
├── routes/                           # API route definitions
│   └── routes.go                     # Router configuration
│
├── utils/                            # Helper functions & middleware
│   ├── time.go                       # Time manipulation helpers
│   ├── price.go                      # Price calculations & validation
│   └── middleware.go                 # Logging & error handling
│
├── docs/                             # Documentation (this structure)
│   ├── README.md                     # Main documentation
│   ├── QUICKSTART.md                 # Quick start guide
│   ├── API_TESTING.md                # API testing guide
│   ├── DATABASE_SETUP.md             # Database setup guide
│   ├── DEPLOYMENT.md                 # Deployment guide
│   ├── EDGE_CASES.md                 # Edge cases documentation
│   └── PROJECT_STRUCTURE.md          # This file
│
└── Stocky_Postman_Collection.json    # Postman API collection

```

## File Descriptions

### Root Level Files

#### main.go
- Application entry point
- Server initialization
- Graceful shutdown handling
- Price update scheduler
- Logger configuration

**Key Functions:**
- `main()` - Starts the application
- `initLogger()` - Configures logging
- `startPriceUpdateScheduler()` - Schedules hourly price updates

#### go.mod
- Go module definition
- Dependency management
- Module name: `stocky-backend`

**Key Dependencies:**
- `gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `sirupsen/logrus` - Logging
- `shopspring/decimal` - Decimal precision

#### .env / .env.example
Environment configuration files.

**Variables:**
- `DATABASE_URL` - PostgreSQL connection string
- `SERVER_PORT` - HTTP server port
- `GIN_MODE` - debug/release mode
- `PRICE_API_URL` - Price service URL
- `STOCKS` - Comma-separated stock symbols

### Controllers Layer

#### controllers/reward_controller.go
HTTP request handlers for all reward endpoints.

**Endpoints:**
- `CreateReward()` - POST /api/reward
- `GetTodayStocks()` - GET /api/today-stocks/:userId
- `GetHistoricalINR()` - GET /api/historical-inr/:userId
- `GetStats()` - GET /api/stats/:userId
- `GetPortfolio()` - GET /api/portfolio/:userId
- `HealthCheck()` - GET /api/health

**Responsibilities:**
- Request validation
- JSON parsing
- Response formatting
- Error handling

### Services Layer

#### services/reward_service.go
Core business logic for reward management.

**Key Functions:**
- `CreateReward()` - Creates reward with validation
- `GetTodayRewards()` - Fetches today's rewards
- `GetHistoricalINR()` - Calculates historical valuations
- `GetStats()` - Compiles user statistics
- `GetPortfolio()` - Builds portfolio view

**Responsibilities:**
- Business rule enforcement
- Transaction management
- Data aggregation
- Edge case handling

#### services/ledger_service.go
Double-entry ledger operations.

**Key Functions:**
- `CreateLedgerEntries()` - Creates STOCK/CASH/FEE entries
- `GetUserStockHoldings()` - Calculates total holdings
- `GetUserStockHoldingsUpToDate()` - Historical holdings
- `CreateReversalEntries()` - Reward cancellation

**Responsibilities:**
- Ledger integrity
- Fee calculations
- Holdings aggregation
- Reversal logic

#### services/price_service.go
Stock price management.

**Key Functions:**
- `GetCurrentPrice()` - Fetches latest price
- `GetPriceAtTime()` - Historical price lookup
- `SavePrice()` - Stores price in database
- `UpdateAllPrices()` - Generates prices for all stocks
- `GetPricesForDate()` - Date-specific prices

**Responsibilities:**
- Price generation (mock)
- Price caching
- Fallback handling
- Hourly updates

### Models Layer

#### models/models.go
GORM database models.

**Models:**
- `RewardEvent` - Stock reward records
- `LedgerEntry` - Double-entry ledger
- `PriceHistory` - Historical stock prices
- `StockConfig` - Stock configuration (splits, etc.)

**Features:**
- Timestamps (CreatedAt, UpdatedAt)
- Soft deletes (DeletedAt)
- Relationships (foreign keys)
- Custom table names

### Database Layer

#### db/database.go
Database connection and migrations.

**Key Functions:**
- `Initialize()` - Connects to PostgreSQL
- `runMigrations()` - Auto-creates tables
- `createIndexes()` - Composite indexes
- `initializeStockConfigs()` - Seeds stock data
- `Close()` - Closes connection

**Responsibilities:**
- Connection pooling
- Auto-migration
- Index management
- Initial data seeding

### Routes Layer

#### routes/routes.go
API route configuration.

**Features:**
- Route grouping (/api)
- Middleware setup
- CORS configuration
- Service initialization

**Middleware:**
- Recovery (panic handling)
- ErrorHandler (custom errors)
- LoggingMiddleware (request logging)
- CORS (cross-origin requests)

### Utils Layer

#### utils/time.go
Time manipulation helpers.

**Functions:**
- `StartOfDay()` / `EndOfDay()` - Day boundaries
- `IsToday()` - Date comparison
- `GetDateString()` - Date formatting
- `GetPastDates()` - Date range generation
- `NowUTC()` - Current UTC time

#### utils/price.go
Price calculations and validation.

**Functions:**
- `PriceGenerator` - Mock price generation
- `CalculateFees()` - Brokerage/STT/GST
- `ValidateQuantity()` - Input validation
- `ValidateStockSymbol()` - Symbol validation
- `RoundINR()` / `RoundQuantity()` - Precision rounding

#### utils/middleware.go
HTTP middleware.

**Middleware:**
- `LoggingMiddleware()` - Logs all requests
- `ErrorHandler()` - Panic recovery

## Data Flow

### Request Flow
```
Client Request
    ↓
Middleware (Logging, CORS, Error Handler)
    ↓
Router (routes/routes.go)
    ↓
Controller (controllers/reward_controller.go)
    ↓
Service (services/*.go)
    ↓
Database (db/database.go)
    ↓
Response (JSON)
```

### Reward Creation Flow
```
POST /api/reward
    ↓
RewardController.CreateReward()
    ↓
RewardService.CreateReward()
    ├── Validate inputs
    ├── Check duplicates
    ├── Get current price (PriceService)
    ├── Begin transaction
    ├── Create RewardEvent
    ├── Create LedgerEntries (LedgerService)
    └── Commit transaction
    ↓
Return success/error
```

### Portfolio Calculation Flow
```
GET /api/portfolio/:userId
    ↓
RewardController.GetPortfolio()
    ↓
RewardService.GetPortfolio()
    ├── Get holdings (LedgerService)
    │   └── Sum STOCK entries per symbol
    ├── For each stock:
    │   ├── Get current price (PriceService)
    │   └── Calculate value = price × quantity
    └── Return portfolio summary
    ↓
Return JSON response
```

## Database Schema Relationships

```
reward_events
    ├── id (PK)
    ├── user_id
    ├── stock_symbol
    ├── quantity
    └── timestamp
        │
        ├─────────────────────┐
        │                     │
        ▼                     ▼
ledger_entries          stock_config
    ├── id (PK)             ├── id (PK)
    ├── reward_event_id (FK)├── stock_symbol (UNIQUE)
    ├── entry_type          ├── multiplier
    ├── stock_symbol        └── is_active
    ├── quantity
    ├── amount_inr
    └── timestamp

price_history
    ├── id (PK)
    ├── stock_symbol
    ├── price_inr
    └── timestamp
```

## Architectural Patterns

### Clean Architecture
- **Controllers**: Handle HTTP
- **Services**: Business logic
- **Models**: Data structures
- **Database**: Persistence layer

### Dependency Injection
Services are initialized and injected into controllers:
```go
priceService := services.NewPriceService()
ledgerService := services.NewLedgerService()
rewardService := services.NewRewardService(priceService, ledgerService)
rewardController := controllers.NewRewardController(rewardService)
```

### Transaction Management
All reward operations use database transactions for ACID compliance:
```go
tx := db.DB.Begin()
// operations...
tx.Commit()
```

### Error Handling
Layered error handling:
1. Validation errors → 400 Bad Request
2. Business rule violations → 409 Conflict
3. Database errors → 500 Internal Server Error

## Code Organization Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Single Responsibility**: Each file/function does one thing well
3. **DRY (Don't Repeat Yourself)**: Shared logic in utils
4. **Testability**: Services can be tested independently
5. **Maintainability**: Clear structure, good naming conventions

## Development Workflow

```bash
# 1. Make changes to code
vim services/reward_service.go

# 2. Run the application
go run main.go

# 3. Test with cURL/Postman
curl http://localhost:8080/api/health

# 4. Check logs
# (logs appear in console)

# 5. Commit changes
git add .
git commit -m "Add feature X"
git push
```

## Adding New Features

### Example: Add new endpoint

1. **Create model** (if needed): `models/models.go`
2. **Add service logic**: `services/new_service.go`
3. **Create controller**: `controllers/new_controller.go`
4. **Register route**: `routes/routes.go`
5. **Test endpoint**: Use Postman/cURL
6. **Update documentation**: `README.md`, `API_TESTING.md`

---

**Project structure complete! 📁**
