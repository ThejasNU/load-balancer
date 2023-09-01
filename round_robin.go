package roundRobin

import "sync"

type LoadBalancer struct {
	Port    string
	Servers []Server
	mu      sync.Mutex
	index   int
}

type Server struct {
	Port     string
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

func (lb *LoadBalancer) GetNextServer() *Server{
	lb.mu.Lock()
	defer lb.mu.Unlock()
	
	numServers:= len(lb.Servers)
	nextIndex:=(lb.index+1)%numServers

	//find next server which is alive, if all servers are dead return nil
	var count int=0
	for lb.Servers[nextIndex].GetAliveStatus()==false && count<numServers{
		nextIndex=(nextIndex+1)%numServers
		count++
	}
	
	if count==numServers{
		return nil
	} else {
		return &lb.Servers[nextIndex]
	}
}
