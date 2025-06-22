package algorithms

import (
	"errors"
	"sync"

	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

type WeightedRoundRobin struct {
	currentIndex int
	mu           sync.Mutex
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		currentIndex: -1,
	}
}

func (wrr *WeightedRoundRobin) NextServer(servers []server.Server) (server.Server, error) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	if len(servers) == 0 {
		return nil, errors.New("no servers available")
	}

	aliveServers := make([]server.Server, 0, len(servers))
	totalWeight := 0
	for _, s := range servers {
		if s.IsAlive() {
			aliveServers = append(aliveServers, s)
			totalWeight += s.GetWeight()
		}
	}

	if len(aliveServers) == 0 {
		return nil, errors.New("no alive servers available")
	}

	var selected server.Server
	maxEffectiveWeight := -1.0

	wrr.currentIndex = (wrr.currentIndex + 1) % len(aliveServers)
	startIndex := wrr.currentIndex

	for i := 0; i < len(aliveServers); i++ {
		idx := (startIndex + i) % len(aliveServers)
		s := aliveServers[idx]
		effectiveWeight := s.GetEffectiveWeight()

		if effectiveWeight > maxEffectiveWeight {
			maxEffectiveWeight = effectiveWeight
			selected = s
			wrr.currentIndex = idx
		}
	}

	if selected != nil {
		return selected, nil
	}

	wrr.currentIndex = (wrr.currentIndex + 1) % len(aliveServers)
	return aliveServers[wrr.currentIndex], nil
}

func (wrr *WeightedRoundRobin) Name() string {
	return "Weighted Round Robin"
}
