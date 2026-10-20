package ports

import (
	"context"
	pb "piryth/ias/api/proto"
)

type GrpcServer interface {
	ApplyWorkflow(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error)
}
