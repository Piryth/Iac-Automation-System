package cli

import (
	"context"
	"log"
	pb "piryth/ias/api/proto"
	"piryth/ias/internal/cli/parser"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Runs diff between repo and remote",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Evaluating current project configuration")
		config, workflowsProto, err := parser.ParseProject(".")

		// 1. Get all workflow from remote
		if config.Server.Address == "" {
			log.Fatalf("Error: Could not find a valid 'server.address' in your CUE configuration.")
		}
		log.Printf("Found destination server: %s\n", config.Server.Address)

		// 2. Initialize gRPC connection (Once for the entire batch)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(
			ctx,
			config.Server.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("Error: Could not connect to gRPC server %s: %v", config.Server.Address, err)
		}
		defer conn.Close()

		req := &pb.ListRequest{}
		client := pb.NewAutomateServiceClient(conn)
		res, err := client.ListWorkflows(ctx, req)

		log.Println(res)

		//2. Get all workflow from repo
		if err != nil {
			log.Fatalf("Error: CUE compilation failed: %v", err)
		}

		if len(workflowsProto) == 0 {
			log.Fatalf("Warning: No workflows were defined in this CUE package.")
		}
	},
}
