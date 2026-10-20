package ports

import (
	"context"
)

// ExecutionContext contains the runtime data (passed between nodes later)
type ExecutionContext struct {
	WorkflowID string
	NodeID     string
	// Future: Data map for passing inputs/outputs between nodes
}

// NodeExecutor is the interface all custom nodes must satisfy
type NodeExecutor interface {
	// Execute takes the node's specific config map and performs the action
	Execute(ctx context.Context, execCtx ExecutionContext, config map[string]string) (string, error)
}
