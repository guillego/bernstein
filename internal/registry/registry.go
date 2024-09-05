package registry

// registry implements a Registry for the nodes and containers
// this is a very basic implementation and will probably need to be broken up

import (
	"fmt"
	"log"
	"sync"
)

// NodeData contains the information relevant to a node
// cpu and ram indicate the *available* capacity in the node
type NodeData struct {
	name       string
	ip         string
	status     string
	cpu        uint16
	ram        uint16
	containers []string
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

var ErrNodeNotFound = fmt.Errorf("node not found")

// GetNode retrieves the node data from the Registry
func (r *Registry) GetNode(name string) (NodeData, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[name]
	if exists {
		log.Printf("GetNode: name=%s", name)
	} else {
		log.Printf("GetNode: name=%s not found", name)
	}
	return node, exists
}

func (r *Registry) AddNode(name string, ip string, status string, cpu uint16, ram uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[name] = NodeData{name: name, ip: ip, status: status, cpu: cpu, ram: ram, containers: make([]string, 0)}
	log.Printf("AddNode: name=%s - ip=%s, status=%s, cpu=%d, ram=%d", name, ip, status, cpu, ram)
}

// HACK! This is an extremely naive implementation, bad for race conditions
func (r *Registry) AddContainerToNode(name string, container string, cpu_req uint16, ram_req uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[name]
	if !exists {
		return ErrNodeNotFound
	}

	node.containers = append(node.containers, container)
	node.cpu = node.cpu - cpu_req
	node.ram = node.ram - ram_req

	log.Printf("AddContainerToNode: name=%s - ip=%s, status=%s, cpu=%d, ram=%d", node.name, node.ip, node.status, node.cpu, node.ram)

	return nil
}

func (r *Registry) UpdateNodeStatus(name string, ip string, status string, cpu uint16, ram uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes[name] = NodeData{name: name, ip: ip, status: status, cpu: cpu, ram: ram}
	log.Printf("UpdateNode: name=%s - ip=%s, status=%s, cpu=%d, ram=%d", name, ip, status, cpu, ram)

}

func (r *Registry) DeleteNode(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.nodes, name)

	log.Printf("DeleteNode: name=%s", name)
}
