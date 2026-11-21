package cli

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestStockTestConnStartsCustomTransportBeforeRPC(t *testing.T) {
	listener := bufconn.Listen(bufSize)
	t.Cleanup(func() { _ = listener.Close() })
	server := grpc.NewServer()
	t.Cleanup(server.Stop)
	go func() { _ = server.Serve(listener) }()
	addresses := make(chan string, 1)
	conn, err := newStockTestConn(func(ctx context.Context, address string) (net.Conn, error) {
		select {
		case addresses <- address:
		default:
		}
		return listener.DialContext(ctx)
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case address := <-addresses:
		require.Equal(t, "bufnet", address)
	case <-ctx.Done():
		t.Fatal("custom transport was not started before first RPC")
	}
	_, err = grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	require.Equal(t, codes.Unimplemented, status.Code(err), "RPC must reach bufconn server")
}
