#!/bin/bash

# 设置输出目录
OUTPUT_DIR="dist"
mkdir -p $OUTPUT_DIR

# 构建Mac Intel版本(默认SSE模式)
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-mac-intel ./cmd/bazi-mcp

# 构建Mac Apple Silicon版本(默认SSE模式)
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-mac-apple ./cmd/bazi-mcp

# 构建Windows 64位版本(默认SSE模式)
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-windows-amd64.exe ./cmd/bazi-mcp

# 构建Windows ARM64版本
GOOS=windows GOARCH=arm64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-windows-arm64.exe ./cmd/bazi-mcp

# 构建Windows 32位版本
GOOS=windows GOARCH=386 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-windows-386.exe ./cmd/bazi-mcp

# 构建Linux Debian/CentOS版本(默认SSE模式)
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-linux-amd64 ./cmd/bazi-mcp

# 构建Linux ARM64版本
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-linux-arm64 ./cmd/bazi-mcp

# 构建Linux ARMv7版本
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-linux-armv7 ./cmd/bazi-mcp

# 构建Linux 32位版本
GOOS=linux GOARCH=386 go build -ldflags "-X main.transport=stdio" -o $OUTPUT_DIR/bazi-mcp-linux-386 ./cmd/bazi-mcp

echo "构建完成，输出文件在 $OUTPUT_DIR 目录"