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
	mutex       *sync.Mutex
	currentcons int
	weight   int
}

type healthCheck struct {
	url             string
	interval        time.Duration
	timeout         time.Duration
	healthy         bool
	lastCheckedTime time.Time
	mutex           sync.Mutex
}

func newServer(addr string, healthCheckURL string, healthCheckInterval, healthCheckTimeout time.Duration) *SimpleServer {
	serverURL, err := url.Parse(addr)
	handleErr(err)

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
        mutex:&sync.Mutex{},
	}
}

func handleErr(err error) {
	panic("unimplemented")
}



func (s *SimpleServer) Address() string {
	return s.addr
}

func (s *SimpleServer) IsAlive() bool {
	s.healthCheck.mutex.Lock()
	defer s.healthCheck.mutex.Unlock()
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
	s.healthCheck.mutex.Lock()
	defer s.healthCheck.mutex.Unlock()

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
