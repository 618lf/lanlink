package config

import (
	"encoding/json"
	"os"

	"github.com/618lf/lanlink/hardware"
)

// Config 应用配置
type Config struct {
	DeviceName           string `json:"deviceName"`           // 设备名，空则基于硬件ID自动生成
	DomainSuffix         string `json:"domainSuffix"`         // 域名后缀
	MulticastAddr        string `json:"multicastAddr"`        // 组播地址
	MulticastPort        int    `json:"multicastPort"`        // 组播端口
	HeartbeatIntervalSec int    `json:"heartbeatIntervalSec"` // 心跳间隔（秒）
	OfflineTimeoutSec    int    `json:"offlineTimeoutSec"`    // 离线超时（秒）
	LogLevel             string `json:"logLevel"`             // 日志级别
}

// Default 默认配置
func Default() *Config {
	// 使用硬件ID生成设备名，格式: {platform}-{序列号后6位}
	deviceName, err := hardware.GenerateDeviceName()
	if err != nil {
		// 如果获取硬件ID失败，回退到主机名
		deviceName, _ = os.Hostname()
	}

	return &Config{
		DeviceName:           deviceName,
		DomainSuffix:         "coobee.local",
		MulticastAddr:        "239.255.0.1",
		MulticastPort:        9527,
		HeartbeatIntervalSec: 10,
		OfflineTimeoutSec:    30,
		LogLevel:             "info",
	}
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	cfg := Default()

	// 如果文件不存在，使用默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 如果配置中deviceName为空，基于硬件ID自动生成
	if cfg.DeviceName == "" {
		deviceName, err := hardware.GenerateDeviceName()
		if err != nil {
			// 如果获取硬件ID失败，回退到主机名
			deviceName, _ = os.Hostname()
		}
		cfg.DeviceName = deviceName
	}

	return cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

