package ports

import "context"

type Trigger interface {
	// Start begins listening for the event.
	// executeCallback is the function the trigger calls when it fires.
	Start(ctx context.Context, workflowID string, config map[string]string, executeCallback func(workflowID string)) error
	Stop(ctx context.Context, workflowID string) error
}
