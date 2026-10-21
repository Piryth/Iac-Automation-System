package services

import (
	"context"
	"fmt"
	"piryth/ias/internal/core/ports"
)

type TriggerRegistry struct {
	triggers map[string]ports.Trigger
}

func NewTriggerRegistry() *TriggerRegistry {
	return &TriggerRegistry{
		triggers: make(map[string]ports.Trigger),
	}
}

func (r *TriggerRegistry) Register(triggerType string, trigger ports.Trigger) {
	r.triggers[triggerType] = trigger
}

func (r *TriggerRegistry) StartTrigger(ctx context.Context, triggerType string, workflowID string, config map[string]string, executeFn func(id string)) error {
	trigger, exists := r.triggers[triggerType]
	if !exists {
		return fmt.Errorf("unsupported trigger type: %s", triggerType)
	}
	return trigger.Start(ctx, workflowID, config, executeFn)
}

func (r *TriggerRegistry) StopTrigger(ctx context.Context, triggerType string, workflowID string) error {
	trigger, exists := r.triggers[triggerType]
	if !exists {
		return fmt.Errorf("unsupported trigger type: %s", triggerType)
	}

	// Delegate to the specific adapter (e.g., the CronTrigger)
	return trigger.Stop(ctx, workflowID)
}
