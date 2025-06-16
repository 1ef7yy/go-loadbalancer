package algorithms

import "github.com/1ef7yy/go-loadbalancer/internal/server"

type Algorithm interface {
	NextServer([]server.Server) (server.Server, error)
	Name() string
}
