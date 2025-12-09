//go:build darwin

package hardware

import (
	"fmt"
	"os/exec"
	"strings"
)

// getSerialNumber macOS 平台获取硬件序列号
func getSerialNumber() (string, error) {
	// 方法1: 使用 system_profiler 获取硬件序列号
	serial, err := getSerialFromSystemProfiler()
	if err == nil && serial != "" {
		return serial, nil
	}

	// 方法2: 使用 ioreg 获取 IOPlatformSerialNumber
	serial, err = getSerialFromIoreg()
	if err == nil && serial != "" {
		return serial, nil
	}

	return "", fmt.Errorf("无法获取硬件序列号")
}

// getSerialFromSystemProfiler 从 system_profiler 获取序列号
func getSerialFromSystemProfiler() (string, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Serial Number") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("未找到序列号")
}

// getSerialFromIoreg 从 ioreg 获取序列号
func getSerialFromIoreg() (string, error) {
	cmd := exec.Command("ioreg", "-l")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IOPlatformSerialNumber") {
			// 格式: "IOPlatformSerialNumber" = "XXXXXXXXXXXX"
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				return parts[3], nil
			}
		}
	}

	return "", fmt.Errorf("未找到序列号")
}

