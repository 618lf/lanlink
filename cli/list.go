package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/618lf/lanlink/internal"
)

// ListNodes 列出所有节点
func ListNodes(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	onlineOnly := fs.Bool("online", false, "仅显示在线节点")
	offlineOnly := fs.Bool("offline", false, "仅显示离线节点")
	jsonOutput := fs.Bool("json", false, "JSON格式输出")
	fs.Parse(args)

	nodes, err := internal.GetNodes()
	if err != nil {
		Error("获取节点列表失败: %v", err)
		return err
	}

	// 过滤节点
	filtered := make([]internal.NodeInfo, 0)
	online := 0
	offline := 0

	for _, node := range nodes {
		if *onlineOnly && node.Status != "online" {
			continue
		}
		if *offlineOnly && node.Status == "online" {
			continue
		}
		filtered = append(filtered, node)

		if node.Status == "online" {
			online++
		} else {
			offline++
		}
	}

	// JSON 输出
	if *jsonOutput {
		output := map[string]interface{}{
			"nodes": filtered,
			"summary": map[string]int{
				"online":  online,
				"offline": offline,
				"total":   len(nodes),
			},
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(output)
	}

	// 表格输出
	Header(fmt.Sprintf("节点列表 (在线: %d, 离线: %d)", online, offline))

	if len(filtered) == 0 {
		Warn("未发现任何节点")
		Footer()
		return nil
	}

	// 表头
	fmt.Printf("%-6s %-20s %-17s %-15s %-20s\n",
		"状态", "域名", "IP 地址", "主机名", "最后心跳")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 节点列表
	for _, node := range filtered {
		statusIcon := "✗"
		if node.Status == "online" {
			statusIcon = color(ColorGreen, "✓")
		} else {
			statusIcon = color(ColorRed, "✗")
		}

		hostname := node.Hostname
		if hostname == "" {
			hostname = "-"
		}

		lastSeen := formatLastSeen(node.LastSeen)
		if node.IsLocal {
			lastSeen = "(本机)"
		}

		fmt.Printf("%-6s %-20s %-17s %-15s %-20s\n",
			statusIcon,
			node.Domain,
			node.IP,
			hostname,
			lastSeen)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n共 %d 个节点\n\n", len(filtered))

	return nil
}

// formatLastSeen 格式化最后心跳时间
func formatLastSeen(t time.Time) string {
	if t.IsZero() {
		return "未知"
	}

	elapsed := time.Since(t)
	if elapsed < time.Minute {
		return fmt.Sprintf("%.0f秒前", elapsed.Seconds())
	}
	if elapsed < time.Hour {
		return fmt.Sprintf("%.0f分钟前", elapsed.Minutes())
	}
	return fmt.Sprintf("%.1f小时前", elapsed.Hours())
}

