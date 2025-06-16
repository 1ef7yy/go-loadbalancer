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

	lb_algo, ok := os.LookupEnv("LB_ALGORITHM")
	if !ok {
		lb_algo = "round-robin"
	}
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

	// Get backend URLs from environment variable
	backendURLs := strings.Split(os.Getenv("BACKEND_URLS"), ",")
	for _, url := range backendURLs {
		server := server.NewServer(url)
		lb.AddServer(server)

	}

	// Start health checks in background
	go func() {
		for {
			lb.HealthCheck()
			time.Sleep(30 * time.Second)
		}
	}()

	// Start the load balancer
	fmt.Printf("load balancer running on :8080 with algorithm %s", lb.GetAlgorithm())
	if err := http.ListenAndServe(":8080", lb); err != nil {
		fmt.Printf("error starting server: %v\n", err)
		os.Exit(1)
	}
}
