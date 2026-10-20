package services

import (
	"context"
	"log"
	"piryth/ias/internal/core/domain"
	"piryth/ias/internal/core/ports"
)

type workflowService struct {
	repo     ports.WorkflowRepository
	runner   ports.CommandRunner
	registry *NodeRegistry // Should be a pointer ?
}

func NewWorkflowService(repo ports.WorkflowRepository, registry NodeRegistry) ports.WorkflowService {
	return &workflowService{
		repo:     repo,
		registry: &registry,
	}
}

func (s *workflowService) Apply(ctx context.Context, wf domain.Workflow) error {
	// Business logic validation happens here (e.g., checking HDAG cycles)

	// Save via the port, knowing nothing about PostgreSQL
	return s.repo.Save(ctx, wf)
}

func (s *workflowService) Execute(ctx context.Context, id string) error {
	// Fetching workflow
	wf, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Getting entry node
	entryNode := wf.Graph.Nodes[wf.Graph.Entrypoint]

	// Finding the entry node executor
	executor, err := s.registry.Get(entryNode.Type)
	if err != nil {
		return err
	}

	// Node execution
	execCtx := ports.ExecutionContext{WorkflowID: wf.ID, NodeID: entryNode.ID}
	output, err := executor.Execute(ctx, execCtx, entryNode.Config)
	if err != nil {
		log.Fatal("Error executing workflow: ", err)
		return err
	}

	log.Println("Workflow %s success\nOutput: %s\n", wf.ID, output)
	return nil
}
