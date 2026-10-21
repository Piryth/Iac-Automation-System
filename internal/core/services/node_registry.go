package services

import (
	"fmt"
	"piryth/ias/internal/core/ports"
)

type NodeRegistry struct {
	executors map[string]ports.NodeExecutor
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		executors: make(map[string]ports.NodeExecutor),
	}
}

// Register adds a new custom node type to the platform
func (r *NodeRegistry) Register(nodeType string, executor ports.NodeExecutor) {
	r.executors[nodeType] = executor
}

// Get retrieves the executor for a specific node type
func (r *NodeRegistry) Get(nodeType string) (ports.NodeExecutor, error) {
	executor, exists := r.executors[nodeType]
	if !exists {
		return nil, fmt.Errorf("unsupported node type: %s", nodeType)
	}
	return executor, nil
}
