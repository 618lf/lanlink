package node

import (
	"sync"
	"time"
)

// Node 节点信息
type Node struct {
	DeviceID  string    // 设备ID（MAC地址）
	Domain    string    // 域名
	IP        string    // IP地址
	Hostname  string    // 主机名
	LastSeen  time.Time // 最后心跳时间
	IsLocal   bool      // 是否是本机节点
}

// Manager 节点管理器
type Manager struct {
	mu            sync.RWMutex
	nodes         map[string]*Node // key: deviceID
	offlineTimeout time.Duration
	onNodeChange  func(*Node, bool) // 节点变化回调：(node, isOnline)
}

// NewManager 创建节点管理器
func NewManager(offlineTimeout time.Duration) *Manager {
	return &Manager{
		nodes:          make(map[string]*Node),
		offlineTimeout: offlineTimeout,
	}
}

// SetChangeCallback 设置节点变化回调
func (m *Manager) SetChangeCallback(callback func(*Node, bool)) {
	m.onNodeChange = callback
}

// AddOrUpdate 添加或更新节点
func (m *Manager) AddOrUpdate(deviceID, domain, ip, hostname string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[deviceID]
	now := time.Now()

	if !exists {
		// 新节点
		node = &Node{
			DeviceID: deviceID,
			Domain:   domain,
			IP:       ip,
			Hostname: hostname,
			LastSeen: now,
			IsLocal:  false,
		}
		m.nodes[deviceID] = node

		// 触发回调
		if m.onNodeChange != nil {
			m.onNodeChange(node, true)
		}
		return true
	}

	// 更新现有节点
	changed := false
	if node.IP != ip {
		node.IP = ip
		changed = true
	}
	if node.Domain != domain {
		node.Domain = domain
		changed = true
	}
	if node.Hostname != hostname {
		node.Hostname = hostname
		changed = true
	}
	node.LastSeen = now

	// 如果IP或域名变化，触发回调
	if changed && m.onNodeChange != nil {
		m.onNodeChange(node, true)
	}

	return changed
}

// Remove 移除节点
func (m *Manager) Remove(deviceID string) *Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, exists := m.nodes[deviceID]
	if !exists {
		return nil
	}

	delete(m.nodes, deviceID)

	// 触发回调
	if m.onNodeChange != nil {
		m.onNodeChange(node, false)
	}

	return node
}

// Get 获取节点
func (m *Manager) Get(deviceID string) (*Node, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	node, exists := m.nodes[deviceID]
	return node, exists
}

// GetAll 获取所有节点
func (m *Manager) GetAll() []*Node {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nodes := make([]*Node, 0, len(m.nodes))
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// CheckOffline 检查离线节点
func (m *Manager) CheckOffline() []*Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	offlineNodes := make([]*Node, 0)

	for deviceID, node := range m.nodes {
		// 跳过本机节点
		if node.IsLocal {
			continue
		}

		// 检查是否超时
		if now.Sub(node.LastSeen) > m.offlineTimeout {
			offlineNodes = append(offlineNodes, node)
			delete(m.nodes, deviceID)

			// 触发回调
			if m.onNodeChange != nil {
				m.onNodeChange(node, false)
			}
		}
	}

	return offlineNodes
}

// SetLocal 设置本机节点
func (m *Manager) SetLocal(deviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node, exists := m.nodes[deviceID]; exists {
		node.IsLocal = true
	}
}

