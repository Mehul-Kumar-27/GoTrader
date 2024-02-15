package routes

import "sync"

type Manager struct {
	ClientsList NseClientList
	mu          sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		ClientsList: make(NseClientList),
		mu:          sync.RWMutex{},
	}
	return m
}

func (m *Manager) AddClient(client *ClientHandeller) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ClientsList[client] = true
}

func (m *Manager) RemoveClient(client *ClientHandeller) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ClientsList, client)
	client.conn.Close()
}
