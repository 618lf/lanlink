package hardware

import (
	"regexp"
	"runtime"
	"strings"
)

// GetPlatform 获取当前平台标识
func GetPlatform() string {
	switch runtime.GOOS {
	case "windows":
		return "win"
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	default:
		return "unknown"
	}
}

// GetSerialNumber 获取硬件序列号（跨平台）
// 优先级：
// - Windows: BIOS序列号 → 主板序列号 → UUID
// - macOS: 硬件序列号
// - Linux: 产品序列号 → machine-id
func GetSerialNumber() (string, error) {
	return getSerialNumber()
}

// GenerateHostID 生成主机标识（序列号后6位）
func GenerateHostID(serial string) string {
	// 转小写，移除非字母数字字符
	clean := sanitize(serial)

	// 取后6位
	if len(clean) > 6 {
		clean = clean[len(clean)-6:]
	}

	// 不足6位用0填充
	for len(clean) < 6 {
		clean = "0" + clean
	}

	return clean
}

// GenerateHostname 生成完整的主机域名
// 格式: {platform}-{序列号后6位}.{suffix}
func GenerateHostname(suffix string) (string, error) {
	platform := GetPlatform()
	serial, err := GetSerialNumber()
	if err != nil {
		return "", err
	}

	hostID := GenerateHostID(serial)
	return platform + "-" + hostID + "." + suffix, nil
}

// GenerateDeviceName 生成设备名称（不含后缀）
// 格式: {platform}-{序列号后6位}
func GenerateDeviceName() (string, error) {
	platform := GetPlatform()
	serial, err := GetSerialNumber()
	if err != nil {
		return "", err
	}

	hostID := GenerateHostID(serial)
	return platform + "-" + hostID, nil
}

// sanitize 清理字符串，只保留字母数字，转小写
func sanitize(s string) string {
	// 移除非字母数字字符
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)
	clean := reg.ReplaceAllString(s, "")

	// 转小写
	return strings.ToLower(clean)
}

