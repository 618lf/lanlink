package config

import (
	"encoding/json"
	"os"
)

// Config 应用配置
type Config struct {
	DeviceName        string `json:"deviceName"`        // 设备名，空则自动获取主机名
	DomainSuffix      string `json:"domainSuffix"`      // 域名后缀
	MulticastAddr     string `json:"multicastAddr"`     // 组播地址
	MulticastPort     int    `json:"multicastPort"`     // 组播端口
	HeartbeatInterval int    `json:"heartbeatInterval"` // 心跳间隔（秒）
	OfflineTimeout    int    `json:"offlineTimeout"`    // 离线超时（秒）
	LogLevel          string `json:"logLevel"`          // 日志级别
}

// Default 默认配置
func Default() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		DeviceName:        hostname,
		DomainSuffix:      "local",
		MulticastAddr:     "239.255.0.1",
		MulticastPort:     9527,
		HeartbeatInterval: 10,
		OfflineTimeout:    30,
		LogLevel:          "info",
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

	// 如果配置中deviceName为空，使用主机名
	if cfg.DeviceName == "" {
		hostname, _ := os.Hostname()
		cfg.DeviceName = hostname
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

