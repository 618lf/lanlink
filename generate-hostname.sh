#!/bin/bash
#
# generate-hostname.sh - 基于主板序列号生成唯一主机域名
#
# 使用方式:
#   ./generate-hostname.sh    # 输出域名
#
# 域名格式: {平台}-{序列号后6位}.coobee.local
#

DOMAIN_SUFFIX="coobee.local"

# 获取平台类型
get_platform() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "linux"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "msys"* ]] || [[ "$OSTYPE" == "cygwin"* ]]; then
        echo "win"
    else
        echo "unknown"
    fi
}

# 获取主板序列号
get_serial_number() {
    local platform=$(get_platform)
    local serial=""
    
    case $platform in
        macos)
            serial=$(system_profiler SPHardwareDataType 2>/dev/null | grep "Serial Number (system)" | awk -F': ' '{print $2}')
            if [[ -z "$serial" ]]; then
                serial=$(ioreg -l 2>/dev/null | grep IOPlatformSerialNumber | awk -F'"' '{print $4}')
            fi
            ;;
        linux)
            if [[ -f /sys/class/dmi/id/product_serial ]]; then
                serial=$(cat /sys/class/dmi/id/product_serial 2>/dev/null)
            fi
            if [[ -z "$serial" ]] || [[ "$serial" == "To Be Filled By O.E.M." ]]; then
                serial=$(cat /etc/machine-id 2>/dev/null)
            fi
            ;;
        win)
            serial=$(wmic bios get serialnumber 2>/dev/null | tail -2 | head -1 | tr -d ' \r\n')
            ;;
    esac
    
    echo "$serial"
}

# 生成主机标识（序列号后6位）
generate_host_id() {
    local serial=$1
    # 转小写，移除非字母数字字符
    local cleaned=$(echo "$serial" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]//g')
    
    # 取后6位
    local len=${#cleaned}
    if [[ $len -gt 6 ]]; then
        host_id="${cleaned: -6}"
    else
        host_id="$cleaned"
    fi
    
    # 不足6位用0填充
    while [[ ${#host_id} -lt 6 ]]; do
        host_id="0${host_id}"
    done
    
    echo "$host_id"
}

# 主函数
main() {
    local platform=$(get_platform)
    local serial=$(get_serial_number)
    local host_id=$(generate_host_id "$serial")
    
    echo "${platform}-${host_id}.${DOMAIN_SUFFIX}"
}

main
