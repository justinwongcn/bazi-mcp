@echo off
chcp 65001

:: 设置输出目录
set OUTPUT_DIR=dist
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

:: Windows平台构建
:: 64位版本(amd64)
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-windows-amd64.exe" ./cmd/bazi-mcp

:: 64位版本(arm64)
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-windows-arm64.exe" ./cmd/bazi-mcp

:: 32位版本
set GOOS=windows
set GOARCH=386
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-windows-386.exe" ./cmd/bazi-mcp

:: Mac平台构建
:: Intel版本
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-mac-intel" ./cmd/bazi-mcp

:: Apple Silicon版本
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-mac-apple" ./cmd/bazi-mcp

:: Linux平台构建
:: AMD64架构
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-linux-amd64" ./cmd/bazi-mcp

:: ARM64架构
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-linux-arm64" ./cmd/bazi-mcp

:: ARMv7架构
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-linux-armv7" ./cmd/bazi-mcp

:: 386架构(32位)
set GOOS=linux
set GOARCH=386
go build -ldflags "-X main.transport=stdio" -o "%OUTPUT_DIR%\bazi-mcp-linux-386" ./cmd/bazi-mcp

echo 构建完成，输出文件在 %OUTPUT_DIR% 目录
pause