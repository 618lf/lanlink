package cli

import (
	"fmt"
	"time"

	"github.com/618lf/lanlink/config"
	"github.com/618lf/lanlink/internal"
)

// ShowStatus 显示运行状态
func ShowStatus() error {
	Header("LanLink 状态概览")

	// 1. 运行状态
	Section("运行状态")
	info, err := internal.GetRuntimeInfo()
	if err != nil {
		Error("获取运行信息失败: %v", err)
		return err
	}

	if info.IsRunning {
		Success("进程运行中 (PID: %d)", info.PID)
		KeyValue("内存占用", fmt.Sprintf("%.1f MB", float64(info.Memory)/1024/1024))
		if !info.StartTime.IsZero() {
			uptime := time.Since(info.StartTime)
			KeyValue("运行时长", formatDuration(uptime))
		}
	} else {
		Warn("进程未运行")
	}

	// 2. 日志状态
	Section("日志状态")
	modTime, err := internal.GetLogFileModTime()
	if err != nil {
		Warn("日志文件不存在")
	} else {
		secondsAgo := time.Since(modTime).Seconds()
		Success("日志文件存在")
		KeyValue("最后更新", fmt.Sprintf("%.0f 秒前", secondsAgo))

		if secondsAgo < 30 {
			Success("日志活跃（系统正常运行）")
		} else {
			Warn("日志停滞（可能异常）")
		}
	}

	// 3. 配置信息
	Section("配置信息")
	cfg, err := config.Load("config.json")
	if err == nil {
		KeyValue("设备名", cfg.DeviceName)
		KeyValue("域名后缀", cfg.DomainSuffix)
		KeyValue("组播地址", fmt.Sprintf("%s:%d", cfg.MulticastAddr, cfg.MulticastPort))
		KeyValue("心跳间隔", fmt.Sprintf("%d 秒", cfg.HeartbeatIntervalSec))
		KeyValue("离线超时", fmt.Sprintf("%d 秒", cfg.OfflineTimeoutSec))
	}

	// 4. 节点统计
	Section("节点统计")
	nodes, err := internal.GetNodes()
	if err != nil {
		Warn("无法读取节点信息: %v", err)
	} else {
		online := 0
		for _, node := range nodes {
			if node.Status == "online" {
				online++
			}
		}
		KeyValue("在线节点", fmt.Sprintf("%d 个", online))
		KeyValue("总节点", fmt.Sprintf("%d 个", len(nodes)))
	}

	// 5. 最近活动
	Section("最近活动")
	logs, err := internal.GetRecentLogs(5)
	if err == nil && len(logs) > 0 {
		for _, log := range logs {
			if log.Level == "INFO" || log.Level == "WARN" || log.Level == "ERROR" {
				fmt.Printf("  [%s] %s\n", log.Time, log.Message)
			}
		}
	} else {
		Info("  无最近活动")
	}

	Footer()

	// 总结
	fmt.Println()
	if info.IsRunning {
		Success("系统运行正常")
	} else {
		Warn("系统未运行")
	}
	fmt.Println()

	return nil
}

// formatDuration 格式化时长
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f 秒", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f 分钟", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f 小时", d.Hours())
	}
	return fmt.Sprintf("%.1f 天", d.Hours()/24)
}

