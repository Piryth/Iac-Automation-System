-- name: UpsertWorkflow :exec
INSERT INTO workflows (id, name, description, trigger, graph, updated_at)
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO UPDATE SET
                               name = EXCLUDED.name,
                               description = EXCLUDED.description,
                               trigger = EXCLUDED.trigger,
                               graph = EXCLUDED.graph,
                               updated_at = CURRENT_TIMESTAMP;

-- name: GetWorkflowByID :one
SELECT id, name, description, trigger, graph FROM workflows WHERE id = $1;

-- name: GetAllWorkflows :many
SELECT id, name, description, trigger, graph FROM workflows;