package algorithms

import (
	customerrors "github.com/1ef7yy/go-loadbalancer/internal/custom_errors"
	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

type LeastConnections struct{}

func NewLeastConnections() *LeastConnections {
	return &LeastConnections{}
}

func (lc *LeastConnections) NextServer(pool []server.Server) (server.Server, error) {
	if len(pool) == 0 {
		return nil, customerrors.ErrNoServersAvailable
	}

	minConnections := -1
	var selected server.Server

	for _, server := range pool {
		if !server.IsAlive() {
			continue
		}

		conns := server.GetConnections()

		if minConnections == -1 || conns < minConnections {
			minConnections = conns
			selected = server
		}
	}

	if selected == nil {
		return nil, customerrors.ErrNoHealthyServersAvailable
	}

	return selected, nil
}

func (lc *LeastConnections) Name() string {
	return "least-connections"
}
