package grpcadapter

import (
	pb "piryth/ias/api/proto"
	"piryth/ias/internal/core/domain"
)

// ToDomainWorkflow converts a Protobuf Workflow to a Core Domain Workflow
func ToDomainWorkflow(pbWf *pb.Workflow) domain.Workflow {
	if pbWf == nil {
		return domain.Workflow{}
	}

	// 1. Map Trigger
	var domainTrigger domain.Trigger
	if t := pbWf.GetTrigger(); t != nil {
		domainTrigger = domain.Trigger{
			Type:   t.GetType(),
			Config: t.GetConfig(),
		}
	}

	// 2. Map Graph and Nodes
	var domainGraph domain.Graph
	if g := pbWf.GetGraph(); g != nil {
		domainGraph = domain.Graph{
			Entrypoint: g.GetEntrypoint(),
			Nodes:      make(map[string]domain.Node),
		}
		for _, n := range g.GetNodes() {
			domainGraph.Nodes[n.GetId()] = domain.Node{
				ID:        n.GetId(),
				Type:      n.GetType(),
				OnSuccess: n.GetOnSuccess(),
				Config:    n.GetConfig(),
			}
		}
	}

	return domain.Workflow{
		ID:          pbWf.GetId(),
		Name:        pbWf.GetName(),
		Description: pbWf.GetDescription(),
		Trigger:     domainTrigger,
		Graph:       domainGraph,
	}
}

func ToProtoWorkflow(dwf domain.Workflow) *pb.Workflow {
	var pbNodes []*pb.Node
	for _, n := range dwf.Graph.Nodes {
		pbNodes = append(pbNodes, &pb.Node{
			Id:        n.ID,
			Type:      n.Type,
			OnSuccess: n.OnSuccess,
			Config:    n.Config,
		})
	}

	return &pb.Workflow{
		Id:          dwf.ID,
		Name:        dwf.Name,
		Description: dwf.Description,
		Trigger: &pb.Trigger{
			Type:   dwf.Trigger.Type,
			Config: dwf.Trigger.Config,
		},
		Graph: &pb.Graph{
			Entrypoint: dwf.Graph.Entrypoint,
			Nodes:      pbNodes,
		},
	}
}
