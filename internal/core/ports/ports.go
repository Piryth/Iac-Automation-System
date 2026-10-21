package ports

import (
	"context"
	"piryth/ias/internal/core/domain"
)

type WorkflowService interface {
	Apply(ctx context.Context, wf domain.Workflow) error
	Execute(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Workflow, error)
	Delete(ctx context.Context, workflowID string) error
}

type WorkflowRepository interface {
	Save(ctx context.Context, wf domain.Workflow) error
	GetByID(ctx context.Context, id string) (domain.Workflow, error)
	GetAll(ctx context.Context) ([]domain.Workflow, error)
	DeleteByID(ctx context.Context, id string) error
}

type CommandRunner interface {
	Run(ctx context.Context, command string) (string, error)
}
