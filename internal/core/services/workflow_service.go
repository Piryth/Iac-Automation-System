package services

import (
	"context"
	"fmt"
	"log"
	"piryth/ias/internal/core/domain"
	"piryth/ias/internal/core/ports"
)

type workflowService struct {
	repo            ports.WorkflowRepository
	registry        *NodeRegistry
	triggerRegistry *TriggerRegistry
}

func NewWorkflowService(repo ports.WorkflowRepository, registry NodeRegistry, triggerRegistry TriggerRegistry) ports.WorkflowService {
	return &workflowService{
		repo:            repo,
		registry:        &registry,
		triggerRegistry: &triggerRegistry,
	}
}

func (s *workflowService) Apply(ctx context.Context, wf domain.Workflow) error {
	// Workflow validation here

	// 1. Find workflow names/ids and compare to existing db
	wfs, err := s.repo.GetAll(ctx)
	names := make([]string, len(wfs))
	if err != nil {
		return err
	}
	// Get list of Names
	for i, _ := range wfs {
		names[i] = wfs[i].Name

		// RULE 1 : Workflow in DB MUST be the same than
		if wfs[i].Name != wf.Name {
			// delete
			err := s.repo.DeleteByID(ctx, names[i])
			if err != nil {
				return err
			}
		}
	}

	executeCallback := func(wfID string) {
		// Call the Execute phase in a background context
		_ = s.Execute(context.Background(), wfID)
	}

	// Save workflow to database
	if err := s.repo.Save(ctx, wf); err != nil {
		return err
	}

	return s.triggerRegistry.StartTrigger(
		ctx,
		wf.Trigger.Type,
		wf.ID,
		wf.Trigger.Config,
		executeCallback,
	)

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

func (s *workflowService) Delete(ctx context.Context, workflowID string) error {
	// 1. Fetch the workflow to know what trigger it uses
	wf, err := s.repo.GetByID(ctx, workflowID)
	if err != nil {
		return err // Or return nil if it already doesn't exist
	}

	// 2. Stop the active trigger in memory
	if err := s.triggerRegistry.StopTrigger(ctx, wf.Trigger.Type, wf.ID); err != nil {
		return fmt.Errorf("failed to stop trigger: %w", err)
	}

	// 3. Delete from the database
	return s.repo.DeleteByID(ctx, workflowID)
}

func (s *workflowService) List(ctx context.Context) ([]domain.Workflow, error) {
	return s.repo.GetAll(ctx)
}
