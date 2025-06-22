package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/1ef7yy/go-loadbalancer"
	"github.com/1ef7yy/go-loadbalancer/internal/algorithms"
	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

func main() {
	var algo algorithms.Algorithm

	lb_algo := os.Getenv("LB_ALGORITHM")
	switch lb_algo {
	case "round-robin":
		algo = algorithms.NewRoundRobin()
	case "random":
		algo = algorithms.NewRandom()
	case "least-connections":
		algo = algorithms.NewLeastConnections()
	default:
		fmt.Printf("LB_ALGORITHM of value %s does not match any type of LB algorithm", lb_algo)
		os.Exit(1)
	}
	lb := loadbalancer.NewLoadBalancer(algo)

	backendURLs := strings.Split(os.Getenv("BACKEND_URLS"), ",")
	for i, url := range backendURLs {
		healthFunc := func(url string) (ok bool) {
			resp, err := http.Get(fmt.Sprintf("%s/health", url))
			if err != nil {
				return false
			}
			return resp.StatusCode == http.StatusOK
		}
		server := server.NewServer(
			url,
			server.WithHealthFunc(healthFunc),
			server.WithInitialAliveStatus(true),
			server.WithWeight(i))
		lb.AddServer(server)

	}

	// background health checks
	go func() {
		for {
			lb.HealthCheck()
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the load balancer
	fmt.Printf("load balancer running on :8080 with algorithm %s", lb.GetAlgorithm())
	if err := http.ListenAndServe(":8080", lb); err != nil {
		fmt.Printf("error starting server: %v\n", err)
		os.Exit(1)
	}
}
