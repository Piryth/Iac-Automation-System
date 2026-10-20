package cron

import (
	"context"
	"fmt"
	"piryth/ias/internal/core/ports"

	"github.com/robfig/cron/v3"
)

type SchedulerAdapter struct {
	cron    *cron.Cron
	service ports.WorkflowService
	repo    ports.WorkflowRepository
}

func NewSchedulerAdapter(service ports.WorkflowService, repo ports.WorkflowRepository) *SchedulerAdapter {
	return &SchedulerAdapter{
		cron:    cron.New(),
		service: service,
		repo:    repo,
	}
}

func (s *SchedulerAdapter) Start(ctx context.Context) error {
	workflows, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, wf := range workflows {
		// In MVP, we assume wf.Trigger.Schedule holds the cron string
		if wf.Trigger.Type == "cron" {
			schedule := wf.Trigger.Config["schedule"]
			workflowID := wf.ID

			// Add job to cron engine
			_, err := s.cron.AddFunc(schedule, func() {
				s.service.Execute(context.Background(), workflowID)
			})

			if err != nil {
				fmt.Printf("Failed to schedule workflow %s: %v\n", workflowID, err)
			} else {
				fmt.Printf("Scheduled workflow %s with cron '%s'\n", workflowID, schedule)
			}
		}
	}

	s.cron.Start()
	return nil
}
