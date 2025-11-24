# Database Setup Guide

## Prerequisites

- PostgreSQL 12 or higher installed
- Access to PostgreSQL (username and password)

## Quick Setup

### 1. Install PostgreSQL

#### Windows
Download and install from: https://www.postgresql.org/download/windows/

#### macOS
```bash
brew install postgresql
brew services start postgresql
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 2. Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE assignment;

# Exit psql
\q
```

### 3. Update Connection String

Edit `.env` file with your PostgreSQL credentials:

```env
DATABASE_URL=postgres://your_username:your_password@localhost:5432/assignment?sslmode=disable
```

**Example:**
```env
DATABASE_URL=postgres://postgres:password123@localhost:5432/assignment?sslmode=disable
```

### 4. Verify Connection

Start the application:
```bash
go run main.go
```

You should see:
```
INFO[...] Connecting to database...
INFO[...] Database connected successfully
INFO[...] Running database migrations...
INFO[...] Database migrations completed successfully
```

## Database Schema

The application will automatically create the following tables:

### reward_events
```sql
CREATE TABLE reward_events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    stock_symbol VARCHAR(20) NOT NULL,
    quantity NUMERIC(18,6) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_user_rewards ON reward_events(user_id);
CREATE INDEX idx_stock_symbol ON reward_events(stock_symbol);
CREATE INDEX idx_timestamp ON reward_events(timestamp);
CREATE INDEX idx_reward_user_timestamp ON reward_events(user_id, timestamp DESC);
CREATE UNIQUE INDEX idx_reward_dedup ON reward_events(user_id, stock_symbol, quantity, timestamp) 
    WHERE deleted_at IS NULL;
```

### ledger_entries
```sql
CREATE TABLE ledger_entries (
    id SERIAL PRIMARY KEY,
    reward_event_id INTEGER NOT NULL REFERENCES reward_events(id),
    entry_type VARCHAR(10) NOT NULL,
    stock_symbol VARCHAR(20),
    quantity NUMERIC(18,6) NOT NULL DEFAULT 0,
    amount_inr NUMERIC(18,4) NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_reward_event ON ledger_entries(reward_event_id);
CREATE INDEX idx_ledger_entry_type ON ledger_entries(entry_type);
```

### price_history
```sql
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    stock_symbol VARCHAR(20) NOT NULL,
    price_inr NUMERIC(18,4) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_symbol_time ON price_history(stock_symbol, timestamp DESC);
CREATE INDEX idx_price_symbol_timestamp ON price_history(stock_symbol, timestamp DESC);
```

### stock_config
```sql
CREATE TABLE stock_config (
    id SERIAL PRIMARY KEY,
    stock_symbol VARCHAR(20) NOT NULL UNIQUE,
    multiplier NUMERIC(18,6) NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
```

## Troubleshooting

### Connection Refused
```
Error: failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Solution:**
1. Ensure PostgreSQL is running:
   ```bash
   # Check status
   sudo systemctl status postgresql  # Linux
   brew services list                 # macOS
   ```

2. Start PostgreSQL if not running:
   ```bash
   sudo systemctl start postgresql   # Linux
   brew services start postgresql    # macOS
   ```

### Authentication Failed
```
Error: pq: password authentication failed for user "postgres"
```

**Solution:**
1. Verify your PostgreSQL username and password
2. Update the DATABASE_URL in `.env` with correct credentials
3. Reset password if needed:
   ```sql
   ALTER USER postgres PASSWORD 'new_password';
   ```

### Database Does Not Exist
```
Error: pq: database "assignment" does not exist
```

**Solution:**
```sql
CREATE DATABASE assignment;
```

### Permission Denied
```
Error: pq: permission denied for database
```

**Solution:**
Grant necessary permissions:
```sql
GRANT ALL PRIVILEGES ON DATABASE assignment TO your_username;
```

## Manual Schema Creation (Optional)

If you prefer to create tables manually instead of auto-migration, run:

```sql
-- Connect to assignment database
\c assignment

-- Create tables (see schema above)
-- Then disable auto-migration in code
```

## Viewing Data

### Using psql
```bash
psql -U postgres -d assignment

-- List all tables
\dt

-- View reward events
SELECT * FROM reward_events;

-- View ledger entries
SELECT * FROM ledger_entries;

-- View price history
SELECT * FROM price_history ORDER BY timestamp DESC LIMIT 10;

-- View stock config
SELECT * FROM stock_config;
```

### Using pgAdmin
1. Download pgAdmin: https://www.pgadmin.org/download/
2. Connect to localhost:5432
3. Navigate to assignment database
4. Browse tables under Schemas > public > Tables

## Backup and Restore

### Backup
```bash
pg_dump -U postgres assignment > backup.sql
```

### Restore
```bash
psql -U postgres assignment < backup.sql
```

## Performance Tuning

### Analyze Tables
```sql
ANALYZE reward_events;
ANALYZE ledger_entries;
ANALYZE price_history;
```

### Check Index Usage
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;
```

## Resetting Database

To start fresh:

```bash
# Drop and recreate database
psql -U postgres

DROP DATABASE assignment;
CREATE DATABASE assignment;
\q

# Restart application to run migrations
go run main.go
```
