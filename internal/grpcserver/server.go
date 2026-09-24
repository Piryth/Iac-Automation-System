package grpcserver

import (
	"context"
	"log"
	pb "piryth/ias/api/proto"
)

type Server struct {
	pb.UnimplementedAutomateServiceServer
}

// NewServer est le constructeur (Factory) de ton service gRPC
func NewServer() *Server {
	return &Server{}
}

func (s *Server) ApplyWorkflow(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	log.Printf("[grpcserver] request recieved with id %s", req.GetWorkflow().GetId())

	return &pb.ApplyResponse{
		Success: true,
		Message: "not implemented",
	}, nil
}
