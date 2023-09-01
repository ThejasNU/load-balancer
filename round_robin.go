package roundRobin

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
	index   int
}

type Server struct {
	Url    string
	IsAlive bool
	mu      sync.RWMutex
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

// ############ LOAD BALANCER MAIN FUNCTION #######################################################
func (lb *LoadBalancer) GetNextServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	numServers := len(lb.Servers)
	nextIndex := (lb.index + 1) % numServers

	//find next server which is alive, if all servers are dead return nil
	var count int = 0
	for lb.Servers[nextIndex].GetAliveStatus() == false && count < numServers {
		nextIndex = (nextIndex + 1) % numServers
		count++
	}

	if count == numServers {
		return nil
	} else {
		return &lb.Servers[nextIndex]
	}
}

//############ HEALTH CHECKER #######################################################
func (lb *LoadBalancer)checkHealth(){
	t:=time.NewTicker(time.Minute*1)
	
	for{
		select{
		case <-t.C:
			for i:=0;i<len(lb.Servers);i++{
				var server *Server=&lb.Servers[i]
				//make sure that server details are not being updated
				server.mu.RLock()
				pingUrl,err:=url.Parse(server.Url)
				if err!=nil{
					log.Fatal(err.Error())
				}
				server.mu.RUnlock()
				//call helper which pings the url
				isAlive:=checkStatus(pingUrl)
				//setAliveStatus has lock implementation inside it so need of lock here
				server.SetAliveStatus(isAlive)
			}
		}
	}
}

func checkStatus(url *url.URL) bool {
	conn,err:=net.DialTimeout("tcp",url.Host,time.Second*10)
	defer conn.Close()
	
	if err!=nil{
		return false
		} else {
			return true
		}
	}
	
//############ MAIN FUNCTION #######################################################
func main(){
	//create dummy servers and add
	var lb LoadBalancer=LoadBalancer{
		Port: "8080",
		Servers: []Server{},
	}

	//execute health checker on different go routine
	go lb.checkHealth()

	for i:=0;i<5;i++{
		server:=lb.GetNextServer()
		println(server.Url)
	}
}