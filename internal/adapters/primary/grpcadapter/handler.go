package grpcadapter

import (
	"context"
	pb "piryth/ias/api/proto"
	"piryth/ias/internal/core/domain"
	"piryth/ias/internal/core/ports"
)

type Handler struct {
	pb.UnimplementedAutomateServiceServer
	service ports.WorkflowService
}

func NewHandler(service ports.WorkflowService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ApplyWorkflow(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	pbWf := req.GetWorkflow()
	if pbWf == nil {
		return &pb.ApplyResponse{Success: false, Message: "Workflow payload is empty"}, nil
	}

	// 1. Map the Trigger safely (checking for nil)
	var domainTrigger domain.Trigger
	if t := pbWf.GetTrigger(); t != nil {
		domainTrigger = domain.Trigger{
			Type:   t.GetType(),
			Config: t.GetConfig(),
		}
	}

	// 2. Map the Graph and convert the Nodes slice to a map
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

	// 3. Construct the final Core Domain entity
	wfDomain := domain.Workflow{
		ID:          pbWf.GetId(),
		Name:        pbWf.GetName(),
		Description: pbWf.GetDescription(),
		Trigger:     domainTrigger,
		Graph:       domainGraph,
	}

	// 4. Pass to your core service
	err := h.service.Apply(ctx, wfDomain)
	if err != nil {
		return &pb.ApplyResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.ApplyResponse{Success: true, Message: "Workflow successfully applied"}, nil
}

func (h *Handler) DeleteWorkflow(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := h.service.Delete(ctx, req.GetId())
	if err != nil {
		return &pb.DeleteResponse{Success: false}, nil
	}
	return &pb.DeleteResponse{Success: true}, nil
}

func (h *Handler) ListWorkflows(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	wfs, err := h.service.List(ctx)
	if err != nil {
		return &pb.ListResponse{Success: false, Message: err.Error()}, nil
	}

	var pbWfs []*pb.Workflow
	for _, wf := range wfs {
		pbWfs = append(pbWfs, ToProtoWorkflow(wf))
	}

	return &pb.ListResponse{Success: false, Message: "success", Workflows: pbWfs}, nil
}
