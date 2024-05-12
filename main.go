package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	address string
	proxy   *httputil.ReverseProxy
}

type LoadBalancer struct {
	servers         []Server
	roundRobinCount int
	mutex           sync.Mutex
}

func (s *Server) isAlive() bool {
	_, err := http.Get(s.address)
	if err != nil {
		fmt.Printf("Invalid address: %v\n", err)
		return false
	}
	return true
}

func (lb *LoadBalancer) getNextServer() *httputil.ReverseProxy {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]

	for !server.isAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	// fmt.Printf("Next Server: %v\n", server.address)
	return server.proxy
}

func (lb *LoadBalancer) serve(w http.ResponseWriter, r *http.Request) {

	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	server := lb.getNextServer()
	lb.roundRobinCount = lb.roundRobinCount % len(lb.servers)
	fmt.Printf("Current server: %v\n", lb.servers[lb.roundRobinCount].address)
	server.ServeHTTP(w, r)
}

func createServer(address string) Server {
	targetUrl, err := url.Parse(address)

	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)

	return Server{
		address: address,
		proxy:   reverseProxy,
	}
}

func createLoadBalancer(servers []Server) LoadBalancer {
	return LoadBalancer{
		servers:         servers,
		roundRobinCount: 0,
	}
}

func main() {
	fmt.Println("Hello world")

	servers := []Server{
		createServer("http://localhost:3001"),
		createServer("http://localhost:3002"),
		createServer("http://localhost:3003"),
	}

	lb := createLoadBalancer(servers)

	http.HandleFunc("/", lb.serve)

	http.ListenAndServe(":8000", nil)
}
