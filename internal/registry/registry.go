package registry

import (
	"fmt"
	"log"
	"sync"
)

// NodeData contains the information relevant to a node
type NodeData struct {
	Name       string
	IP         string
	Status     string
	CPU        uint16
	RAM        uint16
	Containers []string
}

// Registry manages a collection of nodes
// Nodes are mutable
type Registry struct {
	mu    sync.RWMutex
	nodes map[string]*NodeData
}

// NewRegistry creates a new Registry instance
func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]*NodeData),
	}
}

var (
	ErrNodeNotFound      = fmt.Errorf("node not found")
	ErrNodeAlreadyExists = fmt.Errorf("node already exists")
)

// GetNode retrieves the node data from the Registry
func (r *Registry) GetNode(name string) (*NodeData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	node, exists := r.nodes[name]
	if !exists {
		log.Printf("GetNode: name=%s not found", name)
		return nil, ErrNodeNotFound
	}
	log.Printf("GetNode: name=%s", name)
	return node, nil
}

// AddNode adds a new node to the Registry
func (r *Registry) AddNode(name, ip, status string, cpu, ram uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[name]; exists {
		return fmt.Errorf("%w: %s", ErrNodeAlreadyExists, name)
	}

	r.nodes[name] = &NodeData{
		Name:       name,
		IP:         ip,
		Status:     status,
		CPU:        cpu,
		RAM:        ram,
		Containers: []string{},
	}
	log.Printf("AddNode: name=%s - ip=%s, status=%s, cpu=%d, ram=%d", name, ip, status, cpu, ram)
	return nil
}

// AddContainerToNode adds a container to a node and updates its resources
// HACK! This is an extremely naive implementation, to be updated
func (r *Registry) AddContainerToNode(name, container string, cpuReq, ramReq uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, name)
	}

	// Ensure the node has enough resources to allocate to the container
	if node.CPU < cpuReq || node.RAM < ramReq {
		return fmt.Errorf("insufficient resources on node %s", name)
	}

	node.Containers = append(node.Containers, container)
	node.CPU -= cpuReq
	node.RAM -= ramReq

	log.Printf("AddContainerToNode: name=%s - container=%s, cpuReq=%d, ramReq=%d", node.Name, container, cpuReq, ramReq)
	return nil
}

// UpdateNodeStatus updates the status and resources of a node
func (r *Registry) UpdateNodeStatus(name, ip, status string, cpu, ram uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, name)
	}

	node.IP = ip
	node.Status = status
	node.CPU = cpu
	node.RAM = ram
	log.Printf("UpdateNodeStatus: name=%s - ip=%s, status=%s, cpu=%d, ram=%d", name, ip, status, cpu, ram)
	return nil
}

// DeleteNode removes a node from the Registry
func (r *Registry) DeleteNode(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[name]; !exists {
		return fmt.Errorf("%w: %s", ErrNodeNotFound, name)
	}

	delete(r.nodes, name)
	log.Printf("DeleteNode: name=%s", name)
	return nil
}
