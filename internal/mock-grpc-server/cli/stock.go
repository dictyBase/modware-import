package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/logger"
	stocksvc "github.com/dictyBase/modware-import/internal/mock-grpc-server/service/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/storage/pebble"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RunStockServer starts the mock stock gRPC server
func RunStockServer(cliCtx *cli.Context) error {
	// Setup logging
	loggerEntry, err := logger.NewCliLogger(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	loggerInstance := loggerEntry.Logger

	loggerInstance.WithFields(logrus.Fields{
		"port":             cliCtx.Int("port"),
		"log_level":        cliCtx.String("log-level"),
		"data_dir":         cliCtx.String("data-dir"),
		"strain_ontology":  cliCtx.String("strain-ontology"),
		"strain_term":      cliCtx.String("strain-term"),
		"plasmid_ontology": cliCtx.String("plasmid-ontology"),
		"plasmid_term":     cliCtx.String("plasmid-term"),
	}).Info("Starting stock mock gRPC server")

	// Initialize Pebble storage
	storageBackend, err := pebble.NewStockStorage(&pebble.Config{
		DataDir: cliCtx.String("data-dir"),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer storageBackend.Close()

	logMode := "persistent"
	if cliCtx.String("data-dir") == "" {
		logMode = "in-memory"
	}
	loggerInstance.WithField("mode", logMode).Info("Storage initialized")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register stock service
	stockService := stocksvc.NewStockService(
		storageBackend,
		&stocksvc.ServiceConfig{
			StrainOntology:  cliCtx.String("strain-ontology"),
			StrainTerm:      cliCtx.String("strain-term"),
			PlasmidOntology: cliCtx.String("plasmid-ontology"),
			PlasmidTerm:     cliCtx.String("plasmid-term"),
			Logger:          loggerInstance,
		},
	)

	stock.RegisterStockServiceServer(grpcServer, stockService)

	// Enable reflection for debugging
	if cliCtx.Bool("reflection") {
		reflection.Register(grpcServer)
		loggerInstance.Debug("gRPC reflection enabled")
	}

	// Setup listener
	address := fmt.Sprintf(":%d", cliCtx.Int("port"))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		loggerInstance.WithField(
			"signal",
			sig.String(),
		).Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	loggerInstance.WithField("address", lis.Addr().String()).
		Info("Stock gRPC server started")

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
