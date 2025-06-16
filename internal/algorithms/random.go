// internal/algorithms/random.go
package algorithms

import (
	"math/rand"
	"sync"
	"time"

	customerrors "github.com/1ef7yy/go-loadbalancer/internal/custom_errors"
	"github.com/1ef7yy/go-loadbalancer/internal/server"
)

type Random struct {
	r  *rand.Rand
	mu sync.Mutex
}

func NewRandom() *Random {
	return &Random{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Random) NextServer(servers []server.Server) (server.Server, error) {
	if len(servers) == 0 {
		return nil, customerrors.ErrNoServersAvailable
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var aliveServers []server.Server
	for _, s := range servers {
		if s.IsAlive() {
			aliveServers = append(aliveServers, s)
		}
	}

	if len(aliveServers) == 0 {
		return nil, customerrors.ErrNoHealthyServersAvailable
	}

	return aliveServers[r.r.Intn(len(aliveServers))], nil
}

func (r *Random) Name() string {
	return "random"
}
