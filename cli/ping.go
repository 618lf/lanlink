package cli

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"

	"github.com/618lf/lanlink/internal"
)

// PingNode 测试节点连接
func PingNode(args []string) error {
	if len(args) == 0 {
		Error("缺少参数 <domain>")
		fmt.Println("用法: lanlink ping <domain>")
		fmt.Println("示例: lanlink ping server.local")
		return fmt.Errorf("missing domain argument")
	}

	domain := args[0]

	Header(fmt.Sprintf("测试连接: %s", domain))

	// 1. DNS 解析测试
	Section("DNS 解析")
	ips, err := net.LookupHost(domain)
	if err != nil {
		Error("解析失败: %v", err)
	} else if len(ips) > 0 {
		Success("解析成功: %s -> %s", domain, ips[0])
	}

	// 2. Hosts 记录检查
	Section("Hosts 记录")
	nodes, err := internal.GetNodes()
	found := false
	var targetIP string
	if err == nil {
		for _, node := range nodes {
			if node.Domain == domain {
				found = true
				targetIP = node.IP
				Success("找到记录: %s", node.IP)
				KeyValue("状态", node.Status)
				if node.Hostname != "" {
					KeyValue("主机名", node.Hostname)
				}
				break
			}
		}
	}
	if !found {
		Warn("Hosts 文件中未找到记录")
	}

	// 3. Ping 测试
	Section("Ping 测试")
	if len(ips) > 0 || targetIP != "" {
		target := domain
		if targetIP != "" {
			target = targetIP
		}

		if pingHost(target) {
			Success("Ping 成功")
		} else {
			Error("Ping 失败（设备可能离线或防火墙阻止）")
		}
	} else {
		Warn("跳过（无法解析域名）")
	}

	Footer()

	// 总结
	fmt.Println()
	if found && len(ips) > 0 {
		Success("连接正常 ✓")
	} else {
		Warn("连接异常 ✗")
		fmt.Println("\n建议:")
		fmt.Println("  1. 确认目标设备是否运行 LanLink")
		fmt.Println("  2. 检查网络连接")
		fmt.Println("  3. 查看日志: lanlink logs")
	}
	fmt.Println()

	return nil
}

// pingHost 执行 ping 命令
func pingHost(host string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
	}

	err := cmd.Run()
	return err == nil
}

