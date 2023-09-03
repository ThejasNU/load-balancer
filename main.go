package main

import (
	leastconnections "github.com/ThejasNU/load-balancer/least-connections"
	roundrobin "github.com/ThejasNU/load-balancer/round-robin"
)

func main(){
	roundrobin.RoundrobinMain()
	leastconnections.LeastconnectionsMain()
}