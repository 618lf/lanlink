package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

// Action 消息动作类型
const (
	ActionHeartbeat = "heartbeat"
	ActionOffline   = "offline"
)

// Message 组播消息
type Message struct {
	Action    string `json:"action"`    // heartbeat/offline
	Domain    string `json:"domain"`    // 域名
	IP        string `json:"ip"`        // IP地址
	DeviceID  string `json:"deviceId"`  // 设备ID
	Hostname  string `json:"hostname"`  // 主机名
	Timestamp int64  `json:"timestamp"` // 时间戳
}

// MulticastClient 组播客户端
type MulticastClient struct {
	addr       string
	port       int
	conn       *net.UDPConn
	packetConn *ipv4.PacketConn
	localIP    string
	onMessage  func(*Message) // 消息接收回调
}

// NewMulticastClient 创建组播客户端
func NewMulticastClient(addr string, port int) (*MulticastClient, error) {
	localIP, err := getLocalIP()
	if err != nil {
		return nil, fmt.Errorf("获取本机IP失败: %v", err)
	}

	return &MulticastClient{
		addr:    addr,
		port:    port,
		localIP: localIP,
	}, nil
}

// SetMessageCallback 设置消息回调
func (c *MulticastClient) SetMessageCallback(callback func(*Message)) {
	c.onMessage = callback
}

// Start 启动组播监听
func (c *MulticastClient) Start() error {
	groupAddr := &net.UDPAddr{
		IP:   net.ParseIP(c.addr),
		Port: c.port,
	}

	// 监听所有接口
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: c.port,
	})
	if err != nil {
		return fmt.Errorf("创建UDP连接失败: %v", err)
	}

	c.conn = conn
	c.packetConn = ipv4.NewPacketConn(conn)

	// 加入组播组
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("获取网络接口失败: %v", err)
	}

	// 尝试在所有可用接口上加入组播组
	joined := false
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}
		if err := c.packetConn.JoinGroup(&iface, groupAddr); err == nil {
			joined = true
		}
	}

	if !joined {
		return fmt.Errorf("无法加入组播组")
	}

	// 启动接收协程
	go c.receiveLoop()

	return nil
}

// Send 发送消息
func (c *MulticastClient) Send(msg *Message) error {
	msg.Timestamp = time.Now().Unix()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(c.addr),
		Port: c.port,
	}

	_, err = c.conn.WriteToUDP(data, addr)
	return err
}

// Close 关闭连接
func (c *MulticastClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetLocalIP 获取本机IP
func (c *MulticastClient) GetLocalIP() string {
	return c.localIP
}

// receiveLoop 接收消息循环
func (c *MulticastClient) receiveLoop() {
	buffer := make([]byte, 1024)

	for {
		n, _, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			continue
		}

		// 忽略自己发送的消息
		if msg.IP == c.localIP {
			continue
		}

		// 触发回调
		if c.onMessage != nil {
			c.onMessage(&msg)
		}
	}
}

// getLocalIP 获取本机局域网IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的局域网IP")
}

// GetMACAddress 获取MAC地址作为设备ID
func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过回环接口和没有MAC地址的接口
		if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
			continue
		}

		// 返回第一个有效的MAC地址
		mac := iface.HardwareAddr.String()
		if mac != "" {
			return "mac-" + mac, nil
		}
	}

	return "", fmt.Errorf("未找到有效的MAC地址")
}

