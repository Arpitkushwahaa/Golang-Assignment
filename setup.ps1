# Stocky Backend Initialization Script (PowerShell)
# This script helps set up the project quickly on Windows

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "    Stocky Backend Setup Script" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    Write-Host "   Download from: https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Check if PostgreSQL is installed
try {
    $psqlVersion = psql --version
    Write-Host "✅ PostgreSQL client found" -ForegroundColor Green
} catch {
    Write-Host "⚠️  PostgreSQL client not found. Make sure PostgreSQL is installed." -ForegroundColor Yellow
    Write-Host "   Download from: https://www.postgresql.org/download/" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 1: Installing Go dependencies..." -ForegroundColor Cyan
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Dependencies installed successfully" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to install dependencies" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Step 2: Checking environment configuration..." -ForegroundColor Cyan
if (-Not (Test-Path .env)) {
    Write-Host "⚠️  .env file not found. Copying from .env.example..." -ForegroundColor Yellow
    Copy-Item .env.example .env
    Write-Host "✅ .env file created" -ForegroundColor Green
    Write-Host "⚠️  Please edit .env file with your database credentials" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "   Edit the following line in .env:" -ForegroundColor Yellow
    Write-Host "   DATABASE_URL=postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/assignment?sslmode=disable" -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Press Enter after updating .env file"
} else {
    Write-Host "✅ .env file exists" -ForegroundColor Green
}

Write-Host ""
Write-Host "Step 3: Database setup..." -ForegroundColor Cyan
Write-Host "To create the database, run the following command in another terminal:" -ForegroundColor Yellow
Write-Host 'psql -U postgres -c "CREATE DATABASE assignment;"' -ForegroundColor White
Write-Host ""
$dbCreated = Read-Host "Have you created the database? (y/n)"
if ($dbCreated -eq 'y' -or $dbCreated -eq 'Y') {
    Write-Host "✅ Database setup confirmed" -ForegroundColor Green
} else {
    Write-Host "⚠️  Please create the database before running the application" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 4: Building the application..." -ForegroundColor Cyan
go build -o stocky-backend.exe main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Application built successfully" -ForegroundColor Green
} else {
    Write-Host "❌ Build failed. Check the error messages above." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "    Setup Complete! 🎉" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "To start the server, run:" -ForegroundColor Yellow
Write-Host "  .\stocky-backend.exe" -ForegroundColor White
Write-Host ""
Write-Host "Or use:" -ForegroundColor Yellow
Write-Host "  go run main.go" -ForegroundColor White
Write-Host ""
Write-Host "The server will start on: http://localhost:8080" -ForegroundColor Cyan
Write-Host ""
Write-Host "Test the API:" -ForegroundColor Yellow
Write-Host "  curl http://localhost:8080/api/health" -ForegroundColor White
Write-Host ""
Write-Host "For more information, see:" -ForegroundColor Yellow
Write-Host "  - README.md (comprehensive documentation)" -ForegroundColor White
Write-Host "  - QUICKSTART.md (quick start guide)" -ForegroundColor White
Write-Host "  - API_TESTING.md (API examples)" -ForegroundColor White
Write-Host ""
Write-Host "Happy coding! 🚀" -ForegroundColor Green
