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
			&cli.StringFlag{
				Name:    "storage-type",
				Aliases: []string{"s"},
				Value:   "leveldb",
				Usage:   "Storage backend type (leveldb uses in-memory storage, memory uses simple map)",
				EnvVars: []string{"STORAGE_TYPE"},
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
		Port:        ctx.Int("port"),
		LogLevel:    ctx.String("log-level"),
		StorageType: ctx.String("storage-type"),
	}

	// Setup logging
	loggerEntry, err := logger.NewCliLogger(ctx)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	loggerInstance := loggerEntry.Logger

	loggerInstance.WithFields(logrus.Fields{
		"port":         cfg.Port,
		"log_level":    cfg.LogLevel,
		"storage_type": cfg.StorageType,
	}).Info("Starting mock gRPC server")

	// Initialize storage with fallback
	storageBackend, err := createStorage(cfg, loggerInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Generate mock data
	mockData := mock.GenerateFeatureAnnotations()
	for _, annotation := range mockData {
		if err = storageBackend.Create(annotation); err != nil {
			loggerInstance.WithError(err).
				Error("Failed to create mock annotation")
		}
	}
	loggerInstance.WithField("count", len(mockData)).
		Info("Loaded mock feature annotations")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	featureService := server.NewFeatureAnnotationServer(
		storageBackend,
		loggerInstance,
	)
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

	loggerInstance.WithField("address", lis.Addr().String()).
		Info("gRPC server started")

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func createStorage(
	cfg *config.Config,
	logger *logrus.Logger,
) (storage.FeatureAnnotationStorage, error) {
	switch cfg.StorageType {
	case "leveldb":
		leveldbStorage, err := storage.NewLevelDBStorage(logger)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to initialize LevelDB storage: %w",
				err,
			)
		}
		logger.Debug("Using LevelDB in-memory storage")
		return leveldbStorage, nil

	case "memory":
		logger.Debug("Using simple memory storage")
		return storage.NewMemoryStorage(logger), nil

	default:
		logger.WithField("storage_type", cfg.StorageType).
			Warn("Unknown storage type, falling back to memory storage")
		return storage.NewMemoryStorage(logger), nil
	}
}
