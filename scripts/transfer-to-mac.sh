#!/bin/bash

# LanLink 传输到 Mac 的辅助脚本
# 在 Windows 编译后，使用此脚本传输到 Mac

echo "LanLink Mac 传输工具"
echo "===================="
echo ""

# 检查是否提供了 Mac 地址
if [ -z "$1" ]; then
    echo "用法: ./transfer-to-mac.sh <user@mac-ip>"
    echo "示例: ./transfer-to-mac.sh admin@192.168.1.100"
    exit 1
fi

MAC_HOST=$1

# 检测 Mac 架构
echo "正在检测 Mac CPU 架构..."
ARCH=$(ssh $MAC_HOST "uname -m")

if [ "$ARCH" = "arm64" ]; then
    FILE="dist/lanlink-mac-arm64"
    echo "检测到: Apple Silicon (M1/M2/M3)"
elif [ "$ARCH" = "x86_64" ]; then
    FILE="dist/lanlink-mac-amd64"
    echo "检测到: Intel"
else
    echo "未知架构: $ARCH"
    exit 1
fi

# 检查文件是否存在
if [ ! -f "$FILE" ]; then
    echo "错误: 文件不存在: $FILE"
    echo "请先运行 build.cmd 编译"
    exit 1
fi

# 传输文件
echo ""
echo "正在传输 $FILE 到 $MAC_HOST ..."
scp $FILE $MAC_HOST:~/lanlink

if [ $? -eq 0 ]; then
    echo "✓ 传输成功"
    echo ""
    echo "在 Mac 上运行:"
    echo "  chmod +x ~/lanlink"
    echo "  sudo ~/lanlink install"
else
    echo "✗ 传输失败"
    exit 1
fi

