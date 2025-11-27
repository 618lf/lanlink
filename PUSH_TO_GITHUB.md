# 推送到 GitHub 指南

## 📋 准备工作

确保您已经在 GitHub 上创建了仓库：https://github.com/618lf/lanlink

## 🚀 推送步骤

### 1. 初始化 Git 仓库（如果还没有）

```bash
git init
```

### 2. 添加所有文件

```bash
git add .
```

### 3. 提交

```bash
git commit -m "Initial commit: LanLink v1.0.0 - 局域网域名自动映射工具"
```

### 4. 关联远程仓库

```bash
git remote add origin https://github.com/618lf/lanlink.git
```

### 5. 推送到 GitHub

```bash
git branch -M main
git push -u origin main
```

## 🏷️ 发布版本（可选）

创建标签并推送，会自动触发 GitHub Actions 编译发布：

```bash
# 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送标签
git push origin v1.0.0
```

## 📦 GitHub Actions 自动发布

已配置 `.github/workflows/release.yml`，当推送标签时会自动：

1. 编译 Linux/Mac/Windows 版本
2. 创建 GitHub Release
3. 上传编译好的二进制文件

## 📝 项目说明

在 GitHub 仓库设置中添加项目描述：

```
局域网域名自动映射工具 - 无中心节点，跨平台，基于 UDP 组播实现设备 IP 自动发现与同步
```

标签：
```
golang, lan, udp-multicast, hosts-management, network-tools, cross-platform
```

## 🎉 完成

推送成功后，访问：https://github.com/618lf/lanlink

