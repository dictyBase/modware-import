package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/mock"
	"github.com/dictyBase/modware-import/internal/server"
	"github.com/dictyBase/modware-import/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := &cli.App{
		Name:  "mock-grpc-server",
		Usage: "Mock gRPC server for feature annotation service integration testing",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   9000,
				Usage:   "Server port",
				EnvVars: []string{"GRPC_PORT"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-format",
				Value:   "json",
				Usage:   "Log format (json, text)",
				EnvVars: []string{"LOG_FORMAT"},
			},
			&cli.StringFlag{
				Name:    "log-file",
				Usage:   "Log file path (optional)",
				EnvVars: []string{"LOG_FILE"},
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {
	// Parse configuration
	cfg := &config.Config{
		Port:     ctx.Int("port"),
		LogLevel: ctx.String("log-level"),
	}

	// Setup logging
	loggerEntry, err := logger.NewCliLogger(ctx)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	loggerInstance := loggerEntry.Logger

	loggerInstance.WithFields(logrus.Fields{
		"port":      cfg.Port,
		"log_level": cfg.LogLevel,
	}).Info("Starting mock gRPC server")

	// Initialize storage
	storage := storage.NewMemoryStorage(loggerInstance)

	// Generate mock data
	mockData := mock.GenerateFeatureAnnotations()
	for _, annotation := range mockData {
		if err = storage.Create(annotation); err != nil {
			loggerInstance.WithError(err).Error("Failed to create mock annotation")
		}
	}
	loggerInstance.WithField("count", len(mockData)).
		Info("Loaded mock feature annotations")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	featureService := server.NewFeatureAnnotationServer(storage, loggerInstance)
	feature.RegisterFeatureAnnotationServiceServer(grpcServer, featureService)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Setup listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		loggerInstance.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	loggerInstance.WithField("address", lis.Addr().String()).Info("gRPC server started")

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
