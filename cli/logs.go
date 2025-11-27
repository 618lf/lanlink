package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/618lf/lanlink/internal"
)

// ShowLogs 显示日志
func ShowLogs(args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)
	lines := fs.Int("n", 50, "显示最后N行")
	follow := fs.Bool("f", false, "实时跟踪日志")
	level := fs.String("level", "", "过滤日志级别 (debug/info/warn/error)")
	fs.Parse(args)

	if *follow {
		return followLogs(*level)
	}

	logs, err := internal.GetRecentLogs(*lines)
	if err != nil {
		Error("读取日志失败: %v", err)
		return err
	}

	Header(fmt.Sprintf("最近 %d 行日志", len(logs)))

	for _, log := range logs {
		// 过滤日志级别
		if *level != "" && !strings.EqualFold(log.Level, *level) {
			continue
		}

		printLogEntry(log)
	}

	Footer()
	return nil
}

// followLogs 实时跟踪日志
func followLogs(levelFilter string) error {
	file, err := os.Open("lanlink.log")
	if err != nil {
		Error("打开日志文件失败: %v", err)
		return err
	}
	defer file.Close()

	// 跳到文件末尾
	file.Seek(0, 2)

	Info("实时跟踪日志... (按 Ctrl+C 退出)\n")

	scanner := bufio.NewScanner(file)
	for {
		if scanner.Scan() {
			line := scanner.Text()
			entry := internal.LogEntry{Message: line}

			// 简单解析
			if strings.Contains(line, "[INFO]") {
				entry.Level = "INFO"
			} else if strings.Contains(line, "[WARN]") {
				entry.Level = "WARN"
			} else if strings.Contains(line, "[ERROR]") {
				entry.Level = "ERROR"
			} else if strings.Contains(line, "[DEBUG]") {
				entry.Level = "DEBUG"
			}

			// 过滤级别
			if levelFilter != "" && !strings.EqualFold(entry.Level, levelFilter) {
				continue
			}

			fmt.Println(line)
		} else {
			// 没有新行，等待一会
			time.Sleep(100 * time.Millisecond)

			// 检查文件是否有新内容
			stat, _ := os.Stat("lanlink.log")
			currentPos, _ := file.Seek(0, 1)
			if stat.Size() < currentPos {
				// 文件被截断，重新打开
				file.Close()
				file, _ = os.Open("lanlink.log")
				scanner = bufio.NewScanner(file)
			}
		}
	}
}

// printLogEntry 打印日志条目
func printLogEntry(log internal.LogEntry) {
	levelColor := ColorReset
	switch log.Level {
	case "ERROR":
		levelColor = ColorRed
	case "WARN":
		levelColor = ColorYellow
	case "INFO":
		levelColor = ColorGreen
	case "DEBUG":
		levelColor = ColorGray
	}

	if log.Time != "" {
		fmt.Printf("[%s] ", color(ColorGray, log.Time))
	}
	if log.Level != "" {
		fmt.Printf("[%s] ", color(levelColor, log.Level))
	}
	fmt.Println(log.Message)
}

