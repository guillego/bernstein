package registry

import (
	"errors"
	"testing"
)

func TestRegistry_AddNode(t *testing.T) {
	r := NewRegistry()

	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	node, err := r.GetNode("foo")
	if err != nil {
		t.Fatalf("expected node to exist, but got error: %v", err)
	}
	if node.IP != "172.26.0.1" {
		t.Errorf("expected node IP '172.26.0.1', got '%s'", node.IP)
	}
	if node.Status != "ready" {
		t.Errorf("expected node status 'ready', got '%s'", node.Status)
	}
	if node.CPU != 4 || node.RAM != 512 {
		t.Errorf("expected node resources to be CPU=4, RAM=512; got CPU=%d, RAM=%d", node.CPU, node.RAM)
	}
}

func TestRegistry_AddNode_AlreadyExists(t *testing.T) {
	r := NewRegistry()
	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	err = r.AddNode("foo", "172.26.0.2", "active", 8, 1024)
	if !errors.Is(err, ErrNodeAlreadyExists) {
		t.Fatalf("expected error '%v', got '%v'", ErrNodeAlreadyExists, err)
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

	node, err := r.GetNode("foo")
	if err != nil {
		t.Fatalf("expected node to exist after adding container, but got error: %v", err)
	}
	if node.CPU != 3 {
		t.Errorf("expected remaining CPU to be 3, got %d", node.CPU)
	}
	if node.RAM != 384 {
		t.Errorf("expected remaining RAM to be 384, got %d", node.RAM)
	}
	if len(node.Containers) != 1 {
		t.Errorf("expected 1 container, got %d", len(node.Containers))
	}
	if node.Containers[0] != "bar" {
		t.Errorf("expected container 'bar', got '%s'", node.Containers[0])
	}
}

func TestRegistry_AddContainerToNode_NodeNotFound(t *testing.T) {
	r := NewRegistry()

	err := r.AddContainerToNode("foo", "bar", 1, 128)
	if !errors.Is(err, ErrNodeNotFound) {
		t.Fatalf("expected error '%v', got '%v'", ErrNodeNotFound, err)
	}
}

func TestRegistry_AddContainerToNode_InsufficientResources(t *testing.T) {
	r := NewRegistry()
	err := r.AddNode("foo", "172.26.0.1", "ready", 2, 128)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	err = r.AddContainerToNode("foo", "bar", 4, 256)
	if err == nil {
		t.Fatalf("expected error due to insufficient resources, but got none")
	}
}

func TestRegistry_UpdateNodeStatus(t *testing.T) {
	r := NewRegistry()
	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	err = r.UpdateNodeStatus("foo", "172.26.0.2", "active", 8, 1024)
	if err != nil {
		t.Fatalf("unexpected error updating node: %v", err)
	}

	node, err := r.GetNode("foo")
	if err != nil {
		t.Fatalf("expected node to exist after update, but got error: %v", err)
	}
	if node.IP != "172.26.0.2" {
		t.Errorf("expected updated IP '172.26.0.2', got '%s'", node.IP)
	}
	if node.Status != "active" {
		t.Errorf("expected updated status 'active', got '%s'", node.Status)
	}
	if node.CPU != 8 || node.RAM != 1024 {
		t.Errorf("expected updated resources CPU=8, RAM=1024; got CPU=%d, RAM=%d", node.CPU, node.RAM)
	}
}

func TestRegistry_DeleteNode(t *testing.T) {
	r := NewRegistry()
	err := r.AddNode("foo", "172.26.0.1", "ready", 4, 512)
	if err != nil {
		t.Fatalf("unexpected error adding node: %v", err)
	}

	err = r.DeleteNode("foo")
	if err != nil {
		t.Fatalf("unexpected error deleting node: %v", err)
	}

	_, err = r.GetNode("foo")
	if err != ErrNodeNotFound {
		t.Fatalf("expected error '%v', got '%v'", ErrNodeNotFound, err)
	}
}
