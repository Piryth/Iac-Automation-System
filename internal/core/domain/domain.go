package domain

type Workflow struct {
	ID          string
	Name        string
	Description string
	Trigger     Trigger
	Graph       Graph
}

type Trigger struct {
	Type   string
	Config map[string]string
}

type Graph struct {
	Entrypoint string
	Nodes      map[string]Node
}

type Node struct {
	ID        string
	Type      string
	OnSuccess string            // Empty if it's the last node in the graph
	Config    map[string]string // Holds type-specific configuration
}
