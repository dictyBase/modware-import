package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func setupTestServer(t *testing.T) (feature.FeatureAnnotationServiceClient, func()) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	storage := storage.NewMemoryStorage(logger)
	server := NewFeatureAnnotationServer(storage, logger)

	grpcServer := grpc.NewServer()
	feature.RegisterFeatureAnnotationServiceServer(grpcServer, server)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := feature.NewFeatureAnnotationServiceClient(conn)

	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
	}

	return client, cleanup
}

func TestCreateFeatureAnnotation(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	req := &feature.NewFeatureAnnotation{
		Type: "gene",
		Id:   "TEST_001",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:         "testGene",
			Publications: []string{"10.1000/test.2023.001"},
			Pubmed:       []string{"12345678"},
		},
		CreatedBy: "test@dictybase.org",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp, err := client.CreateFeatureAnnotation(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "TEST_001", resp.Id)
	assert.Equal(t, "testGene", resp.Attributes.Name)
	assert.Equal(t, "test@dictybase.org", resp.CreatedBy)
	assert.False(t, resp.IsObsolete)
}

func TestGetFeatureAnnotation(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	// First create an annotation
	createReq := &feature.NewFeatureAnnotation{
		Type: "gene",
		Id:   "TEST_002",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "testGene2",
		},
		CreatedBy: "test@dictybase.org",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := client.CreateFeatureAnnotation(ctx, createReq)
	require.NoError(t, err)

	// Now get it
	getReq := &feature.FeatureAnnotationId{Id: "TEST_002"}
	resp, err := client.GetFeatureAnnotation(ctx, getReq)
	require.NoError(t, err)
	assert.Equal(t, "TEST_002", resp.Id)
	assert.Equal(t, "testGene2", resp.Attributes.Name)
}

func TestAddTag(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	// First create an annotation
	createReq := &feature.NewFeatureAnnotation{
		Type: "gene",
		Id:   "TEST_003",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "testGene3",
		},
		CreatedBy: "test@dictybase.org",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := client.CreateFeatureAnnotation(ctx, createReq)
	require.NoError(t, err)

	// Add a tag
	addTagReq := &feature.AddTagRequest{
		Id: "TEST_003",
		Tag: &feature.TagPropertyCreate{
			Tag:       "function",
			Value:     "test function",
			CreatedBy: "test@dictybase.org",
		},
	}

	resp, err := client.AddTag(ctx, addTagReq)
	require.NoError(t, err)
	assert.Len(t, resp.Attributes.Properties, 1)
	assert.Equal(t, "function", resp.Attributes.Properties[0].Tag)
	assert.Equal(t, "test function", resp.Attributes.Properties[0].Value)
}

func TestListFeatureAnnotationsByPubmedId(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	// Create annotations with the same PubMed ID
	pubmedId := "99999999"

	for i := 0; i < 3; i++ {
		createReq := &feature.NewFeatureAnnotation{
			Type: "gene",
			Id:   fmt.Sprintf("TEST_PUBMED_%d", i),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:   fmt.Sprintf("testGene%d", i),
				Pubmed: []string{pubmedId},
			},
			CreatedBy: "test@dictybase.org",
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err := client.CreateFeatureAnnotation(ctx, createReq)
		require.NoError(t, err)
	}

	// Query by PubMed ID
	listReq := &feature.PubmedId{Id: pubmedId}
	resp, err := client.ListFeatureAnnotationsByPubmedId(ctx, listReq)
	require.NoError(t, err)
	assert.Len(t, resp.Data, 3)
}

func TestListFeatureAnnotationsByDOI(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	// Create annotations with the same DOI
	doi := "10.1000/test.2023.999"

	for i := 0; i < 2; i++ {
		createReq := &feature.NewFeatureAnnotation{
			Type: "gene",
			Id:   fmt.Sprintf("TEST_DOI_%d", i),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:         fmt.Sprintf("testGene%d", i),
				Publications: []string{doi},
			},
			CreatedBy: "test@dictybase.org",
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err := client.CreateFeatureAnnotation(ctx, createReq)
		require.NoError(t, err)
	}

	// Query by DOI
	listReq := &feature.DOI{Id: doi}
	resp, err := client.ListFeatureAnnotationsByDOI(ctx, listReq)
	require.NoError(t, err)
	assert.Len(t, resp.Data, 2)
}

func TestDeleteFeatureAnnotation(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()
	now := timestamppb.New(time.Now())

	// Create an annotation
	createReq := &feature.NewFeatureAnnotation{
		Type: "gene",
		Id:   "TEST_DELETE",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "testGeneDelete",
		},
		CreatedBy: "test@dictybase.org",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := client.CreateFeatureAnnotation(ctx, createReq)
	require.NoError(t, err)

	// Soft delete
	deleteReq := &feature.DeleteFeatureAnnotationRequest{
		Id:    "TEST_DELETE",
		Purge: false,
	}

	_, err = client.DeleteFeatureAnnotation(ctx, deleteReq)
	require.NoError(t, err)

	// Try to get - should fail because it's obsolete
	getReq := &feature.FeatureAnnotationId{Id: "TEST_DELETE"}
	_, err = client.GetFeatureAnnotation(ctx, getReq)
	assert.Error(t, err)
}

func TestValidationErrors(t *testing.T) {
	client, cleanup := setupTestServer(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("invalid email", func(t *testing.T) {
		req := &feature.NewFeatureAnnotation{
			Type: "gene",
			Id:   "TEST_INVALID",
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "testGene",
			},
			CreatedBy: "invalid-email", // Invalid email
		}

		_, err := client.CreateFeatureAnnotation(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a valid email")
	})

	t.Run("invalid DOI", func(t *testing.T) {
		req := &feature.NewFeatureAnnotation{
			Type: "gene",
			Id:   "TEST_INVALID_DOI",
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:         "testGene",
				Publications: []string{"invalid-doi-format"},
			},
			CreatedBy: "test@dictybase.org",
		}

		_, err := client.CreateFeatureAnnotation(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid DOI format")
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := &feature.NewFeatureAnnotation{
			Type: "gene",
			// Missing ID, Attributes, CreatedBy
		}

		_, err := client.CreateFeatureAnnotation(ctx, req)
		assert.Error(t, err)
	})
}
