package load

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	cacher "power-project/cache"
	"power-project/responsewriter"
	serve "power-project/server"
	"sync"
)

type LoadBalancingAlgorithm int

const (
	RoundRobinCount LoadBalancingAlgorithm = iota
	WeightedRoundRobin
	LeastConnections
)

type LoadBalancer struct {
	port       string
	algorithm  LoadBalancingAlgorithm
	servers    []serve.Server
	mutex      *sync.Mutex
	connection int
	cache      *cacher.Cache
}

func NewLoadBalancer(port string, servers []serve.Server, algorithm LoadBalancingAlgorithm,  cache *cacher.Cache) *LoadBalancer {
	for _, server := range servers {
		server.(*serve.SimpleServer).StartHealthCheck()
	} 

	return &LoadBalancer{
		port:      port,
		algorithm: algorithm,
		servers:   servers,
		cache:     cache,
		mutex:     &sync.Mutex{},

		
	}
}
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error:%v \n", err)
		os.Exit(1)
	}
}

func (lb *LoadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	// Check cache first
	cacheKey := r.URL.String()
	//response, err := lb.cache.Get(cacheKey)
	// if err !=nil {
	// 	fmt.Println("Cache hit!")
	// 	w.Write(response)
	// 	return
	//}

	// if lb.cache.Exists(cacheKey) {
	// 	fmt.Println("Cache hit!")
    //    response:=lb.cache.Get()
	// 	w.Write(response)
	// 	return
	response,exists:=lb.cache.Get(cacheKey)
	if exists{
		fmt.Println("Cache hit!")
		w.Write(response)
		return
	}
	

	targetServer := lb.getAvailableServerFunc(r)
	fmt.Printf("Forwarding request to address %q\n", targetServer.Address())

	// Perform the request
	respWriter := responsewriter.NewResponseWriter(w)
	targetServer.Serve(respWriter, r)

	// Store response in cache
	if respWriter.GetStatus() == http.StatusOK {
		lb.cache.Set(cacheKey, respWriter.Body())
	}
}

func (lb *LoadBalancer) getAvailableServerFunc(r *http.Request) serve.Server {
	switch lb.algorithm {
	case RoundRobinCount:
		return lb.getAvailableServerRoundRobin()
	case WeightedRoundRobin:
		return lb.getAvailableServerWeightedRoundRobin()
	case LeastConnections:
		return lb.getAvailableServerLeastConnections()
	default:
		return lb.getAvailableServerRoundRobin()
	}
}

func (lb *LoadBalancer) getAvailableServerRoundRobin() serve.Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	server := lb.servers[lb.connection%len(lb.servers)]
	lb.connection++
	return server
}

func (lb *LoadBalancer) getAvailableServerWeightedRoundRobin() serve.Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	var totalWeight int
	for _, server := range lb.servers {
		simpleServer := server.(*serve.SimpleServer)
		if simpleServer.IsAlive() {
			totalWeight += 1 // Consider each server with a weight of 1
		}
	}
	var selectedServer serve.Server
	for _, server := range lb.servers {
		simpleServer := &serve.SimpleServer{}
		simpleServer.mutex.Lock()
		if simpleServer.IsAlive() {
			selectedServer = server
			simpleServer.mutex.Unlock()
			break
		}
		simpleServer.mutex.Unlock()
	}
	lb.connection++
	return selectedServer
}

func (lb *LoadBalancer) getAvailableServerLeastConnections() serve.Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	var minConnections int
	var selectedServer serve.Server
	for _, server := range lb.servers {
		simpleServer := &serve.SimpleServer{}
		simpleServer.mutex.Lock()
		if simpleServer.IsAlive() {
			if minConnections == 0 || simpleServer.currentcons < minConnections {
				minConnections = simpleServer.currentcons
				selectedServer = server
			}
		}
		simpleServer.mutex.Unlock()
	}
	lb.connection++
	return selectedServer
}

func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("failed to parse url %v\n", err)
		os.Exit(1)
	}
	return parsedURL
}
