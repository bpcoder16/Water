package server

import "sync"

// Manager TODO 启动定时服务去定时清理无效的 client
type Manager struct {
	clients *sync.Map
}

func NewManager() *Manager {
	return &Manager{
		clients: new(sync.Map),
	}
}

func (m *Manager) Len() int {
	cnt := 0
	m.clients.Range(func(k, v interface{}) bool {
		cnt++
		return true
	})
	return cnt
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

func (m *Manager) Range(f func(key, value any) bool) {
	m.clients.Range(f)
}
