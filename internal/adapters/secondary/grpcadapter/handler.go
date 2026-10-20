package grpcadapter

import (
	"context"
	"log"
	pb "piryth/ias/api/proto"
	"piryth/ias/internal/core/domain"
	"piryth/ias/internal/core/ports"
)

type Server struct {
	pb.UnimplementedAutomateServiceServer
}

type Handler struct {
	pb.UnimplementedAutomateServiceServer
	service ports.WorkflowService
}

func NewHandler(service ports.WorkflowService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ApplyWorkflow(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	log.Printf("[grpcserver] request recieved with id %s", req.GetWorkflow().GetId())

	log.Print(req.GetWorkflow().Description)

	wf := domain.Workflow{
		ID:          req.GetWorkflow().GetId(),
		Name:        req.GetWorkflow().GetName(),
		Description: req.GetWorkflow().GetDescription(),
		Trigger:     req.GetWorkflow().GetTrigger(),
		Graph:       req.GetWorkflow().GetGraph(),
	}

	h.service.Apply(ctx, req.Workflow)

	return &pb.ApplyResponse{
		Success: true,
		Message: "not implemented",
	}, nil
}
