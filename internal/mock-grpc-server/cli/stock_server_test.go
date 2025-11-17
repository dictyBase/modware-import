package cli

import (
	"context"
	"flag"
	"fmt"
	"net"
	"testing"
	"time"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// realServerFixture holds real server test dependencies
type realServerFixture struct {
	port    int
	client  stockpb.StockServiceClient
	cancel  context.CancelFunc
	cleanup func()
	dataDir string
}

// setupRealServer starts actual gRPC server on random port
func setupRealServer(t *testing.T, persistent bool) *realServerFixture {
	// Get random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	var dataDir string
	if persistent {
		dataDir = t.TempDir() // Persistent pebble storage
	}

	// Create CLI context
	app := &cli.App{}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Int("port", port, "")
	flagSet.String("log-level", "error", "")
	flagSet.String("log-format", "json", "")
	flagSet.String("log-file", "", "")
	flagSet.String("data-dir", dataDir, "")
	flagSet.Bool("reflection", true, "")
	flagSet.String("strain-ontology", "dicty_strain_property", "")
	flagSet.String("strain-term", "general strain", "")
	flagSet.String("plasmid-ontology", "plasmid_keywords", "")
	flagSet.String("plasmid-term", "cloning vector", "")

	cliCtx := cli.NewContext(app, flagSet, nil)

	// Run server in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	serverErr := make(chan error, 1)

	go func() {
		if nerr := RunStockServer(cliCtx); nerr != nil {
			serverErr <- nerr
		}
	}()

	// Wait for server to be ready
	require.NoError(t, waitForServerReady(ctx, port, 5*time.Second))

	// Create client
	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return &realServerFixture{
		port:    port,
		client:  stockpb.NewStockServiceClient(conn),
		cancel:  cancel,
		dataDir: dataDir,
		cleanup: func() {
			conn.Close()
			cancel()
			// Wait for server shutdown
			select {
			case <-serverErr:
			case <-time.After(2 * time.Second):
				t.Log("Server shutdown timeout")
			}
		},
	}
}

// waitForServerReady polls until server accepts connections
func waitForServerReady(
	ctx context.Context,
	port int,
	timeout time.Duration,
) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp",
			fmt.Sprintf("127.0.0.1:%d", port),
			100*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
	return fmt.Errorf("server not ready after %v", timeout)
}

// TestRealServer_StartupShutdown tests server startup and shutdown
func TestRealServer_StartupShutdown(t *testing.T) {
	fix := setupRealServer(t, false)
	defer fix.cleanup()

	// Verify server is responsive
	ctx := context.Background()
	strain := NewStrainBuilder().WithLabel("test").Build()

	created, err := fix.client.CreateStrain(ctx, strain)
	require.NoError(t, err)
	require.NotEmpty(t, created.Data.Id)
}

// TestRealServer_GracefulShutdown tests graceful shutdown behavior
func TestRealServer_GracefulShutdown(t *testing.T) {
	fix := setupRealServer(t, false)

	// Start long-running operation
	ctx := context.Background()
	done := make(chan struct{})

	go func() {
		_, _ = fix.client.ListStrains(ctx, &stockpb.StockParameters{})
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	fix.cancel() // Trigger shutdown

	select {
	case <-done:
		// Operation completed
	case <-time.After(2 * time.Second):
		t.Fatal("Graceful shutdown timeout")
	}
}

// setupServerWithDataDir creates a server with specified data directory
func setupServerWithDataDir(
	t *testing.T,
	dataDir string,
) (*realServerFixture, func()) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	app := &cli.App{}
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	flagSet.Int("port", port, "")
	flagSet.String("log-level", "error", "")
	flagSet.String("log-format", "json", "")
	flagSet.String("log-file", "", "")
	flagSet.String("data-dir", dataDir, "")
	flagSet.Bool("reflection", true, "")
	flagSet.String("strain-ontology", "dicty_strain_property", "")
	flagSet.String("strain-term", "general strain", "")
	flagSet.String("plasmid-ontology", "plasmid_keywords", "")
	flagSet.String("plasmid-term", "cloning vector", "")

	cliCtx := cli.NewContext(app, flagSet, nil)

	ctx, cancel := context.WithCancel(context.Background())
	serverErr := make(chan error, 1)

	go func() {
		if gerr := RunStockServer(cliCtx); gerr != nil {
			serverErr <- gerr
		}
	}()

	require.NoError(t, waitForServerReady(ctx, port, 5*time.Second))

	conn, err := grpc.NewClient(
		fmt.Sprintf("127.0.0.1:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	cleanup := func() {
		conn.Close()
		cancel()
		select {
		case <-serverErr:
		case <-time.After(2 * time.Second):
		}
	}

	return &realServerFixture{
		port:    port,
		client:  stockpb.NewStockServiceClient(conn),
		cancel:  cancel,
		dataDir: dataDir,
		cleanup: cleanup,
	}, cleanup
}

// TestRealServer_PersistentStorage tests persistent storage across restarts
func TestRealServer_PersistentStorage(t *testing.T) {
	// First server with persistent storage
	fix1 := setupRealServer(t, true)
	dataDir := fix1.dataDir

	// Create strain
	ctx := context.Background()
	strain := NewStrainBuilder().WithLabel("persistent").Build()
	created, err := fix1.client.CreateStrain(ctx, strain)
	require.NoError(t, err)
	strainID := created.Data.Id

	fix1.cleanup()
	time.Sleep(100 * time.Millisecond) // Allow port to be released

	// Second server with same data dir
	fix2, cleanup2 := setupServerWithDataDir(t, dataDir)
	defer cleanup2()

	// Verify data persisted
	retrieved, err := fix2.client.GetStrain(ctx, &stockpb.StockId{Id: strainID})
	require.NoError(t, err)
	require.Equal(t, "persistent", retrieved.Data.Attributes.Label)
}

// TestRealServer_PortConflict tests port binding behavior
func TestRealServer_PortConflict(t *testing.T) {
	fix := setupRealServer(t, false)
	defer fix.cleanup()

	// Try to bind to same port - should fail
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", fix.port))
	require.Error(t, err)
	require.Contains(t, err.Error(), "address already in use")
	if listener != nil {
		listener.Close()
	}
}

// TestRealServer_MultipleOperations tests multiple operations on real server
func TestRealServer_MultipleOperations(t *testing.T) {
	fix := setupRealServer(t, false)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strain
	strain := NewStrainBuilder().WithLabel("multi-op test").Build()
	created, err := fix.client.CreateStrain(ctx, strain)
	require.NoError(t, err)

	// Update it
	updateReq := &stockpb.StrainUpdate{
		Data: &stockpb.StrainUpdate_Data{
			Id:   created.Data.Id,
			Type: "strain",
			Attributes: &stockpb.StrainUpdateAttributes{
				UpdatedBy: "updater@dictybase.org",
				Label:     "updated label",
			},
		},
	}

	updated, err := fix.client.UpdateStrain(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, "updated label", updated.Data.Attributes.Label)

	// List it
	result, err := fix.client.ListStrains(ctx, &stockpb.StockParameters{})
	require.NoError(t, err)
	require.Len(t, result.Data, 1)

	// Remove it
	_, err = fix.client.RemoveStock(ctx, &stockpb.StockId{Id: created.Data.Id})
	require.NoError(t, err)

	// Verify gone
	_, err = fix.client.GetStrain(ctx, &stockpb.StockId{Id: created.Data.Id})
	require.Error(t, err)
}

// TestRealServer_InMemoryMode tests in-memory mode (no persistence)
func TestRealServer_InMemoryMode(t *testing.T) {
	fix := setupRealServer(t, false) // in-memory
	defer fix.cleanup()

	ctx := context.Background()

	// Create strain
	strain := NewStrainBuilder().WithLabel("in-memory").Build()
	created, err := fix.client.CreateStrain(ctx, strain)
	require.NoError(t, err)

	// Verify it exists
	retrieved, err := fix.client.GetStrain(
		ctx,
		&stockpb.StockId{Id: created.Data.Id},
	)
	require.NoError(t, err)
	require.Equal(t, "in-memory", retrieved.Data.Attributes.Label)

	// After server restart, data should be gone (but we can't easily test this
	// without more complex server lifecycle management)
}
