package parser

import (
	"fmt"
	"path/filepath"
	pb "piryth/ias/api/proto"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// --- Go Structs representing the CUE Schema ---

type ProjectConfig struct {
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
	Workflows map[string]WorkflowDef `json:"workflows"`
}

type WorkflowDef struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Trigger     TriggerDef `json:"trigger"`
	Graph       GraphDef   `json:"graph"`
}

type TriggerDef struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type GraphDef struct {
	Entrypoint string             `json:"entrypoint"`
	Nodes      map[string]NodeDef `json:"nodes"`
}

type NodeDef struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	OnSuccess string            `json:"on_success,omitempty"`
	Config    map[string]string `json:"config"`
}

// --- Parsing Logic ---

// ParseProject loads the CUE package from the given directory, merges it with
// the embedded schema, and exports it as Protobuf messages.
func ParseProject(repoDir string) (*ProjectConfig, []*pb.Workflow, error) {
	ctx := cuecontext.New()
	absRepoDir, err := filepath.Abs(repoDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 1. Overlay the embedded schema into the virtual filesystem
	virtualSchemaPath := filepath.Join(absRepoDir, "_embedded_schema.cue")
	cfg := &load.Config{
		Dir: repoDir,
		Overlay: map[string]load.Source{
			// Note: EmbeddedSchema is defined in embed.go via //go:embed
			virtualSchemaPath: load.FromBytes(EmbeddedSchema),
		},
	}

	// 3. Load the CUE packages (returns one instance per directory)
	bis := load.Instances([]string{"./..."}, cfg)
	if len(bis) == 0 {
		return nil, nil, fmt.Errorf("no CUE packages found in %s", absRepoDir)
	}

	// 4. Build and merge ALL directories into a single CUE Value
	var finalVal cue.Value
	for i, b := range bis {
		if b.Err != nil {
			return nil, nil, fmt.Errorf("cue load error in %s: %w", b.Dir, b.Err)
		}

		v := ctx.BuildInstance(b)
		if v.Err() != nil {
			return nil, nil, fmt.Errorf("cue build error in %s: %w", b.Dir, v.Err())
		}

		if i == 0 {
			finalVal = v
		} else {
			finalVal = finalVal.Unify(v) // Merge this directory with the root
		}
	}

	// 5. Validate the combined result
	if err := finalVal.Validate(); err != nil {
		return nil, nil, fmt.Errorf("schema validation failed: %w", err)
	}

	// 6. Decode into our Go structs
	var config ProjectConfig
	if err := finalVal.Decode(&config); err != nil {
		return nil, nil, fmt.Errorf("failed to decode CUE to struct: %w", err)
	}

	// 7. Convert to Protobuf
	return &config, buildProtoWorkflows(&config), nil
}

// --- Conversion Helpers ---

// buildProtoWorkflows iterates over the CUE-decoded ProjectConfig
// and maps it into a slice of Protobuf Workflow messages.
func buildProtoWorkflows(config *ProjectConfig) []*pb.Workflow {
	var pbWorkflows []*pb.Workflow

	for wfID, w := range config.Workflows {
		pbWorkflows = append(pbWorkflows, &pb.Workflow{
			Id:          wfID,
			Name:        w.Name,
			Description: w.Description,
			Trigger: &pb.Trigger{
				Type:   w.Trigger.Type,
				Config: w.Trigger.Config,
			},
			Graph: buildGraphProto(&w.Graph),
		})
	}

	return pbWorkflows
}

// buildGraphProto converts the GraphDef struct into the Protobuf Graph message.
func buildGraphProto(g *GraphDef) *pb.Graph {
	if g == nil {
		return nil
	}

	var pbNodes []*pb.Node
	for nodeID, n := range g.Nodes {
		pbNodes = append(pbNodes, &pb.Node{
			// Use the map key (nodeID) as the ID, matching the CUE schema pattern
			Id:        nodeID,
			Type:      n.Type,
			OnSuccess: n.OnSuccess,
			Config:    n.Config,
		})
	}

	return &pb.Graph{
		Entrypoint: g.Entrypoint,
		Nodes:      pbNodes,
	}
}
