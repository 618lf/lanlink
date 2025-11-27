package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// RunInteractive 运行交互式命令行
func RunInteractive() {
	// 打印欢迎信息
	printWelcome()

	// 创建输入扫描器
	scanner := bufio.NewScanner(os.Stdin)

	// 主循环
	for {
		// 显示提示符
		printPrompt()

		// 读取用户输入
		if !scanner.Scan() {
			break
		}

		// 获取输入的命令
		input := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if input == "" {
			continue
		}

		// 分割命令和参数
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		// 执行命令
		shouldExit := executeCommand(cmd, args)
		if shouldExit {
			break
		}

		// 打印空行分隔输出
		fmt.Println()
	}

	// 打印退出信息
	printGoodbye()
}

// printWelcome 打印欢迎信息
func printWelcome() {
	clearScreen()
	fmt.Print(color(ColorCyan, "\n╔════════════════════════════════════════════════════════╗\n"))
	fmt.Print(color(ColorCyan, "║          LanLink Interactive Shell                    ║\n"))
	fmt.Print(color(ColorCyan, "╚════════════════════════════════════════════════════════╝\n"))
	fmt.Printf("  版本: %s\n", "1.0.0")
	fmt.Printf("  构建: %s\n\n", "latest")
	fmt.Print(color(ColorYellow+ColorBold, "💡 提示：\n"))
	fmt.Println("  - 输入命令直接执行，无需前缀 'lanlink'")
	fmt.Println("  - 输入 'help' 查看所有命令")
	fmt.Println("  - 输入 'clear' 清屏")
	fmt.Println("  - 输入 'exit' 或 'quit' 退出")
	fmt.Println()
}

// printPrompt 打印命令提示符
func printPrompt() {
	fmt.Print(color(ColorGreen, "lanlink"))
	fmt.Print(color(ColorCyan, "> "))
}

// executeCommand 执行命令
func executeCommand(cmd string, args []string) bool {
	switch strings.ToLower(cmd) {
	case "exit", "quit", "q":
		return true

	case "clear", "cls":
		clearScreen()
		return false

	case "help", "h", "?":
		printInteractiveHelp()
		return false

	case "status", "st":
		ShowStatus()
		return false

	case "list", "ls":
		ListNodes(args)
		return false

	case "logs", "log":
		// 构建参数数组
		logArgs := []string{}
		for _, arg := range args {
			logArgs = append(logArgs, arg)
		}
		ShowLogs(logArgs)
		return false

	case "ping":
		if len(args) == 0 {
			Error("请指定要 ping 的域名")
			fmt.Println("   用法: ping <domain>")
			fmt.Println("   示例: ping server1.local")
		} else {
			PingNode(args)
		}
		return false

	case "version", "v":
		ShowVersion()
		return false

	case "install":
		Warn("install 命令需要管理员权限，建议退出交互模式后执行：")
		if runtime.GOOS == "windows" {
			fmt.Println("   以管理员身份运行: lanlink install")
		} else {
			fmt.Println("   sudo lanlink install")
		}
		return false

	case "uninstall":
		Warn("uninstall 命令需要管理员权限，建议退出交互模式后执行：")
		if runtime.GOOS == "windows" {
			fmt.Println("   以管理员身份运行: lanlink uninstall")
		} else {
			fmt.Println("   sudo lanlink uninstall")
		}
		return false

	case "service", "svc":
		Warn("service 命令需要管理员权限，建议退出交互模式后执行：")
		if runtime.GOOS == "windows" {
			fmt.Println("   以管理员身份运行: lanlink service <subcommand>")
		} else {
			fmt.Println("   sudo lanlink service <subcommand>")
		}
		fmt.Println("\n可用子命令:")
		fmt.Println("   install   - 安装服务")
		fmt.Println("   uninstall - 卸载服务")
		fmt.Println("   start     - 启动服务")
		fmt.Println("   stop      - 停止服务")
		fmt.Println("   status    - 查看服务状态")
		return false

	case "refresh", "reload":
		Info("🔄 刷新配置...")
		fmt.Println("提示: 当前版本需要重启 LanLink 服务来重新加载配置")
		return false

	default:
		Error(fmt.Sprintf("未知命令: %s", cmd))
		fmt.Println("   输入 'help' 查看所有可用命令")
		return false
	}
}

// printInteractiveHelp 打印交互式帮助
func printInteractiveHelp() {
	fmt.Print(color(ColorCyan+ColorBold, "📚 可用命令：\n\n"))

	commands := []struct {
		name    string
		aliases string
		desc    string
	}{
		{"status", "st", "查看 LanLink 运行状态"},
		{"list", "ls", "列出所有已发现的节点"},
		{"logs", "log", "查看日志（-f 实时跟踪，-n 指定行数）"},
		{"ping", "", "测试与指定节点的连接"},
		{"version", "v", "显示版本信息"},
		{"help", "h, ?", "显示此帮助信息"},
		{"clear", "cls", "清屏"},
		{"exit", "quit, q", "退出交互模式"},
	}

	fmt.Println("🔍 查询命令:")
	for _, cmd := range commands[:5] {
		if cmd.aliases != "" {
			fmt.Printf("  %-12s %-12s %s\n", cmd.name, fmt.Sprintf("[%s]", cmd.aliases), cmd.desc)
		} else {
			fmt.Printf("  %-12s %-12s %s\n", cmd.name, "", cmd.desc)
		}
	}

	fmt.Println("\n🛠️  管理命令 (需要管理员权限，建议退出后执行):")
	fmt.Println("  install                   安装到系统 PATH")
	fmt.Println("  uninstall                 从系统 PATH 卸载")
	fmt.Println("  service install           安装为系统服务（开机自启）")
	fmt.Println("  service start             启动服务")
	fmt.Println("  service stop              停止服务")
	fmt.Println("  service status            查看服务状态")

	fmt.Println("\n💡 使用示例:")
	fmt.Println("  lanlink> status           # 查看运行状态")
	fmt.Println("  lanlink> list             # 列出所有节点")
	fmt.Println("  lanlink> logs -f          # 实时查看日志")
	fmt.Println("  lanlink> ping srv1.local  # 测试连接")
	fmt.Println("  lanlink> clear            # 清屏")
	fmt.Println("  lanlink> exit             # 退出")
}

// clearScreen 清屏
func clearScreen() {
	if runtime.GOOS == "windows" {
		// Windows 使用 cls
		fmt.Print("\033[H\033[2J")
	} else {
		// Unix/Linux/Mac 使用 clear
		fmt.Print("\033[H\033[2J")
	}
}

// printGoodbye 打印退出信息
func printGoodbye() {
	fmt.Println()
	fmt.Print(color(ColorCyan+ColorBold, "👋 再见！感谢使用 LanLink\n"))
}

