package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	pb "piryth/ias/api/proto"
	"piryth/ias/internal/adapters/primary/cron"
	grpc2 "piryth/ias/internal/adapters/secondary/grpcadapter"
	"piryth/ias/internal/adapters/secondary/nodes"
	"piryth/ias/internal/adapters/secondary/postgres"
	"piryth/ias/internal/adapters/secondary/runner"
	"piryth/ias/internal/config"
	"piryth/ias/internal/core/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Loading configuration
	cfg := config.Load()

	// Connecting to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	fmt.Println("Connected to PostgreSQL")

	// Seeding database
	if err := postgres.InitSchema(db); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}
	fmt.Println("Database schema verified")

	// Initialising adapters
	repo := postgres.NewPostgresWorkflowRepository(db)
	osRunner := runner.NewOSRunner()

	// Initialising node registry
	registry := services.NewNodeRegistry()
	registry.Register("command", nodes.NewCommandNode(osRunner))

	workflowService := services.NewWorkflowService(repo, *registry)

	// Init primary adapters
	scheduler := cron.NewSchedulerAdapter(workflowService, repo)
	if err := scheduler.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	fmt.Println("Cron scheduler started")

	// Rcp server initialization
	grpcHandler := grpc.NewGrpcHandler(workflowService)
	grpcServer := grpc.NewServer()
	port := cfg.GRPCPort
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s : %v", port, err)
	}
	server := grpc2.NewServer()

	pb.RegisterAutomateServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Printf("Started grpc server on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server failed to start : %v", err)
	}
}
