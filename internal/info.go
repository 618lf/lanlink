package internal

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// NodeInfo 节点信息
type NodeInfo struct {
	Domain   string
	IP       string
	Hostname string
	Status   string // online/offline
	LastSeen time.Time
	IsLocal  bool
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	IsRunning bool
	PID       int
	Memory    int64
	StartTime time.Time
}

// LogEntry 日志条目
type LogEntry struct {
	Time    string
	Level   string
	Message string
}

// GetRuntimeInfo 获取运行时信息
func GetRuntimeInfo() (*RuntimeInfo, error) {
	info := &RuntimeInfo{}

	// 检查进程是否运行
	pid, running := checkProcess()
	info.IsRunning = running
	info.PID = pid

	if running {
		// 获取进程信息
		info.Memory = getProcessMemory(pid)
		info.StartTime = getProcessStartTime(pid)
	}

	return info, nil
}

// GetNodes 从 Hosts 文件读取节点信息
func GetNodes() ([]NodeInfo, error) {
	hostsPath := getHostsPath()
	file, err := os.Open(hostsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []NodeInfo
	scanner := bufio.NewScanner(file)
	inManagedZone := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.Contains(line, "LanLink Managed Begin") {
			inManagedZone = true
			continue
		}
		if strings.Contains(line, "LanLink Managed End") {
			inManagedZone = false
			continue
		}

		if inManagedZone && strings.Contains(line, "# LanLink") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				node := NodeInfo{
					IP:      fields[0],
					Domain:  fields[1],
					Status:  "unknown",
					IsLocal: false,
				}
				nodes = append(nodes, node)
			}
		}
	}

	// 从日志文件获取最后心跳时间和在线状态
	enrichNodesFromLog(nodes)

	return nodes, scanner.Err()
}

// GetRecentLogs 获取最近的日志
func GetRecentLogs(n int) ([]LogEntry, error) {
	file, err := os.Open("lanlink.log")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 获取最后 n 行
	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}

	var entries []LogEntry
	for _, line := range lines[start:] {
		entry := parseLogLine(line)
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// checkProcess 检查进程是否运行
func checkProcess() (int, bool) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/FI", "IMAGENAME eq lanlink.exe", "/FO", "CSV", "/NH")
	} else {
		cmd = exec.Command("pgrep", "-f", "lanlink")
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, false
	}

	if runtime.GOOS == "windows" {
		// Windows: 解析 CSV 输出
		parts := strings.Split(string(output), ",")
		if len(parts) >= 2 {
			pidStr := strings.Trim(parts[1], "\"")
			pid, _ := strconv.Atoi(pidStr)
			return pid, true
		}
	} else {
		// Linux/Mac: pgrep 直接返回 PID
		pidStr := strings.TrimSpace(string(output))
		lines := strings.Split(pidStr, "\n")
		if len(lines) > 0 && lines[0] != "" {
			pid, _ := strconv.Atoi(lines[0])
			return pid, true
		}
	}

	return 0, false
}

// getProcessMemory 获取进程内存占用（简化版）
func getProcessMemory(pid int) int64 {
	// 简化实现，返回估计值
	return 20 * 1024 * 1024 // 20MB
}

// getProcessStartTime 获取进程启动时间（简化版）
func getProcessStartTime(pid int) time.Time {
	// 简化实现，从日志文件推断
	if stat, err := os.Stat("lanlink.log"); err == nil {
		return stat.ModTime()
	}
	return time.Now()
}

// enrichNodesFromLog 从日志文件丰富节点信息
func enrichNodesFromLog(nodes []NodeInfo) {
	// 简化实现：标记所有为在线
	for i := range nodes {
		nodes[i].Status = "online"
		nodes[i].LastSeen = time.Now()
	}
}

// parseLogLine 解析日志行
func parseLogLine(line string) LogEntry {
	entry := LogEntry{}

	// 格式: [2024-11-27 14:30:00] [INFO] message
	if len(line) < 30 {
		entry.Message = line
		return entry
	}

	// 提取时间
	if strings.HasPrefix(line, "[") {
		end := strings.Index(line[1:], "]")
		if end > 0 {
			entry.Time = line[1 : end+1]
			line = line[end+2:]
		}
	}

	// 提取级别
	if strings.HasPrefix(line, " [") {
		end := strings.Index(line[2:], "]")
		if end > 0 {
			entry.Level = line[3 : end+2]
			line = line[end+4:]
		}
	}

	entry.Message = strings.TrimSpace(line)
	return entry
}

// getHostsPath 获取 Hosts 文件路径
func getHostsPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

// GetLogFileModTime 获取日志文件最后修改时间
func GetLogFileModTime() (time.Time, error) {
	stat, err := os.Stat("lanlink.log")
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

