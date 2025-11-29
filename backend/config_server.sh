#!/bin/bash

# Configuration script for dṛkka backend server
# Usage: source config_server.sh (to set environment variables)
#        ./config_server.sh (to set variables and start server)

# Production Configuration for codekaryashala.com
export PORT=8080
export DB_PATH=/var/lib/drkka/submissions.db
export STATIC_DIR=../
export ALLOWED_ORIGINS="http://codekaryashala.com,https://codekaryashala.com"

echo "✅ Environment variables configured:"
echo "   PORT=$PORT"
echo "   DB_PATH=$DB_PATH"
echo "   STATIC_DIR=$STATIC_DIR"
echo "   ALLOWED_ORIGINS=$ALLOWED_ORIGINS"
echo ""

# If script is executed (not sourced), start the server
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "🚀 Starting drkka server..."
    ./drkka-server
fi
