#!/bin/bash

# Configuration script for dṛkka backend server - Development
# Usage: source config_dev.sh (to set environment variables)
#        ./config_dev.sh (to set variables and start server)

# Development Configuration for local testing
export PORT=8080
export DB_PATH=./drkka.db
export STATIC_DIR=../frontend/
export ALLOWED_ORIGINS="http://localhost:8080,http://127.0.0.1:8080,http://localhost:3000,http://127.0.0.1:3000"

echo "✅ Development environment variables configured:"
echo "   PORT=$PORT"
echo "   DB_PATH=$DB_PATH"
echo "   STATIC_DIR=$STATIC_DIR"
echo "   ALLOWED_ORIGINS=$ALLOWED_ORIGINS"
echo ""

# If script is executed (not sourced), start the server
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "🚀 Starting drkka server (development mode)..."
    ./drkka-server
fi
