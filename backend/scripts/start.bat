@echo off
echo Starting Online Mall Backend...

REM 检查Go环境
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed
    pause
    exit /b 1
)

REM 检查依赖
echo Checking dependencies...
go mod tidy

REM 创建必要的目录
if not exist "logs" mkdir logs
if not exist "uploads" mkdir uploads

REM 设置环境变量（如果需要）
if not defined APP_ENV set APP_ENV=development

REM 启动服务
echo Starting server on port 8080...
go run cmd/server/main.go

pause