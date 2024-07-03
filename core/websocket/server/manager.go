package server

import "sync"

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

func (m *Manager) Store(uuid string, c *Client) {
	m.clients.Store(uuid, c)
}

func (m *Manager) Delete(uuid string) {
	m.clients.Delete(uuid)
}

func (m *Manager) IsExist(uuid string) bool {
	_, ok := m.clients.Load(uuid)
	return ok
}

func (m *Manager) Range(f func(key, value any) bool) {
	m.clients.Range(f)
}
