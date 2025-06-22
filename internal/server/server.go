package server

import (
	"sync"
)

type server struct {
	URL         string
	Alive       bool
	connections int
	weight      int
	HealthFunc  func(url string) bool
	HealthURL   string
	mu          sync.RWMutex
}

type Server interface {
	GetURL() string
	HealthCheck() bool
	SetAlive(bool)
	IsAlive() bool
	GetConnections() int
	IncrementConnections()
	DecrementConnections()
	GetWeight() int
	SetWeight(weight int)
	GetEffectiveWeight() float64
}

type Option func(*server)

func WithHealthFunc(f func(string) bool) Option {
	return func(s *server) {
		s.HealthFunc = f
	}
}

func WithHealthURL(url string) Option {
	return func(s *server) {
		s.HealthURL = url
	}
}

func WithInitialAliveStatus(alive bool) Option {
	return func(s *server) {
		s.Alive = alive
	}
}

func WithWeight(weight int) Option {
	return func(s *server) {
		s.weight = weight
	}
}

func NewServer(url string, opts ...Option) Server {
	s := &server{
		URL:        url,
		Alive:      true, // default value
		HealthFunc: nil,  // default nil, caller must set this
		weight:     1,    // default weight of 1
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *server) GetURL() string {
	return s.URL
}

func (s *server) HealthCheck() bool {
	if s.HealthFunc == nil {
		return false
	}
	return s.HealthFunc(s.URL)
}

func (s *server) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Alive
}

func (s *server) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Alive = alive
}

func (s *server) GetConnections() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connections
}

func (s *server) IncrementConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections++
}

func (s *server) DecrementConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections--
}

func (s *server) GetWeight() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.weight
}

func (s *server) SetWeight(weight int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weight = weight
}

func (s *server) GetEffectiveWeight() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.connections == 0 {
		return float64(s.weight)
	}
	return float64(s.weight) / float64(s.connections+1)
}
