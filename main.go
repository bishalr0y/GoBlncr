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

	// TODO pass the server proxy
	fmt.Println(lb.servers[1].isAlive())

	http.HandleFunc("/", reverseProxyHandler(lb.servers[1].proxy))

	http.ListenAndServe(":8000", nil)
}

func reverseProxyHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
