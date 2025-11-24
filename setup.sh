#!/bin/bash

# Stocky Backend Initialization Script
# This script helps set up the project quickly

echo "========================================="
echo "    Stocky Backend Setup Script"
echo "========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Download from: https://go.dev/dl/"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null
then
    echo "⚠️  PostgreSQL client not found. Make sure PostgreSQL is installed."
    echo "   Download from: https://www.postgresql.org/download/"
else
    echo "✅ PostgreSQL client found"
fi

echo ""
echo "Step 1: Installing Go dependencies..."
go mod download
if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

echo ""
echo "Step 2: Checking environment configuration..."
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Copying from .env.example..."
    cp .env.example .env
    echo "✅ .env file created"
    echo "⚠️  Please edit .env file with your database credentials"
    echo ""
    echo "   Edit the following line in .env:"
    echo "   DATABASE_URL=postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/assignment?sslmode=disable"
    echo ""
    read -p "Press Enter after updating .env file..."
else
    echo "✅ .env file exists"
fi

echo ""
echo "Step 3: Testing database connection..."
echo "Please enter your PostgreSQL password when prompted"
echo ""

# Extract database credentials from .env
DB_URL=$(grep DATABASE_URL .env | cut -d '=' -f2-)

# Check if database exists
DB_EXISTS=$(psql -U postgres -lqt | cut -d \| -f 1 | grep -w assignment | wc -l)

if [ $DB_EXISTS -eq 0 ]; then
    echo "Creating database 'assignment'..."
    psql -U postgres -c "CREATE DATABASE assignment;"
    if [ $? -eq 0 ]; then
        echo "✅ Database created successfully"
    else
        echo "❌ Failed to create database"
        echo "Please create it manually: psql -U postgres -c 'CREATE DATABASE assignment;'"
    fi
else
    echo "✅ Database 'assignment' already exists"
fi

echo ""
echo "Step 4: Building the application..."
go build -o stocky-backend main.go
if [ $? -eq 0 ]; then
    echo "✅ Application built successfully"
else
    echo "❌ Build failed. Check the error messages above."
    exit 1
fi

echo ""
echo "========================================="
echo "    Setup Complete! 🎉"
echo "========================================="
echo ""
echo "To start the server, run:"
echo "  ./stocky-backend"
echo ""
echo "Or use:"
echo "  go run main.go"
echo ""
echo "The server will start on: http://localhost:8080"
echo ""
echo "Test the API:"
echo "  curl http://localhost:8080/api/health"
echo ""
echo "For more information, see:"
echo "  - README.md (comprehensive documentation)"
echo "  - QUICKSTART.md (quick start guide)"
echo "  - API_TESTING.md (API examples)"
echo ""
echo "Happy coding! 🚀"
