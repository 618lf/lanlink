//go:build windows

package hardware

import (
	"fmt"
	"os/exec"
	"strings"
)

// getSerialNumber Windows 平台获取硬件序列号
func getSerialNumber() (string, error) {
	// 1. 尝试获取 BIOS 序列号
	serial, err := runWmic("bios", "get", "serialnumber")
	if err == nil && isValidSerial(serial) {
		return serial, nil
	}

	// 2. 尝试获取主板序列号
	serial, err = runWmic("baseboard", "get", "serialnumber")
	if err == nil && isValidSerial(serial) {
		return serial, nil
	}

	// 3. 尝试获取 UUID
	serial, err = runWmic("csproduct", "get", "uuid")
	if err == nil && isValidSerial(serial) {
		// 移除 UUID 中的连字符
		serial = strings.ReplaceAll(serial, "-", "")
		return serial, nil
	}

	return "", fmt.Errorf("无法获取硬件序列号")
}

// runWmic 执行 wmic 命令并返回结果
func runWmic(args ...string) (string, error) {
	cmd := exec.Command("wmic", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 解析输出，跳过标题行
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("wmic 输出格式错误")
	}

	// 取第二行（数据行），去除空白字符
	result := strings.TrimSpace(lines[1])
	return result, nil
}

// isValidSerial 检查序列号是否有效
func isValidSerial(serial string) bool {
	if serial == "" {
		return false
	}

	// 排除常见的无效值
	invalidValues := []string{
		"To Be Filled By O.E.M.",
		"To be filled by O.E.M.",
		"Default string",
		"None",
		"N/A",
		"Not Specified",
		"System Serial Number",
		"0",
	}

	for _, invalid := range invalidValues {
		if strings.EqualFold(serial, invalid) {
			return false
		}
	}

	return true
}

