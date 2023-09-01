package roundRobin

import "sync"

type LoadBalancer struct {
	Port    string
	Servers []Server
	mu      sync.Mutex
	index   int
}

type Server struct {
	Url     string
	IsAlive bool
	mu      sync.RWMutex
}

func (server *Server) GetAliveStatus() bool {
	server.mu.RLock()
	defer server.mu.RUnlock()

	isAlive := server.IsAlive
	return isAlive
}

func (server *Server) SetAliveStatus(currentStatus bool) {
	server.mu.Lock()
	defer server.mu.Unlock()

	server.IsAlive = currentStatus
}
