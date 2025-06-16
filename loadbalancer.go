package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/1ef7yy/go-loadbalancer/internal/algorithms"
	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

type loadBalancer struct {
	serverPool []server.Server
	algo       algorithms.Algorithm
	mu         sync.RWMutex
	proxy      *httputil.ReverseProxy
}

type LoadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	AddServer(server.Server)
	RemoveServer(url string)
	SetAlgorithm(algorithms.Algorithm)
	GetAlgorithm() string
	HealthCheck()
}

func NewLoadBalancer(algo algorithms.Algorithm) LoadBalancer {
	return &loadBalancer{
		algo: algo,
		serverPool: make([]server.Server, 0),
		proxy: &httputil.ReverseProxy{
			Director: func(r *http.Request) {},
		},
	}
}

func (lb *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server, err := lb.algo.NextServer(lb.serverPool)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if server == nil {
		http.Error(w, "no available servers", http.StatusServiceUnavailable)
	}

	serverURL, _ := url.Parse(server.GetURL())

	lb.proxy.Director = func(r *http.Request) {
		r.URL.Scheme = serverURL.Scheme
		r.URL.Host = serverURL.Host
		r.Host = serverURL.Host
	}

	server.IncrementConnections()
	defer server.DecrementConnections()

	lb.proxy.ServeHTTP(w, r)
}

func (lb *loadBalancer) AddServer(server server.Server) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.serverPool = append(lb.serverPool, server)
}

func (lb *loadBalancer) RemoveServer(url string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for idx, server := range lb.serverPool {
		if server.GetURL() == url {
			lb.serverPool = append(lb.serverPool[:idx], lb.serverPool[idx+1:]...)
			break
		}
	}
}

func (lb *loadBalancer) SetAlgorithm(algo algorithms.Algorithm) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.algo = algo
}

func (lb *loadBalancer) GetAlgorithm() string {
	return lb.algo.Name()
}

func (lb *loadBalancer) HealthCheck() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, server := range lb.serverPool {
		alive := serverHealthCheck(server.GetURL())
		server.SetAlive(alive)
	}
}

func serverHealthCheck(url string) bool {
	// TODO: implement healthchecks
	return true
}


