package registry

import (
	"testing"
)

func TestRegistry_AddNode(t *testing.T) {
	r := NewRegistry()

	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	node, exists := r.GetNode("foo")
	if !exists {
		t.Fatalf("expected node to exist, but it does not")
	}
	if node.ip != "172.26.0.1" {
		t.Errorf("expected node IP '172.26.0.1', got '%s'", node.ip)
	}
}

func TestRegistry_AddContainerToNode(t *testing.T) {
	r := NewRegistry()
	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	err = r.AddContainerToNode("foo", "bar", 1, 128)
	if err != nil {
		t.Fatalf("unexpected error adding container to node: %v", err)
	}

	node, exists := r.GetNode("foo")
	if !exists {
		t.Fatalf("expected node to exist after adding container, but it does not")
	}
	if node.cpu != 3 {
		t.Errorf("expected remaining CPU to be 3, got %d", node.cpu)
	}
	if len(node.containers) != 1 {
		t.Errorf("expected 1 container, got %d", len(node.containers))
	}
}

