package routes

import "sync"

type Manager struct {
	clientsList ClientList
	mu          sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clientsList: make(ClientList),
		mu:          sync.RWMutex{},
	}
}

func (m *Manager) AddClient(client *ClientHandeller) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clientsList[client] = true
}

func (m *Manager) RemoveClient(client *ClientHandeller) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clientsList, client)
	client.conn.Close()
}
