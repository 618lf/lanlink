# LanLink - 局域网域名自动映射工具

LanLink 是一款无中心节点、跨平台的局域网工具，基于 UDP 组播实现设备 IP 自动发现与同步，自动维护各设备的 Hosts 文件，让用户以**固定域名**访问局域网内设备。

## ✨ 核心特性

- 🚀 **零配置启动** - 开箱即用，自动获取设备信息
- 🌐 **跨平台支持** - Windows/Mac/Linux 单文件运行
- 🔄 **自动同步** - 实时发现局域网设备，自动更新 Hosts
- 🎯 **固定域名** - 告别动态IP，使用 `devicename.local` 访问设备
- 🛡️ **安全可靠** - 仅管理标记区域，不影响用户自定义配置
- 📊 **可观测** - 详细日志记录，运行状态一目了然

## 📦 命令行工具

LanLink 提供强大的命令行工具，方便监控和管理。

```bash
lanlink              # 启动服务（默认）
lanlink status       # 查看运行状态
lanlink list         # 列出所有节点
lanlink logs -f      # 实时查看日志
lanlink ping <domain> # 测试连接
lanlink version      # 显示版本
lanlink help         # 帮助信息
```

详细使用说明请查看：[CLI使用指南](docs/CLI使用指南.md)

---

## 🚀 快速开始

### 1. 下载或编译

```bash
# 克隆仓库
git clone https://github.com/618lf/lanlink.git
cd lanlink

# 编译
go build -o lanlink

# Windows
go build -o lanlink.exe

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o lanlink-linux
GOOS=darwin GOARCH=arm64 go build -o lanlink-mac
GOOS=windows GOARCH=amd64 go build -o lanlink.exe
```

### 2. 运行

**Windows（需要管理员权限）：**
```powershell
# 右键以管理员身份运行，或在管理员PowerShell中：
.\lanlink.exe
```

**Mac/Linux（需要root权限）：**
```bash
sudo ./lanlink
```

### 3. 使用

程序启动后会自动：
- 获取本机 IP 和 MAC 地址
- 生成域名（如 `mypc.local`）
- 加入组播组，监听其他节点
- 定期发送心跳包
- 自动更新 Hosts 文件

在局域网内的其他设备上启动 LanLink 后，就可以通过域名访问：

```bash
# Ping 测试
ping mypc.local

# SSH 连接
ssh user@server.local

# 访问Web服务
curl http://webserver.local:8080
```

## ⚙️ 配置

首次运行会使用默认配置，如需自定义，创建 `config.json`：

```json
{
  "deviceName": "mypc",           // 设备名（空则使用主机名）
  "domainSuffix": "local",        // 域名后缀
  "multicastAddr": "239.255.0.1", // 组播地址
  "multicastPort": 9527,          // 组播端口
  "heartbeatIntervalSec": 10,     // 心跳间隔（秒）
  "offlineTimeoutSec": 30,        // 离线超时（秒）
  "logLevel": "info"              // 日志级别：debug/info/warn/error
}
```

## 📁 项目结构

```
LanLink/
├── main.go              # 程序入口
├── config/             
│   └── config.go        # 配置管理
├── network/            
│   └── multicast.go     # 组播通信
├── hosts/              
│   └── manager.go       # Hosts文件管理
├── node/               
│   └── manager.go       # 节点管理
├── logger/             
│   └── logger.go        # 日志模块
├── config.example.json  # 配置示例
└── README.md
```

## 🔍 工作原理

### 1. 组播通信

- 使用 UDP 组播地址 `239.255.0.1:9527`
- 每 10 秒发送心跳包（携带域名、IP、设备ID）
- 超过 30 秒未收到心跳判定为离线

### 2. Hosts 文件管理

LanLink 在 Hosts 文件中创建专属管理区域：

```text
# === LanLink Managed Begin ===
192.168.1.100  mypc.local      # LanLink
192.168.1.101  server.local    # LanLink
# === LanLink Managed End ===
```

- ✅ 仅在标记区域内操作
- ✅ 修改前自动备份（`hosts.bak`）
- ✅ 不影响用户手动配置

### 3. 节点管理

- 使用 MAC 地址作为唯一设备标识
- 自动检测节点上线/离线
- 处理域名冲突（自动添加后缀）

## 📊 可观测性

### 日志输出

日志同时输出到**控制台**和 `lanlink.log` 文件：

```
[2024-11-27 10:30:00] [INFO] === LanLink 启动 ===
[2024-11-27 10:30:00] [INFO] 设备名称: mypc
[2024-11-27 10:30:00] [INFO] 本机信息: DeviceID=mac-00:11:22:33:44:55, IP=192.168.1.100
[2024-11-27 10:30:05] [INFO] 节点上线: server (server.local -> 192.168.1.101)
[2024-11-27 10:30:05] [INFO] 已更新hosts: server.local -> 192.168.1.101
```

### 日志级别

- `debug`: 详细调试信息（包括每次心跳）
- `info`: 重要事件（节点上下线、hosts更新）
- `warn`: 警告信息（域名冲突）
- `error`: 错误信息

## ⚠️ 注意事项

1. **权限要求**：必须以管理员/root 权限运行才能修改 Hosts 文件
2. **防火墙**：需要放行 UDP 9527 端口（或配置的端口）
3. **组播支持**：确保路由器/交换机支持组播（大部分家庭/办公网络支持）
4. **域名唯一性**：同一局域网内避免设备名重复
5. **域名后缀**：建议使用 `.local`，避免与公网域名冲突

## 🛠️ 故障排查

### 问题：无权限修改 Hosts 文件

**解决**：以管理员/root 身份运行程序

### 问题：无法发现其他节点

**原因**：
- 防火墙阻止 UDP 组播
- 路由器禁用组播
- 设备不在同一局域网

**解决**：
1. 检查防火墙设置，放行 UDP 9527 端口
2. 查看日志确认是否正常发送/接收心跳
3. 确保所有设备在同一子网

### 问题：域名冲突

**现象**：日志中出现域名冲突警告，域名被自动重命名

**解决**：修改 `config.json` 中的 `deviceName`，确保唯一性

## 🚀 后续计划

- [ ] 多网卡支持
- [ ] 网络自动降级（组播失败时切换到广播）
- [ ] 命令行工具（查看节点列表、测试连接）
- [ ] 系统服务安装脚本
- [ ] Web 管理界面
- [ ] 消息加密与设备认证

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**LanLink** - 让局域网设备通过固定域名无缝链接 🔗

