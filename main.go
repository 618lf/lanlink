package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/618lf/lanlink/config"
	"github.com/618lf/lanlink/hosts"
	"github.com/618lf/lanlink/logger"
	"github.com/618lf/lanlink/network"
	"github.com/618lf/lanlink/node"
)

const (
	configFile = "config.json"
	logFile    = "lanlink.log"
)

func main() {
	fmt.Println("LanLink - 局域网域名自动映射工具")
	fmt.Println("Version: 1.0.0")
	fmt.Println()

	// 1. 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	if err := logger.Init(cfg.LogLevel, logFile); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("=== LanLink 启动 ===")
	logger.Info("设备名称: %s", cfg.DeviceName)
	logger.Info("域名后缀: %s", cfg.DomainSuffix)

	// 3. 检查hosts文件权限
	hostsManager := hosts.NewManager()
	if err := hostsManager.CheckPermission(); err != nil {
		logger.Error("%v", err)
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 初始化hosts文件（添加标记区域）
	if err := hostsManager.Initialize(); err != nil {
		logger.Error("初始化hosts文件失败: %v", err)
		os.Exit(1)
	}
	logger.Info("Hosts文件初始化完成")

	// 4. 获取本机信息
	deviceID, err := network.GetMACAddress()
	if err != nil {
		logger.Error("获取MAC地址失败: %v", err)
		os.Exit(1)
	}

	// 5. 创建组播客户端
	client, err := network.NewMulticastClient(cfg.MulticastAddr, cfg.MulticastPort)
	if err != nil {
		logger.Error("创建组播客户端失败: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	localIP := client.GetLocalIP()
	domain := generateDomain(cfg.DeviceName, cfg.DomainSuffix)

	logger.Info("本机信息: DeviceID=%s, IP=%s, Domain=%s", deviceID, localIP, domain)

	// 6. 创建节点管理器
	nodeManager := node.NewManager(time.Duration(cfg.OfflineTimeout) * time.Second)

	// 添加本机节点
	nodeManager.AddOrUpdate(deviceID, domain, localIP, cfg.DeviceName)
	nodeManager.SetLocal(deviceID)

	// 设置节点变化回调（更新hosts文件）
	nodeManager.SetChangeCallback(func(n *node.Node, isOnline bool) {
		if n.IsLocal {
			return
		}

		if isOnline {
			logger.Info("节点上线: %s (%s -> %s)", n.Hostname, n.Domain, n.IP)
			if err := hostsManager.AddOrUpdate(n.IP, n.Domain); err != nil {
				logger.Error("更新hosts失败: %v", err)
			} else {
				logger.Info("已更新hosts: %s -> %s", n.Domain, n.IP)
			}
		} else {
			logger.Info("节点离线: %s (%s)", n.Hostname, n.Domain)
			if err := hostsManager.Remove(n.Domain); err != nil {
				logger.Error("删除hosts失败: %v", err)
			} else {
				logger.Info("已删除hosts: %s", n.Domain)
			}
		}
	})

	// 设置消息接收回调
	client.SetMessageCallback(func(msg *network.Message) {
		logger.Debug("收到消息: Action=%s, From=%s (%s)", msg.Action, msg.Hostname, msg.IP)

		switch msg.Action {
		case network.ActionHeartbeat:
			// 检查域名冲突
			_, exists := nodeManager.Get(msg.DeviceID)
			if !exists && hasDomainConflict(nodeManager, msg.Domain, msg.DeviceID) {
				// 域名冲突，添加后缀
				originalDomain := msg.Domain
				msg.Domain = msg.Domain + "-" + extractMACShort(msg.DeviceID)
				logger.Warn("域名冲突: %s 已被占用，自动重命名为 %s", originalDomain, msg.Domain)
			}

			// 更新节点
			if changed := nodeManager.AddOrUpdate(msg.DeviceID, msg.Domain, msg.IP, msg.Hostname); changed && exists {
				logger.Info("节点信息更新: %s (%s -> %s)", msg.Hostname, msg.Domain, msg.IP)
			}

		case network.ActionOffline:
			nodeManager.Remove(msg.DeviceID)
		}
	})

	// 7. 启动组播监听
	if err := client.Start(); err != nil {
		logger.Error("启动组播监听失败: %v", err)
		os.Exit(1)
	}
	logger.Info("组播监听已启动: %s:%d", cfg.MulticastAddr, cfg.MulticastPort)

	// 8. 发送首次心跳
	sendHeartbeat(client, domain, localIP, deviceID, cfg.DeviceName)

	// 9. 启动定时任务
	heartbeatTicker := time.NewTicker(time.Duration(cfg.HeartbeatInterval) * time.Second)
	offlineCheckTicker := time.NewTicker(5 * time.Second)
	defer heartbeatTicker.Stop()
	defer offlineCheckTicker.Stop()

	// 10. 优雅退出处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("LanLink 运行中，按 Ctrl+C 退出")
	fmt.Println("LanLink 运行中，按 Ctrl+C 退出...")

	// 主循环
	for {
		select {
		case <-heartbeatTicker.C:
			// 发送心跳
			sendHeartbeat(client, domain, localIP, deviceID, cfg.DeviceName)

		case <-offlineCheckTicker.C:
			// 检查离线节点
			offlineNodes := nodeManager.CheckOffline()
			if len(offlineNodes) > 0 {
				logger.Debug("检查到 %d 个离线节点", len(offlineNodes))
			}

		case <-sigChan:
			// 优雅退出
			logger.Info("收到退出信号，正在清理...")
			fmt.Println("\n正在退出...")

			// 发送离线通知
			msg := &network.Message{
				Action:   network.ActionOffline,
				Domain:   domain,
				IP:       localIP,
				DeviceID: deviceID,
				Hostname: cfg.DeviceName,
			}
			client.Send(msg)
			logger.Info("已发送离线通知")

			// 等待消息发送完成
			time.Sleep(100 * time.Millisecond)

			logger.Info("=== LanLink 已退出 ===")
			return
		}
	}
}

// sendHeartbeat 发送心跳
func sendHeartbeat(client *network.MulticastClient, domain, ip, deviceID, hostname string) {
	msg := &network.Message{
		Action:   network.ActionHeartbeat,
		Domain:   domain,
		IP:       ip,
		DeviceID: deviceID,
		Hostname: hostname,
	}

	if err := client.Send(msg); err != nil {
		logger.Error("发送心跳失败: %v", err)
	} else {
		logger.Debug("已发送心跳: %s -> %s", domain, ip)
	}
}

// generateDomain 生成域名
func generateDomain(deviceName, suffix string) string {
	// 将设备名转换为小写，替换空格为连字符
	name := strings.ToLower(deviceName)
	name = strings.ReplaceAll(name, " ", "-")
	return fmt.Sprintf("%s.%s", name, suffix)
}

// hasDomainConflict 检查域名冲突
func hasDomainConflict(manager *node.Manager, domain, deviceID string) bool {
	nodes := manager.GetAll()
	for _, n := range nodes {
		if n.Domain == domain && n.DeviceID != deviceID {
			return true
		}
	}
	return false
}

// extractMACShort 提取MAC地址的短格式（后6位）
func extractMACShort(deviceID string) string {
	// deviceID格式: mac-00:11:22:33:44:55
	parts := strings.Split(deviceID, "-")
	if len(parts) != 2 {
		return deviceID
	}
	mac := strings.ReplaceAll(parts[1], ":", "")
	if len(mac) >= 6 {
		return mac[len(mac)-6:]
	}
	return mac
}

