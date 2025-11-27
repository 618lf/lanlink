# LanLink CLI 设计方案

## 🎯 设计目标

1. **简洁易用** - 符合直觉的命令设计
2. **功能完整** - 覆盖常用运维场景
3. **不过度设计** - 使用标准库，避免复杂依赖
4. **向后兼容** - 保持现有功能不变

## 📋 命令设计

### 核心命令

```bash
# 默认：启动服务（后台运行）
lanlink
lanlink start

# 查看运行状态
lanlink status

# 列出所有节点
lanlink list
lanlink ls

# 查看日志
lanlink logs

# 测试连接
lanlink ping <domain>

# 查看版本
lanlink version
lanlink -v

# 帮助信息
lanlink help
lanlink -h
```

### 高级命令

```bash
# 配置管理
lanlink config show                    # 显示当前配置
lanlink config set <key> <value>       # 设置配置项
lanlink config reset                   # 重置为默认配置

# Hosts管理
lanlink hosts show                     # 显示LanLink管理的条目
lanlink hosts backup                   # 手动备份hosts文件
lanlink hosts restore                  # 恢复hosts备份
lanlink hosts clean                    # 清理离线节点

# 诊断工具
lanlink diagnose                       # 运行诊断检查
lanlink check                          # 别名
```

## 💡 命令详细说明

### 1. `lanlink` / `lanlink start`

**功能**：启动 LanLink 服务（当前默认行为）

**示例**：
```bash
lanlink
# 或
lanlink start
```

**输出**：
```
LanLink - 局域网域名自动映射工具
Version: 1.0.0

[2024-11-27 14:30:00] [INFO] === LanLink 启动 ===
[2024-11-27 14:30:00] [INFO] 设备名称: mypc
...
LanLink 运行中，按 Ctrl+C 退出...
```

---

### 2. `lanlink status`

**功能**：查看运行状态和统计信息

**示例**：
```bash
lanlink status
```

**输出**：
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  LanLink 状态概览
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

运行状态:   ✓ 运行中 (PID: 12345)
启动时间:   2024-11-27 14:00:00 (运行 30 分钟)
内存占用:   18.5 MB
CPU 使用:   0.2%

本机信息:
  设备名:   mypc
  域名:     mypc.local
  IP:       192.168.1.100
  MAC:      00:11:22:33:44:55

网络配置:
  组播地址: 239.255.0.1:9527
  心跳间隔: 10 秒
  离线超时: 30 秒

节点统计:
  在线节点: 3 个
  离线节点: 1 个
  总节点:   4 个

最近活动:
  [14:30:00] server.local 上线 (192.168.1.101)
  [14:30:05] nas.local 上线 (192.168.1.102)
  [14:31:00] pi.local 离线

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**参数**：
```bash
lanlink status --json    # JSON格式输出
lanlink status --simple  # 简化输出
```

---

### 3. `lanlink list` / `lanlink ls`

**功能**：列出所有发现的节点

**示例**：
```bash
lanlink list
# 或
lanlink ls
```

**输出**：
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  节点列表 (在线: 3, 离线: 1)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

状态  域名              IP 地址          主机名      最后心跳         
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓    mypc.local       192.168.1.100    MYPC        (本机)
✓    server.local     192.168.1.101    SERVER      5秒前
✓    nas.local        192.168.1.102    NAS         8秒前
✗    pi.local         192.168.1.103    PI          2分钟前 (离线)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**参数**：
```bash
lanlink list --online     # 仅显示在线节点
lanlink list --offline    # 仅显示离线节点
lanlink list --json       # JSON格式输出
lanlink list --watch      # 实时监控模式
```

**JSON 输出**：
```json
{
  "nodes": [
    {
      "domain": "server.local",
      "ip": "192.168.1.101",
      "hostname": "SERVER",
      "deviceId": "mac-aabbccddeeff",
      "status": "online",
      "lastSeen": "2024-11-27T14:30:05Z"
    }
  ],
  "summary": {
    "online": 3,
    "offline": 1,
    "total": 4
  }
}
```

---

### 4. `lanlink logs`

**功能**：查看日志

**示例**：
```bash
lanlink logs              # 显示最后50行
lanlink logs -n 100       # 显示最后100行
lanlink logs -f           # 实时跟踪（类似 tail -f）
lanlink logs --level error # 仅显示错误
```

**输出**：
```
[2024-11-27 14:30:00] [INFO] === LanLink 启动 ===
[2024-11-27 14:30:00] [INFO] 设备名称: mypc
[2024-11-27 14:30:05] [INFO] 节点上线: server (server.local -> 192.168.1.101)
[2024-11-27 14:30:05] [INFO] 已更新hosts: server.local -> 192.168.1.101
...
```

---

### 5. `lanlink ping <domain>`

**功能**：测试与指定节点的连接

**示例**：
```bash
lanlink ping server.local
```

**输出**：
```
正在测试连接: server.local
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ DNS 解析:     server.local -> 192.168.1.101
✓ Hosts 记录:   已找到
✓ Ping 测试:    成功 (延迟: 0.5ms)
✓ 节点状态:     在线 (最后心跳: 5秒前)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
连接正常 ✓
```

**失败输出**：
```
正在测试连接: unknown.local
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✗ DNS 解析:     失败 (未找到域名)
✗ Hosts 记录:   未找到
- Ping 测试:    跳过
✗ 节点状态:     未知

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
连接失败 ✗

建议: 确认目标设备是否运行 LanLink
```

---

### 6. `lanlink config`

**功能**：配置管理

**示例**：
```bash
# 显示当前配置
lanlink config show

# 设置配置项
lanlink config set deviceName mypc
lanlink config set heartbeatIntervalSec 15

# 重置配置
lanlink config reset
```

**输出**：
```
当前配置:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
deviceName:           mypc
domainSuffix:         local
multicastAddr:        239.255.0.1
multicastPort:        9527
heartbeatIntervalSec: 10
offlineTimeoutSec:    30
logLevel:             info
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

配置文件: D:\git10\LanLink\config.json
```

---

### 7. `lanlink hosts`

**功能**：Hosts 文件管理

**示例**：
```bash
# 显示 LanLink 管理的条目
lanlink hosts show

# 手动备份
lanlink hosts backup

# 恢复备份
lanlink hosts restore

# 清理离线节点
lanlink hosts clean
```

**输出**：
```
LanLink 管理的 Hosts 条目:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
192.168.1.101   server.local    # LanLink
192.168.1.102   nas.local       # LanLink
192.168.1.103   pi.local        # LanLink
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
共 3 个条目

Hosts 文件: C:\Windows\System32\drivers\etc\hosts
备份文件:   C:\Windows\System32\drivers\etc\hosts.bak
```

---

### 8. `lanlink diagnose` / `lanlink check`

**功能**：运行诊断检查

**示例**：
```bash
lanlink diagnose
# 或
lanlink check
```

**输出**：
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  LanLink 系统诊断
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

运行环境:
  ✓ 操作系统:   Windows 11
  ✓ Go 版本:    1.23
  ✓ 管理员权限: 是

配置检查:
  ✓ 配置文件:   存在且有效
  ✓ 日志文件:   正常 (最后更新: 5秒前)

网络检查:
  ✓ 本机 IP:    192.168.1.100
  ✓ MAC 地址:   00:11:22:33:44:55
  ✓ 组播支持:   是
  ⚠ 防火墙:     未检测到规则 (可能影响通信)

Hosts 检查:
  ✓ 文件权限:   可写
  ✓ 标记区域:   存在
  ✓ 备份文件:   存在

进程检查:
  ✓ 运行状态:   正常
  ✓ 内存占用:   18.5 MB (正常)
  ✓ CPU 使用:   0.2% (正常)

节点通信:
  ✓ 在线节点:   3 个
  ✓ 心跳正常:   是
  ✗ 离线节点:   1 个 (pi.local 已离线 2分钟)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
诊断结果: 1 个警告

建议:
  ⚠ 添加防火墙规则允许 UDP 9527 端口
```

---

### 9. `lanlink version`

**功能**：显示版本信息

**示例**：
```bash
lanlink version
# 或
lanlink -v
```

**输出**：
```
LanLink v1.0.0
Build: 2024-11-27
Commit: abc1234
Go Version: go1.23
Platform: windows/amd64
```

---

## 🏗️ 实现架构

### 文件结构

```
LanLink/
├── main.go              # 主入口，路由命令
├── cli/
│   ├── cli.go          # CLI 框架
│   ├── start.go        # start 命令
│   ├── status.go       # status 命令
│   ├── list.go         # list 命令
│   ├── logs.go         # logs 命令
│   ├── ping.go         # ping 命令
│   ├── config.go       # config 命令
│   ├── hosts.go        # hosts 命令
│   ├── diagnose.go     # diagnose 命令
│   └── version.go      # version 命令
├── internal/
│   └── info.go         # 读取运行时信息（日志、进程等）
└── ...
```

### 技术实现

#### 1. 命令行参数解析

使用标准库 `flag`，简单直接：

```go
// main.go
func main() {
    if len(os.Args) < 2 {
        // 默认启动服务
        runService()
        return
    }

    command := os.Args[1]
    switch command {
    case "start":
        runService()
    case "status":
        cli.ShowStatus()
    case "list", "ls":
        cli.ListNodes()
    case "logs":
        cli.ShowLogs(os.Args[2:])
    case "ping":
        cli.PingNode(os.Args[2:])
    case "config":
        cli.ConfigCommand(os.Args[2:])
    case "hosts":
        cli.HostsCommand(os.Args[2:])
    case "diagnose", "check":
        cli.Diagnose()
    case "version", "-v", "--version":
        cli.ShowVersion()
    case "help", "-h", "--help":
        cli.ShowHelp()
    default:
        fmt.Printf("未知命令: %s\n", command)
        cli.ShowHelp()
        os.Exit(1)
    }
}
```

#### 2. 状态信息获取

不需要复杂的 IPC，直接读取文件：

```go
// internal/info.go
type RuntimeInfo struct {
    IsRunning    bool
    PID          int
    Memory       int64
    StartTime    time.Time
    Nodes        []NodeInfo
    LogEntries   []LogEntry
}

func GetRuntimeInfo() (*RuntimeInfo, error) {
    info := &RuntimeInfo{}
    
    // 1. 检查进程
    info.IsRunning, info.PID = checkProcess()
    
    // 2. 读取日志
    info.LogEntries = parseLogFile("lanlink.log")
    
    // 3. 读取 Hosts
    info.Nodes = parseHostsFile()
    
    return info, nil
}
```

#### 3. 美化输出

使用简单的 ANSI 颜色码：

```go
// cli/ui.go
const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
)

func Success(msg string) {
    fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, msg)
}

func Error(msg string) {
    fmt.Printf("%s✗%s %s\n", ColorRed, ColorReset, msg)
}

func Warn(msg string) {
    fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, msg)
}
```

---

## 📊 命令优先级

### MVP (最小可行版本)

**必须实现**：
1. ✅ `lanlink` / `start` - 启动服务（已有）
2. ✅ `lanlink status` - 查看状态
3. ✅ `lanlink list` - 列出节点
4. ✅ `lanlink logs` - 查看日志
5. ✅ `lanlink version` - 版本信息
6. ✅ `lanlink help` - 帮助信息

### V1.1 (增强版)

**可选实现**：
7. ⭐ `lanlink ping` - 测试连接
8. ⭐ `lanlink diagnose` - 诊断检查

### V1.2 (完整版)

**后续实现**：
9. 🔮 `lanlink config` - 配置管理
10. 🔮 `lanlink hosts` - Hosts 管理

---

## 🎨 用户体验设计

### 1. 帮助信息友好

```bash
lanlink help
```

输出：
```
LanLink - 局域网域名自动映射工具

用法:
  lanlink [command] [options]

命令:
  start              启动服务（默认）
  status             查看运行状态
  list, ls           列出所有节点
  logs               查看日志
  ping <domain>      测试连接
  config             配置管理
  hosts              Hosts文件管理
  diagnose, check    运行诊断
  version, -v        显示版本
  help, -h           显示帮助

示例:
  lanlink                      # 启动服务
  lanlink status               # 查看状态
  lanlink list --online        # 仅显示在线节点
  lanlink logs -f              # 实时查看日志
  lanlink ping server.local    # 测试连接

更多信息: https://github.com/618lf/lanlink
```

### 2. 错误提示清晰

```bash
lanlink ping
# 错误: 缺少参数 <domain>
# 用法: lanlink ping <domain>
# 示例: lanlink ping server.local
```

### 3. 进度提示

```bash
lanlink diagnose
# 正在检查运行环境... ✓
# 正在检查网络配置... ✓
# 正在检查节点通信... ✓
# 诊断完成!
```

---

## 🔧 实现细节

### 跨平台兼容

```go
// 检测管理员权限
func isAdmin() bool {
    if runtime.GOOS == "windows" {
        // Windows: 尝试打开需要管理员权限的文件
        _, err := os.OpenFile(`\\.\PHYSICALDRIVE0`, os.O_RDONLY, 0)
        return err == nil
    } else {
        // Linux/Mac: 检查 UID
        return os.Geteuid() == 0
    }
}
```

### 颜色输出兼容

```go
// Windows 需要启用 ANSI 颜色支持
func initColorOutput() {
    if runtime.GOOS == "windows" {
        // 启用 Windows 10+ 的 ANSI 颜色支持
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        setConsoleMode := kernel32.NewProc("SetConsoleMode")
        setConsoleMode.Call(uintptr(syscall.Stdout), 0x0001|0x0002|0x0004)
    }
}
```

---

## 📈 性能考虑

1. **快速响应** - 命令执行 < 100ms
2. **轻量级** - 不引入重量级依赖
3. **缓存信息** - 避免重复读取文件

---

## 🚀 实现计划

### Phase 1: 基础命令 (1-2天)
- [ ] main.go 命令路由
- [ ] status 命令
- [ ] list 命令
- [ ] logs 命令
- [ ] version/help 命令

### Phase 2: 增强功能 (1天)
- [ ] ping 命令
- [ ] diagnose 命令
- [ ] 美化输出

### Phase 3: 高级功能 (可选)
- [ ] config 命令
- [ ] hosts 命令
- [ ] JSON 输出支持

---

**这个设计遵循了"不过度设计"的原则，使用标准库实现，功能完整且易于维护！** ✨

