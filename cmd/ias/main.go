package main

import (
	"log"
	"net"

	pb "piryth/ias/api/proto"
	"piryth/ias/internal/grpcserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Rcp server initialization
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s : %v", port, err)
	}
	grpcEngine := grpc.NewServer()
	server := grpcserver.NewServer()

	pb.RegisterAutomateServiceServer(grpcEngine, server)

	reflection.Register(grpcEngine)

	log.Printf("Started grpc server on port %s", port)

	if err := grpcEngine.Serve(lis); err != nil {
		log.Fatalf("Server failed to start : %v", err)
	}
}
