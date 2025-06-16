package algorithms

import (
	customerrors "github.com/1ef7yy/go-loadbalancer/internal/custom_errors"
	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

type RoundRobin struct {
	current int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		current: -1,
	}
}

func (rr *RoundRobin) NextServer(pool []server.Server) (server.Server, error) {
	if len(pool) == 0 {
		return nil, customerrors.ErrNoServersAvailable
	}

	rr.current = (rr.current + 1) % len(pool)
	return pool[rr.current], nil
}

func (rr *RoundRobin) Name() string {
	return "round-robin"
}
