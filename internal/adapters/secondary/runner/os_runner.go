package runner

import (
	"context"
	"os/exec"
	"piryth/ias/internal/core/ports"
)

type osRunner struct{}

func NewOSRunner() ports.CommandRunner {
	return &osRunner{}
}

func (r *osRunner) Run(ctx context.Context, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
