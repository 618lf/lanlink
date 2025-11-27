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
		
		// Remove UTF-8 BOM if present
		input = strings.TrimPrefix(input, "\ufeff")
		
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

// printWelcome prints welcome message
func printWelcome() {
	clearScreen()
	fmt.Println("\n========================================================")
	fmt.Println("          LanLink Interactive Shell")
	fmt.Println("========================================================")
	fmt.Printf("  Version: %s\n", "1.0.0")
	fmt.Printf("  Build: %s\n\n", "latest")
	fmt.Println("Tips:")
	fmt.Println("  - Enter commands directly without 'lanlink' prefix")
	fmt.Println("  - Type 'help' to see all commands")
	fmt.Println("  - Type 'clear' to clear screen")
	fmt.Println("  - Type 'exit' or 'quit' to exit")
	fmt.Println()
}

// printPrompt prints command prompt
func printPrompt() {
	fmt.Print("lanlink> ")
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
			fmt.Println("[ERROR] Please specify domain to ping")
			fmt.Println("   Usage: ping <domain>")
			fmt.Println("   Example: ping server1.local")
		} else {
			PingNode(args)
		}
		return false

	case "version", "v":
		ShowVersion()
		return false

	case "install":
		fmt.Println("[WARN] install requires admin privileges, please exit and run:")
		if runtime.GOOS == "windows" {
			fmt.Println("   Run as administrator: lanlink install")
		} else {
			fmt.Println("   sudo lanlink install")
		}
		return false

	case "uninstall":
		fmt.Println("[WARN] uninstall requires admin privileges, please exit and run:")
		if runtime.GOOS == "windows" {
			fmt.Println("   Run as administrator: lanlink uninstall")
		} else {
			fmt.Println("   sudo lanlink uninstall")
		}
		return false

	case "service", "svc":
		fmt.Println("[WARN] service commands require admin privileges, please exit and run:")
		if runtime.GOOS == "windows" {
			fmt.Println("   Run as administrator: lanlink service <subcommand>")
		} else {
			fmt.Println("   sudo lanlink service <subcommand>")
		}
		fmt.Println("\nAvailable subcommands:")
		fmt.Println("   install   - Install service")
		fmt.Println("   uninstall - Uninstall service")
		fmt.Println("   start     - Start service")
		fmt.Println("   stop      - Stop service")
		fmt.Println("   status    - Check service status")
		return false

	case "refresh", "reload":
		fmt.Println("[INFO] Refreshing configuration...")
		fmt.Println("Note: Current version needs to restart LanLink service to reload config")
		return false

	default:
		fmt.Printf("[ERROR] Unknown command: %s\n", cmd)
		fmt.Println("   Type 'help' to see all available commands")
		return false
	}
}

// printInteractiveHelp prints interactive help
func printInteractiveHelp() {
	fmt.Println("\nAvailable Commands:")

	commands := []struct {
		name    string
		aliases string
		desc    string
	}{
		{"status", "st", "Show LanLink running status"},
		{"list", "ls", "List all discovered nodes"},
		{"logs", "log", "View logs (-f to follow, -n for lines)"},
		{"ping", "", "Test connection to a node"},
		{"version", "v", "Show version information"},
		{"help", "h, ?", "Show this help message"},
		{"clear", "cls", "Clear screen"},
		{"exit", "quit, q", "Exit interactive mode"},
	}

	fmt.Println("Query Commands:")
	for _, cmd := range commands[:5] {
		if cmd.aliases != "" {
			fmt.Printf("  %-12s %-12s %s\n", cmd.name, fmt.Sprintf("[%s]", cmd.aliases), cmd.desc)
		} else {
			fmt.Printf("  %-12s %-12s %s\n", cmd.name, "", cmd.desc)
		}
	}

	fmt.Println("\nManagement Commands (require admin privileges, exit to run):")
	fmt.Println("  install                   Install to system PATH")
	fmt.Println("  uninstall                 Uninstall from system PATH")
	fmt.Println("  service install           Install as system service (auto-start)")
	fmt.Println("  service start             Start service")
	fmt.Println("  service stop              Stop service")
	fmt.Println("  service status            Check service status")

	fmt.Println("\nExamples:")
	fmt.Println("  lanlink> status           # Check running status")
	fmt.Println("  lanlink> list             # List all nodes")
	fmt.Println("  lanlink> logs -f          # Follow logs in real-time")
	fmt.Println("  lanlink> ping srv1.local  # Test connection")
	fmt.Println("  lanlink> clear            # Clear screen")
	fmt.Println("  lanlink> exit             # Exit")
}

// clearScreen clears the screen
func clearScreen() {
	// Just print some newlines instead of ANSI codes for better compatibility
	fmt.Print("\n\n\n")
}

// printGoodbye prints goodbye message
func printGoodbye() {
	fmt.Println()
	fmt.Println("Goodbye! Thanks for using LanLink")
}

