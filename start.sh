#!/bin/bash
set -e

echo "🚀 ===== STARTING BOT ====="
echo "📁 Current directory: $(pwd)"
echo "📁 Files in directory:"
ls -la

echo "🔧 Building application..."
go build -o bot ./cmd/bot
ls -la bot

echo "✅ Binary built successfully"
chmod +x bot

echo "🗃️ Checking database..."
if [ -n "$DB_URL" ]; then
    echo "📦 DB_URL is set, running migrations..."
    ./bot migrate
else
    echo "❌ DB_URL is NOT set!"
fi

echo "🤖 Starting bot application..."
# Запускаем с подробным выводом
exec ./bot