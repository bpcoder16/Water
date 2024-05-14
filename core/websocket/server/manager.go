package server

import "sync"

// clientManager TODO 启动定时服务去定时清理无效的 client
var clientManager *Manager

type Manager struct {
	clients *sync.Map
}

func NewManager() *Manager {
	return &Manager{
		clients: new(sync.Map),
	}
}

func (m *Manager) Store(c *Client) {
	m.clients.Store(c, struct{}{})
}

func (m *Manager) Delete(c *Client) {
	m.clients.Delete(c)
}

func (m *Manager) IsExist(c *Client) bool {
	_, ok := m.clients.Load(c)
	return ok
}
