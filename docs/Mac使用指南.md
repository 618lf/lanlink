# LanLink Mac 使用指南

## 🍎 Mac 系统安装和使用

### 快速开始

#### 1. 下载或编译

**从 Windows 交叉编译 Mac 版本：**

```powershell
# Mac Intel (amd64)
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o lanlink-mac-amd64

# Mac Apple Silicon (arm64 - M1/M2/M3)
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o lanlink-mac-arm64
```

**或在 Mac 上直接编译：**

```bash
# 会自动检测 CPU 架构
go build -o lanlink
```

#### 2. 赋予执行权限

```bash
chmod +x lanlink-mac-arm64
# 或
chmod +x lanlink-mac-amd64
```

#### 3. 运行

**快速测试（前台运行）：**

```bash
sudo ./lanlink-mac-arm64
```

---

### 完整安装（推荐）

#### 方式一：安装到系统 PATH

```bash
# 1. 安装到 /usr/local/bin
sudo ./lanlink-mac-arm64 install

# 2. 验证
lanlink version
```

#### 方式二：安装为系统服务（开机自启）

```bash
# 1. 先安装到系统
sudo ./lanlink-mac-arm64 install

# 2. 安装为系统服务
sudo lanlink service install

# 3. 启动服务
sudo lanlink service start

# 4. 查看状态
lanlink service status
```

---

### Mac 特定说明

#### 安全性和隐私

Mac 首次运行可能会提示"无法验证开发者"：

**解决方法：**

1. 系统偏好设置 → 安全性与隐私
2. 在"通用"选项卡中点击"仍要打开"

**或使用命令行：**

```bash
# 移除隔离属性
sudo xattr -d com.apple.quarantine lanlink-mac-arm64
```

#### 系统服务管理

Mac 使用 `launchd` 管理服务（类似 systemd）：

```bash
# 查看服务状态
sudo launchctl list | grep lanlink

# 手动加载服务
sudo launchctl load /Library/LaunchDaemons/com.lanlink.plist

# 手动卸载服务
sudo launchctl unload /Library/LaunchDaemons/com.lanlink.plist
```

#### 防火墙设置

如果启用了防火墙，需要允许 LanLink：

```bash
# 允许 lanlink 通过防火墙
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/bin/lanlink
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /usr/local/bin/lanlink
```

---

### 常用命令

```bash
# 查看状态
lanlink status

# 列出节点
lanlink list

# 实时日志
lanlink logs -f

# Ping 测试
lanlink ping server.local

# 服务管理
sudo lanlink service start
sudo lanlink service stop
sudo lanlink service status
```

---

### Hosts 文件位置

Mac 的 Hosts 文件位置：`/etc/hosts`

查看 LanLink 管理的条目：

```bash
sudo cat /etc/hosts | grep -A 20 "LanLink Managed Begin"
```

---

### 卸载

```bash
# 1. 卸载服务
sudo lanlink service uninstall

# 2. 完全卸载
sudo lanlink uninstall
```

---

## 🔨 在 Windows 上交叉编译 Mac 版本

Go 语言支持交叉编译，可以在 Windows 上编译 Mac 版本。

### 方法一：使用 PowerShell

```powershell
# 编译 Mac Intel 版本 (amd64)
$env:GOOS="darwin"
$env:GOARCH="amd64"
go build -o dist/lanlink-mac-amd64

# 编译 Mac Apple Silicon 版本 (arm64 - M1/M2/M3)
$env:GOOS="darwin"
$env:GOARCH="arm64"
go build -o dist/lanlink-mac-arm64

# 重置环境变量
$env:GOOS=""
$env:GOARCH=""
```

### 方法二：使用 CMD

```batch
REM Mac Intel (amd64)
set GOOS=darwin
set GOARCH=amd64
go build -o dist\lanlink-mac-amd64

REM Mac Apple Silicon (arm64)
set GOOS=darwin
set GOARCH=arm64
go build -o dist\lanlink-mac-arm64

REM 重置
set GOOS=
set GOARCH=
```

### 方法三：使用 build.bat（已创建）

直接运行项目中的 `build.bat`，会自动编译所有平台版本：

```batch
build.bat
```

生成的文件在 `dist/` 目录：
- `lanlink-mac-amd64` - Mac Intel 版本
- `lanlink-mac-arm64` - Mac Apple Silicon 版本
- `lanlink-windows-amd64.exe` - Windows 版本
- `lanlink-linux-amd64` - Linux 版本

### 一键编译所有平台（推荐）

创建 `build-all.ps1`:

```powershell
# 创建输出目录
New-Item -ItemType Directory -Force -Path dist | Out-Null

Write-Host "开始编译所有平台版本..." -ForegroundColor Cyan

# Windows amd64
Write-Host "`n编译 Windows (amd64)..." -ForegroundColor Yellow
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -o dist/lanlink-windows-amd64.exe
if ($LASTEXITCODE -eq 0) { Write-Host "✓ Windows 编译成功" -ForegroundColor Green }

# Linux amd64
Write-Host "`n编译 Linux (amd64)..." -ForegroundColor Yellow
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -o dist/lanlink-linux-amd64
if ($LASTEXITCODE -eq 0) { Write-Host "✓ Linux 编译成功" -ForegroundColor Green }

# Mac Intel (amd64)
Write-Host "`n编译 Mac Intel (amd64)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="amd64"
go build -o dist/lanlink-mac-amd64
if ($LASTEXITCODE -eq 0) { Write-Host "✓ Mac Intel 编译成功" -ForegroundColor Green }

# Mac Apple Silicon (arm64)
Write-Host "`n编译 Mac Apple Silicon (arm64)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="arm64"
go build -o dist/lanlink-mac-arm64
if ($LASTEXITCODE -eq 0) { Write-Host "✓ Mac Apple Silicon 编译成功" -ForegroundColor Green }

# 重置环境变量
$env:GOOS=""; $env:GOARCH=""

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "编译完成！" -ForegroundColor Green
Write-Host "`n输出目录: dist/" -ForegroundColor Yellow
Get-ChildItem dist/ | Format-Table Name, Length, LastWriteTime
```

运行：

```powershell
.\build-all.ps1
```

---

## 🔍 如何知道 Mac 使用哪个版本？

### 检查 Mac CPU 架构

在 Mac 上运行：

```bash
# 查看架构
uname -m

# 输出:
# x86_64  → 使用 lanlink-mac-amd64 (Intel)
# arm64   → 使用 lanlink-mac-arm64 (Apple Silicon)
```

### 或者查看系统信息

```bash
# 查看详细信息
sysctl -a | grep machdep.cpu.brand_string

# Intel Mac 会显示: Intel
# Apple Silicon 会显示: Apple M1/M2/M3
```

### 关于这一点 - 点击左上角苹果图标

1. 点击"关于本机"
2. 查看"芯片"或"处理器"
   - 如果是"Apple M1/M2/M3" → 使用 arm64 版本
   - 如果是"Intel Core" → 使用 amd64 版本

---

## 📊 平台对照表

| 平台 | 架构 | 文件名 | 适用设备 |
|------|------|--------|---------|
| Windows | amd64 | `lanlink-windows-amd64.exe` | 所有 Windows PC |
| Linux | amd64 | `lanlink-linux-amd64` | 大部分 Linux PC |
| Mac Intel | amd64 | `lanlink-mac-amd64` | 2020年前的 Mac |
| Mac Apple Silicon | arm64 | `lanlink-mac-arm64` | M1/M2/M3 Mac |

---

## 💡 提示

1. **交叉编译完全支持** - Go 语言的交叉编译非常成熟，不用担心兼容性
2. **无需 Mac 设备** - 可以在 Windows 上编译 Mac 版本，然后传输到 Mac 使用
3. **文件大小** - 编译后的单个文件约 3-5MB，非常轻量
4. **CGO 注意** - LanLink 不使用 CGO，完全支持交叉编译

---

## 🚀 快速部署到 Mac

### 从 Windows 传输到 Mac

```powershell
# Windows 上编译
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o lanlink-mac-arm64

# 使用 SCP 传输到 Mac
scp lanlink-mac-arm64 user@mac-ip:~/

# 或使用共享文件夹、U盘等方式
```

### 在 Mac 上安装

```bash
# 1. 赋予执行权限
chmod +x ~/lanlink-mac-arm64

# 2. 安装
sudo ~/lanlink-mac-arm64 install

# 3. 安装为服务（可选）
sudo lanlink service install

# 4. 启动
sudo lanlink service start

# 5. 验证
lanlink status
```

---

**Mac 和 Windows/Linux 使用完全一致，命令通用！** 🍎✨

