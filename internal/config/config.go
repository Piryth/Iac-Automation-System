package config

import "os"

type Config struct {
	DatabaseURL string
	GRPCPort    string
}

func Load() *Config {
	// Provide sensible defaults for local development
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://username:password@localhost:5432/deployment-tools?sslmode=disable"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	return &Config{
		DatabaseURL: dbURL,
		GRPCPort:    grpcPort,
	}
}
