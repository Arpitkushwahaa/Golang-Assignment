# Deployment Guide

This guide covers deploying the Stocky Backend to production environments.

## Deployment Checklist

- [ ] PostgreSQL database set up
- [ ] Environment variables configured
- [ ] Application tested locally
- [ ] Build successful
- [ ] Database migrations run
- [ ] API endpoints tested
- [ ] Monitoring configured

## Building for Production

### 1. Build Binary

```bash
# Build for current platform
go build -o stocky-backend main.go

# Build for Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o stocky-backend-linux main.go

# Build for Windows (from any OS)
GOOS=windows GOARCH=amd64 go build -o stocky-backend.exe main.go

# Build for macOS (from any OS)
GOOS=darwin GOARCH=amd64 go build -o stocky-backend-macos main.go
```

### 2. Set Production Environment

Update `.env` for production:

```env
# Production Database
DATABASE_URL=postgres://prod_user:secure_password@db-host:5432/assignment?sslmode=require

# Production Server
SERVER_PORT=8080
GIN_MODE=release

# Price API
PRICE_API_URL=https://api.stocky.com/prices

# Stock Symbols
STOCKS=RELIANCE,TCS,INFY,HDFCBANK,ICICIBANK,SBIN,BHARTIARTL,ITC,KOTAKBANK,LT
```

### 3. Run Application

```bash
# Direct execution
./stocky-backend

# With custom port
SERVER_PORT=3000 ./stocky-backend

# Background process
nohup ./stocky-backend > /var/log/stocky.log 2>&1 &
```

## Docker Deployment

### Dockerfile

Create `Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stocky-backend main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/stocky-backend .

# Copy .env file (optional, prefer environment variables)
# COPY .env .

EXPOSE 8080

CMD ["./stocky-backend"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: assignment
      POSTGRES_USER: stocky
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U stocky"]
      interval: 10s
      timeout: 5s
      retries: 5

  stocky-backend:
    build: .
    environment:
      DATABASE_URL: postgres://stocky:${DB_PASSWORD}@postgres:5432/assignment?sslmode=disable
      SERVER_PORT: 8080
      GIN_MODE: release
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:
```

### Build and Run with Docker

```bash
# Build image
docker build -t stocky-backend:latest .

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f stocky-backend

# Stop
docker-compose down
```

## Cloud Deployment

### AWS EC2

#### 1. Launch EC2 Instance
- AMI: Ubuntu 22.04 LTS
- Instance Type: t3.micro (for testing) or t3.medium (production)
- Security Group: Allow ports 22 (SSH), 8080 (HTTP), 5432 (PostgreSQL)

#### 2. Install Dependencies

```bash
# Connect via SSH
ssh -i your-key.pem ubuntu@ec2-instance-ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database
sudo -u postgres psql
CREATE DATABASE assignment;
CREATE USER stocky WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE assignment TO stocky;
\q
```

#### 3. Deploy Application

```bash
# Clone repository
git clone <your-repo-url>
cd Assignment

# Configure environment
cp .env.example .env
nano .env  # Update DATABASE_URL

# Build
go build -o stocky-backend main.go

# Run
./stocky-backend
```

#### 4. Setup Systemd Service

Create `/etc/systemd/system/stocky-backend.service`:

```ini
[Unit]
Description=Stocky Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/Assignment
ExecStart=/home/ubuntu/Assignment/stocky-backend
Restart=on-failure
RestartSec=5s
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable stocky-backend
sudo systemctl start stocky-backend
sudo systemctl status stocky-backend
```

### AWS RDS (PostgreSQL)

#### 1. Create RDS Instance
- Engine: PostgreSQL 15
- Instance Class: db.t3.micro (free tier)
- Storage: 20 GB
- Public Access: Yes (for testing) / No (production with VPC)

#### 2. Update Connection String

```env
DATABASE_URL=postgres://username:password@your-rds-endpoint.amazonaws.com:5432/assignment?sslmode=require
```

### Heroku Deployment

#### 1. Prepare Application

Create `Procfile`:
```
web: ./stocky-backend
```

Create `heroku.yml`:
```yaml
build:
  languages:
    - go
```

#### 2. Deploy

```bash
# Login to Heroku
heroku login

# Create app
heroku create stocky-backend-app

# Add PostgreSQL addon
heroku addons:create heroku-postgresql:mini

# Set environment variables
heroku config:set GIN_MODE=release
heroku config:set SERVER_PORT=$PORT

# Deploy
git push heroku main

# View logs
heroku logs --tail
```

### Google Cloud Run

#### 1. Build Container

```bash
# Build for Cloud Run
gcloud builds submit --tag gcr.io/your-project-id/stocky-backend

# Or use Docker
docker build -t gcr.io/your-project-id/stocky-backend .
docker push gcr.io/your-project-id/stocky-backend
```

#### 2. Deploy

```bash
gcloud run deploy stocky-backend \
  --image gcr.io/your-project-id/stocky-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DATABASE_URL="postgres://..." \
  --set-env-vars GIN_MODE=release
```

### DigitalOcean App Platform

#### 1. Create App

```bash
# Install doctl
brew install doctl  # macOS
sudo snap install doctl  # Linux

# Authenticate
doctl auth init

# Create app from spec
doctl apps create --spec app.yaml
```

#### app.yaml

```yaml
name: stocky-backend
region: nyc
services:
  - name: api
    github:
      repo: your-username/stocky-backend
      branch: main
    build_command: go build -o stocky-backend main.go
    run_command: ./stocky-backend
    environment_slug: go
    http_port: 8080
    instance_count: 1
    instance_size_slug: basic-xxs
    envs:
      - key: GIN_MODE
        value: release
      - key: DATABASE_URL
        value: ${db.DATABASE_URL}

databases:
  - name: db
    engine: PG
    version: "15"
```

## Reverse Proxy (Nginx)

### Install Nginx

```bash
sudo apt install nginx -y
```

### Configure Nginx

Create `/etc/nginx/sites-available/stocky-backend`:

```nginx
server {
    listen 80;
    server_name api.stocky.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable and restart:
```bash
sudo ln -s /etc/nginx/sites-available/stocky-backend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain certificate
sudo certbot --nginx -d api.stocky.com

# Auto-renewal
sudo certbot renew --dry-run
```

## Monitoring

### Setup Logging

```bash
# Create log directory
sudo mkdir -p /var/log/stocky

# Update systemd service to log
[Service]
StandardOutput=append:/var/log/stocky/stdout.log
StandardError=append:/var/log/stocky/stderr.log
```

### Log Rotation

Create `/etc/logrotate.d/stocky-backend`:

```
/var/log/stocky/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 ubuntu ubuntu
    sharedscripts
}
```

### Health Monitoring

```bash
# Cron job to check health
*/5 * * * * curl -f http://localhost:8080/api/health || systemctl restart stocky-backend
```

## Backup Strategy

### Database Backup

```bash
# Create backup script
cat > /home/ubuntu/backup-db.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U stocky assignment > /backups/stocky_$DATE.sql
# Keep only last 7 days
find /backups -name "stocky_*.sql" -mtime +7 -delete
EOF

chmod +x /home/ubuntu/backup-db.sh

# Add to crontab (daily at 2 AM)
0 2 * * * /home/ubuntu/backup-db.sh
```

## Performance Optimization

### 1. Database Connection Pooling

Already configured in `db/database.go`:
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. Gin Release Mode

```env
GIN_MODE=release
```

### 3. PostgreSQL Tuning

```sql
-- Increase shared buffers
ALTER SYSTEM SET shared_buffers = '256MB';

-- Increase work memory
ALTER SYSTEM SET work_mem = '4MB';

-- Reload configuration
SELECT pg_reload_conf();
```

## Security Best Practices

1. **Use strong database passwords**
2. **Enable SSL for database connections** (`sslmode=require`)
3. **Restrict database access** (firewall rules)
4. **Use environment variables** for sensitive data
5. **Regular security updates** (`apt update && apt upgrade`)
6. **Implement rate limiting** (nginx or application level)
7. **Enable HTTPS** (Let's Encrypt)
8. **Use principle of least privilege** for database users

## Rollback Strategy

```bash
# Tag releases
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0

# Rollback to previous version
git checkout v1.0.0
go build -o stocky-backend main.go
sudo systemctl restart stocky-backend
```

## Scaling Considerations

### Horizontal Scaling
- Deploy multiple instances behind load balancer
- Use managed database (RDS, Cloud SQL)
- Implement Redis for caching

### Vertical Scaling
- Increase instance size (CPU, RAM)
- Optimize database queries
- Add read replicas

---

**Deployment Complete! 🚀**
