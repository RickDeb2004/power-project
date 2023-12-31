package serve

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
	

)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	addr        string
	proxy       *httputil.ReverseProxy
	healthCheck *healthCheck
	Mutex       *sync.Mutex
	CurrentCons int
	weight   int
}

type healthCheck struct {
	url             string
	interval        time.Duration
	timeout         time.Duration
	healthy         bool
	lastCheckedTime time.Time
	Mutex           sync.Mutex
}

func NewServer(addr string, healthCheckURL string, healthCheckInterval, healthCheckTimeout time.Duration) *SimpleServer {
	serverURL, err := url.Parse(addr)
	if err!=nil{
  handleErr(err)
	}
	

	healthCheck := &healthCheck{
		url:             healthCheckURL,
		interval:        healthCheckInterval,
		timeout:         healthCheckTimeout,
		healthy:         true,
		lastCheckedTime: time.Now(),
	}

	proxy := httputil.NewSingleHostReverseProxy(serverURL)

	return &SimpleServer{
		addr:        addr,
		proxy:       proxy,
		healthCheck: healthCheck,
        Mutex:  &sync.Mutex{},
	}
}

func handleErr(err error) {
	panic(err)
}



func (s *SimpleServer) Address() string {
	return s.addr
}

func (s *SimpleServer) IsAlive() bool {
	s.healthCheck.Mutex.Lock()
	defer s.healthCheck.Mutex.Unlock()
	return s.healthCheck.healthy
}

func (s *SimpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

func (s *SimpleServer) StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(s.healthCheck.interval)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case <-ticker.C:
				s.checkHealth()
			}
		}
	}()
}

func (s *SimpleServer) checkHealth() {
	s.healthCheck.Mutex.Lock()
	defer s.healthCheck.Mutex.Unlock()

	if time.Since(s.healthCheck.lastCheckedTime) < s.healthCheck.interval {
		return
	}

	client := http.Client{Timeout: s.healthCheck.timeout}
	resp, err := client.Get(s.healthCheck.url)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.healthCheck.healthy = false
	} else {
		s.healthCheck.healthy = true
	}

	s.healthCheck.lastCheckedTime = time.Now();
}
