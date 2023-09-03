package leastconnections

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

// ############ STRUCTS #######################################################
type LoadBalancer struct {
	Port    string
	Servers []Server
	mu      sync.Mutex
}

type Server struct {
	Url     string
	IsAlive bool
	mu      sync.RWMutex
	numConnections int

}

// ############ HELPERS #######################################################
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

func (server *Server) GetNumConnections() int {
	server.mu.RLock()
	defer server.mu.RUnlock()

	return server.numConnections
}

func (server *Server) SetNumConnections(num int){
	server.mu.Lock()
	defer server.mu.Unlock()

	server.numConnections=num
}

// ############ LOAD BALANCER MAIN FUNCTION #######################################################
func (lb *LoadBalancer) GetNextServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	numServers := len(lb.Servers)

	var serverToBeReturned *Server=nil
	var numConnections=100000000000
	for i:=0;i<numServers;i++{
		var curServer *Server=&lb.Servers[i]

		if !curServer.GetAliveStatus(){
			continue
		}

		curServerNumConnections:=curServer.numConnections

		if curServerNumConnections<numConnections{
			serverToBeReturned=curServer
			numConnections=curServerNumConnections
		}
	}

	if serverToBeReturned!=nil{
		serverToBeReturned.SetNumConnections(numConnections+1)
	}

	return serverToBeReturned
}

// ############ HEALTH CHECKER #######################################################
func (lb *LoadBalancer) checkHealth() {
	t := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-t.C:
			for i := 0; i < len(lb.Servers); i++ {
				var server *Server = &lb.Servers[i]
				//make sure that server details are not being updated
				server.mu.RLock()
				pingUrl, err := url.Parse(server.Url)
				if err != nil {
					log.Fatal(err.Error())
				}
				server.mu.RUnlock()
				//call helper which pings the url
				isAlive := checkStatus(pingUrl)
				//setAliveStatus has lock implementation inside it so need of lock here
				server.SetAliveStatus(isAlive)
			}
		}
	}
}

func checkStatus(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Second*10)
	defer conn.Close()

	if err != nil {
		return false
	} else {
		return true
	}
}

// ############ MAIN FUNCTION #######################################################
func LeastconnectionsMain() {
	//create dummy servers and add
	var lb LoadBalancer = LoadBalancer{
		Port: "8080",
		Servers: []Server{
			{
				Url:     "htttp://localhost:8081/",
				IsAlive: true,
				numConnections: 0,
			},
			{
				Url:     "htttp://localhost:8082/",
				IsAlive: true,
				numConnections: 2,
			},
			{
				Url:     "htttp://localhost:8083/",
				IsAlive: true,
				numConnections: 1,
			},
			{
				Url:     "http://localhost:8084/",
				IsAlive: true,
				numConnections: 4,
			},
		},
	}

	//execute health checker on different go routine
	go lb.checkHealth()

	for i := 0; i < 10; i++ {
		server := lb.GetNextServer()
		println(server.Url)
	}
}