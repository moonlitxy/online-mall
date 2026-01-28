#!/bin/bash

# 在线商城后端启动脚本

echo "Starting Online Mall Backend..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# 检查依赖
echo "Checking dependencies..."
go mod tidy

# 创建必要的目录
mkdir -p logs uploads

# 设置环境变量（如果需要）
export APP_ENV=${APP_ENV:-development}

# 启动服务
echo "Starting server on port 8080..."
go run cmd/server/main.go