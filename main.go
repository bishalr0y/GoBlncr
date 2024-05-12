package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	address string
	proxy   *httputil.ReverseProxy
}

type LoadBalancer struct {
	servers         []Server
	roundRobinCount int
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
		server = lb.servers[lb.roundRobinCount]
	}

	fmt.Println(server.address)
	return server.proxy
}

func (lb *LoadBalancer) serve() *httputil.ReverseProxy {
	server := lb.getNextServer()
	return server
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

	for _, server := range lb.servers {
		fmt.Println(server.address)
	}

	http.HandleFunc("/", reverseProxyHandler(lb.serve()))

	http.ListenAndServe(":8000", nil)
}

func reverseProxyHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the next server
		proxy.ServeHTTP(w, r)
	}
}
