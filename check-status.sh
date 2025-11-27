#!/bin/bash

# LanLink 状态检查脚本
# 用于快速检查 LanLink 运行状态

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

echo -e "\n${CYAN}====== LanLink 状态检查 ======${NC}\n"

# 1. 检查进程
echo -e "${YELLOW}1. 进程状态:${NC}"
if pgrep -f lanlink > /dev/null; then
    PID=$(pgrep -f lanlink | head -1)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        MEM=$(ps -o rss= -p $PID | awk '{printf "%.2f", $1/1024}')
    else
        MEM=$(ps -o rss= -p $PID | awk '{printf "%.2f", $1/1024}')
    fi
    echo -e "   ${GREEN}✓ 进程运行中 (PID: $PID)${NC}"
    echo -e "   ${GRAY}内存占用: ${MEM} MB${NC}"
else
    echo -e "   ${RED}✗ 进程未运行${NC}"
fi

# 2. 检查日志
echo -e "\n${YELLOW}2. 日志状态:${NC}"
if [ -f "lanlink.log" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        LAST_MOD=$(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" lanlink.log)
        LAST_MOD_SEC=$(stat -f "%m" lanlink.log)
    else
        LAST_MOD=$(stat -c "%y" lanlink.log | cut -d'.' -f1)
        LAST_MOD_SEC=$(stat -c "%Y" lanlink.log)
    fi
    
    NOW_SEC=$(date +%s)
    DIFF=$((NOW_SEC - LAST_MOD_SEC))
    
    echo -e "   ${GREEN}✓ 日志文件存在${NC}"
    echo -e "   ${GRAY}最后更新: $LAST_MOD (${DIFF}秒前)${NC}"
    
    if [ $DIFF -lt 30 ]; then
        echo -e "   ${GREEN}✓ 日志活跃（系统正常运行）${NC}"
    else
        echo -e "   ${YELLOW}⚠ 日志停滞（可能异常）${NC}"
    fi
    
    echo -e "\n   ${GRAY}最新日志:${NC}"
    tail -5 lanlink.log | while read line; do
        echo -e "   ${GRAY}$line${NC}"
    done
else
    echo -e "   ${RED}✗ 日志文件不存在${NC}"
fi

# 3. 检查 Hosts 文件
echo -e "\n${YELLOW}3. Hosts 文件:${NC}"
if [[ "$OSTYPE" == "darwin"* ]] || [[ "$OSTYPE" == "linux-gnu"* ]]; then
    HOSTS_PATH="/etc/hosts"
else
    HOSTS_PATH="/etc/hosts"
fi

if [ -r "$HOSTS_PATH" ]; then
    IN_ZONE=0
    COUNT=0
    ENTRIES=()
    
    while IFS= read -r line; do
        if [[ $line == *"LanLink Managed Begin"* ]]; then
            IN_ZONE=1
        elif [[ $line == *"LanLink Managed End"* ]]; then
            IN_ZONE=0
        elif [[ $IN_ZONE -eq 1 && $line =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+ ]]; then
            ((COUNT++))
            ENTRIES+=("$line")
        fi
    done < "$HOSTS_PATH"
    
    if [ $COUNT -gt 0 ]; then
        echo -e "   ${GREEN}✓ 发现 $COUNT 个设备:${NC}"
        for entry in "${ENTRIES[@]}"; do
            echo -e "     ${GRAY}$entry${NC}"
        done
    else
        echo -e "   ${YELLOW}⚠ 未发现任何设备（可能还未发现其他节点）${NC}"
    fi
    
    # 4. 网络测试
    echo -e "\n${YELLOW}4. 网络测试:${NC}"
    if [ $COUNT -gt 0 ]; then
        # 获取第一个设备的域名
        FIRST_ENTRY="${ENTRIES[0]}"
        TEST_DOMAIN=$(echo "$FIRST_ENTRY" | awk '{print $2}')
        echo -e "   ${GRAY}测试 $TEST_DOMAIN ...${NC}"
        
        if ping -c 1 -W 1 "$TEST_DOMAIN" > /dev/null 2>&1; then
            echo -e "   ${GREEN}✓ Ping 成功（域名解析正常）${NC}"
        else
            echo -e "   ${YELLOW}⚠ Ping 失败（设备可能离线）${NC}"
        fi
    else
        echo -e "   ${YELLOW}⚠ 无可测试的设备${NC}"
    fi
else
    echo -e "   ${RED}✗ 无法读取 Hosts 文件（需要 root 权限）${NC}"
fi

# 5. 配置检查
echo -e "\n${YELLOW}5. 配置状态:${NC}"
if [ -f "config.json" ]; then
    echo -e "   ${GREEN}✓ 配置文件存在${NC}"
    if command -v jq &> /dev/null; then
        DEVICE_NAME=$(jq -r '.deviceName // "默认"' config.json)
        DOMAIN_SUFFIX=$(jq -r '.domainSuffix // "local"' config.json)
        HEARTBEAT=$(jq -r '.heartbeatIntervalSec // "10"' config.json)
        echo -e "   ${GRAY}设备名: $DEVICE_NAME${NC}"
        echo -e "   ${GRAY}域名后缀: $DOMAIN_SUFFIX${NC}"
        echo -e "   ${GRAY}心跳间隔: $HEARTBEAT 秒${NC}"
    else
        echo -e "   ${GRAY}(安装 jq 可查看详细配置)${NC}"
    fi
else
    echo -e "   ${YELLOW}⚠ 使用默认配置${NC}"
fi

# 总结
echo -e "\n${CYAN}=============================${NC}"
echo -e "${YELLOW}状态总结:${NC}"

ALL_GOOD=true
if ! pgrep -f lanlink > /dev/null; then
    echo -e "  ${RED}✗ 程序未运行${NC}"
    ALL_GOOD=false
fi
if [ ! -f "lanlink.log" ]; then
    echo -e "  ${RED}✗ 日志文件不存在${NC}"
    ALL_GOOD=false
fi

if [ "$ALL_GOOD" = true ]; then
    echo -e "  ${GREEN}✓ 系统运行正常${NC}"
else
    echo -e "  ${YELLOW}⚠ 检测到问题，请检查上述输出${NC}"
fi

echo -e "\n${GRAY}提示: 使用 'tail -f lanlink.log' 实时查看日志${NC}"
echo -e "${CYAN}=============================${NC}"
echo ""

