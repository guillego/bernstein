package registry

import (
	"log"
	"sync"
)

type NodeData struct {
	ip       string
	status   string
	capacity uint16
}

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]NodeData
}

func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]NodeData),
	}
}

func (r *Registry) Add(name string, ip string, status string, capacity uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[name] = NodeData{ip: ip, status: status, capacity: capacity}
	log.Printf("Registered new node: name=%s, ip=%s, status:%s, capacity:%d", name, ip, status, capacity)
}
