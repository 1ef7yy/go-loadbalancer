package server

import (
	"sync"
)

type server struct {
	URL         string
	Alive       bool
	connections int
	mu          sync.RWMutex
}

type Server interface {
	GetURL() string
	IsAlive() bool
	GetConnections() int
	IncrementConnections()
	DecrementConnections()
	SetAlive(bool)
}

func NewServer(url string) Server {
	return &server{
		URL:   url,
		Alive: true,
	}
}

func (s *server) GetURL() string {
	return s.URL
}
func (s *server) IsAlive() bool {
	return s.Alive
}
func (s *server) GetConnections() int {
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *server) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Alive = alive
}
