package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "piryth/ias/api/proto"
	"piryth/ias/internal/parser"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Deploys CUE workflows to the execution engine",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Evaluate the entire CUE project in the current working directory.
		// The parser merges the repo's CUE files with the binary's embedded schema.
		fmt.Println("Evaluating CUE project in current directory...")
		config, workflowsProto, err := parser.ParseProject(".")
		if err != nil {
			log.Fatalf("Error: CUE compilation failed: %v", err)
		}

		if config.Server.Address == "" {
			log.Fatalf("Error: Could not find a valid 'server.address' in your CUE configuration.")
		}

		if len(workflowsProto) == 0 {
			log.Fatalf("Warning: No workflows were defined in this CUE package.")
		}

		fmt.Printf("Found destination server: %s\n", config.Server.Address)
		fmt.Printf("Parsed %d workflow(s) to deploy.\n", len(workflowsProto))

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

		client := pb.NewAutomateServiceClient(conn)

		// 3. Send each compiled Protobuf workflow to the server
		for _, wf := range workflowsProto {
			fmt.Printf("Deploying workflow '%s' (ID: %s)...\n", wf.Name, wf.Id)

			req := &pb.ApplyRequest{Workflow: wf}
			res, err := client.ApplyWorkflow(ctx, req)

			if err != nil {
				fmt.Printf("Error: Failed to send workflow '%s' to server: %v\n", wf.Id, err)
			} else if res.Success {
				fmt.Printf("Success: Server acknowledged workflow '%s': %s\n", wf.Id, res.Message)
			} else {
				fmt.Printf("Error: Server rejected workflow '%s': %s\n", wf.Id, res.Message)
			}
		}
	},
}
