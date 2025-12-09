//go:build linux

package hardware

import (
	"fmt"
	"os"
	"strings"
)

// getSerialNumber Linux 平台获取硬件序列号
func getSerialNumber() (string, error) {
	// 1. 尝试从 DMI 获取产品序列号
	serial, err := readFile("/sys/class/dmi/id/product_serial")
	if err == nil && isValidSerial(serial) {
		return serial, nil
	}

	// 2. 尝试获取主板序列号
	serial, err = readFile("/sys/class/dmi/id/board_serial")
	if err == nil && isValidSerial(serial) {
		return serial, nil
	}

	// 3. 使用 machine-id 作为后备方案
	serial, err = readFile("/etc/machine-id")
	if err == nil && serial != "" {
		return serial, nil
	}

	// 4. 尝试 /var/lib/dbus/machine-id
	serial, err = readFile("/var/lib/dbus/machine-id")
	if err == nil && serial != "" {
		return serial, nil
	}

	return "", fmt.Errorf("无法获取硬件序列号")
}

// readFile 读取文件内容
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
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

