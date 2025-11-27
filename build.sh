#!/bin/bash
# LanLink 编译脚本

echo "开始编译 LanLink..."

# 创建输出目录
mkdir -p dist

# 编译 Linux 版本
echo "编译 Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/lanlink-linux-amd64

# 编译 macOS 版本
echo "编译 macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o dist/lanlink-mac-arm64

echo "编译 macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o dist/lanlink-mac-amd64

# 编译 Windows 版本
echo "编译 Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/lanlink-windows-amd64.exe

echo "编译完成！输出目录: dist/"
ls -lh dist/

