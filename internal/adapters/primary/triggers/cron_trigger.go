package triggers

import (
	"context"
	"fmt"
	_ "piryth/ias/internal/core/ports"
	"sync"

	"github.com/robfig/cron/v3"
)

type CronTrigger struct {
	cronEngine *cron.Cron
	jobIDs     map[string]cron.EntryID
	mu         sync.Mutex // Protects map from concurrent Start/Stop panics
}

func NewCronTrigger() *CronTrigger {
	c := cron.New()
	c.Start() // Start the engine immediately
	return &CronTrigger{
		cronEngine: c,
		jobIDs:     make(map[string]cron.EntryID),
	}
}

func (t *CronTrigger) Start(ctx context.Context, workflowID string, config map[string]string, executeFn func(id string)) error {
	schedule := config["schedule"]

	// Schedule the job and pass the callback function
	entryID, err := t.cronEngine.AddFunc(schedule, func() {
		executeFn(workflowID)
	})

	if err != nil {
		return fmt.Errorf("invalid cron schedule %s: %w", schedule, err)
	}

	t.jobIDs[workflowID] = entryID
	fmt.Printf("Started Cron trigger for workflow %s on schedule: %s\n", workflowID, schedule)
	return nil
}

func (t *CronTrigger) Stop(ctx context.Context, workflowID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 1. Check if the workflow actually has a running cron job
	if entryID, exists := t.jobIDs[workflowID]; exists {
		// 2. Tell the cron engine to unschedule it
		t.cronEngine.Remove(entryID)

		// 3. Remove it from our memory
		delete(t.jobIDs, workflowID)
		fmt.Printf("Successfully stopped cron job for workflow %s\n", workflowID)
	}

	return nil
}
