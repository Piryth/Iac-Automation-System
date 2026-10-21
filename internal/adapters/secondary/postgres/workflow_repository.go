package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"piryth/ias/internal/adapters/secondary/postgres/db"
	"piryth/ias/internal/core/domain"
	"piryth/ias/internal/core/ports"

	_ "github.com/lib/pq" // Postgres driver
)

type pgWorkflowRepository struct {
	queries *db.Queries
}

func NewPostgresWorkflowRepository(sqlDB *sql.DB) ports.WorkflowRepository {
	return &pgWorkflowRepository{
		queries: db.New(sqlDB),
	}
}

func (r *pgWorkflowRepository) Save(ctx context.Context, wf domain.Workflow) error {
	triggerBytes, err := json.Marshal(wf.Trigger)
	if err != nil {
		return fmt.Errorf("failed to marshal trigger: %w", err)
	}

	graphBytes, err := json.Marshal(wf.Graph)
	if err != nil {
		return fmt.Errorf("failed to marshal graph: %w", err)
	}

	// Use sqlc generated method
	err = r.queries.UpsertWorkflow(ctx, db.UpsertWorkflowParams{
		ID:          wf.ID,
		Name:        wf.Name,
		Description: wf.Description,
		Trigger:     triggerBytes, // sqlc jsonb is json.RawMessage
		Graph:       graphBytes,
	})

	return err
}

func (r *pgWorkflowRepository) GetByID(ctx context.Context, id string) (domain.Workflow, error) {
	// Use sqlc generated method
	row, err := r.queries.GetWorkflowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Workflow{}, fmt.Errorf("workflow %s not found", id)
		}
		return domain.Workflow{}, err
	}

	return r.mapToDomain(row.ID, row.Name, row.Description, row.Trigger, row.Graph)
}

func (r *pgWorkflowRepository) GetAll(ctx context.Context) ([]domain.Workflow, error) {
	// Use sqlc generated method
	rows, err := r.queries.GetAllWorkflows(ctx)
	if err != nil {
		return nil, err
	}

	var workflows []domain.Workflow
	for _, row := range rows {
		wf, err := r.mapToDomain(row.ID, row.Name, row.Description, row.Trigger, row.Graph)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

func (r *pgWorkflowRepository) DeleteByID(ctx context.Context, id string) error {
	err := r.queries.DeleteWorkflowByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// mapToDomain is a helper to convert sqlc DB rows back into clean Domain entities
func (r *pgWorkflowRepository) mapToDomain(id, name, desc string, triggerBytes, graphBytes json.RawMessage) (domain.Workflow, error) {
	var trigger domain.Trigger
	var graph domain.Graph

	if err := json.Unmarshal(triggerBytes, &trigger); err != nil {
		return domain.Workflow{}, fmt.Errorf("failed to unmarshal trigger for wf %s: %w", id, err)
	}
	if err := json.Unmarshal(graphBytes, &graph); err != nil {
		return domain.Workflow{}, fmt.Errorf("failed to unmarshal graph for wf %s: %w", id, err)
	}

	return domain.Workflow{
		ID:          id,
		Name:        name,
		Description: desc,
		Trigger:     trigger,
		Graph:       graph,
	}, nil
}
