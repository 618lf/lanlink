package hosts

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	// Hosts文件标记
	beginMarker = "# === LanLink Managed Begin ==="
	endMarker   = "# === LanLink Managed End ==="
	entryMarker = "# LanLink"
)

// Manager Hosts文件管理器
type Manager struct {
	hostsPath string
}

// NewManager 创建Hosts管理器
func NewManager() *Manager {
	return &Manager{
		hostsPath: getHostsPath(),
	}
}

// CheckPermission 检查是否有权限修改hosts文件
func (m *Manager) CheckPermission() error {
	// 尝试打开文件
	file, err := os.OpenFile(m.hostsPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("无权限修改hosts文件，请以管理员/root身份运行: %v", err)
	}
	file.Close()
	return nil
}

// Initialize 初始化hosts文件（添加标记区域）
func (m *Manager) Initialize() error {
	content, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// 如果已经有标记，不需要重复初始化
	if strings.Contains(contentStr, beginMarker) {
		return nil
	}

	// 备份
	if err := m.backup(); err != nil {
		return err
	}

	// 在文件末尾添加标记区域
	file, err := os.OpenFile(m.hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 确保文件末尾有换行
	if !strings.HasSuffix(contentStr, "\n") {
		file.WriteString("\n")
	}

	// 添加标记区域
	file.WriteString(fmt.Sprintf("\n%s\n%s\n", beginMarker, endMarker))

	return nil
}

// AddOrUpdate 添加或更新域名映射
func (m *Manager) AddOrUpdate(ip, domain string) error {
	if err := m.backup(); err != nil {
		return err
	}

	content, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines))
	inManagedZone := false
	found := false
	entry := fmt.Sprintf("%s\t%s\t%s", ip, domain, entryMarker)

	for _, line := range lines {
		if strings.TrimSpace(line) == beginMarker {
			inManagedZone = true
			newLines = append(newLines, line)
			continue
		}

		if strings.TrimSpace(line) == endMarker {
			// 如果在管理区域内没找到，添加新条目
			if inManagedZone && !found {
				newLines = append(newLines, entry)
			}
			inManagedZone = false
			newLines = append(newLines, line)
			continue
		}

		// 在管理区域内，检查是否是要更新的域名
		if inManagedZone && strings.Contains(line, entryMarker) {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[1] == domain {
				// 找到了，更新IP
				newLines = append(newLines, entry)
				found = true
				continue
			}
		}

		newLines = append(newLines, line)
	}

	// 写回文件
	return os.WriteFile(m.hostsPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// Remove 删除域名映射
func (m *Manager) Remove(domain string) error {
	if err := m.backup(); err != nil {
		return err
	}

	content, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines))
	inManagedZone := false

	for _, line := range lines {
		if strings.TrimSpace(line) == beginMarker {
			inManagedZone = true
			newLines = append(newLines, line)
			continue
		}

		if strings.TrimSpace(line) == endMarker {
			inManagedZone = false
			newLines = append(newLines, line)
			continue
		}

		// 在管理区域内，检查是否是要删除的域名
		if inManagedZone && strings.Contains(line, entryMarker) {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[1] == domain {
				// 跳过这一行（删除）
				continue
			}
		}

		newLines = append(newLines, line)
	}

	// 写回文件
	return os.WriteFile(m.hostsPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// List 列出所有LanLink管理的条目
func (m *Manager) List() (map[string]string, error) {
	content, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	inManagedZone := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == beginMarker {
			inManagedZone = true
			continue
		}

		if strings.TrimSpace(line) == endMarker {
			inManagedZone = false
			continue
		}

		if inManagedZone && strings.Contains(line, entryMarker) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				entries[fields[1]] = fields[0] // domain -> ip
			}
		}
	}

	return entries, scanner.Err()
}

// backup 备份hosts文件
func (m *Manager) backup() error {
	backupPath := m.hostsPath + ".bak"

	content, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, content, 0644)
}

// getHostsPath 获取hosts文件路径
func getHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Windows\System32\drivers\etc\hosts`
	default: // linux, darwin
		return "/etc/hosts"
	}
}
