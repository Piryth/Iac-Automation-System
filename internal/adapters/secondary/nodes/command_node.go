package nodes

import (
	"context"
	"fmt"
	"piryth/ias/internal/core/ports"
)

type CommandNode struct {
	runner ports.CommandRunner
}

func NewCommandNode(runner ports.CommandRunner) *CommandNode {
	return &CommandNode{runner: runner}
}

func (n *CommandNode) Execute(ctx context.Context, execCtx ports.ExecutionContext, config map[string]string) (string, error) {
	cmd, exists := config["command"]
	if !exists {
		return "", fmt.Errorf("missing 'command' in config")
	}
	return n.runner.Run(ctx, cmd)
}
