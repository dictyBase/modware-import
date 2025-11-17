package cli

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/service/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/storage/pebble"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

// stockTestFixture holds test dependencies
type stockTestFixture struct {
	client  stockpb.StockServiceClient
	cleanup func()
}

// setupStockTest creates in-memory bufconn server for fast testing
func setupStockTest(t *testing.T) *stockTestFixture {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(io.Discard) // Silent during tests

	// In-memory Pebble storage
	storage, err := pebble.NewStockStorage(&pebble.Config{
		DataDir: "", // Empty = in-memory
	})
	require.NoError(t, err)

	// Create service with test config
	service := stock.NewService(storage, &stock.ServiceConfig{
		StrainOntology:  "dicty_strain_property",
		StrainTerm:      "general strain",
		PlasmidOntology: "plasmid_keywords",
		PlasmidTerm:     "cloning vector",
		Logger:          logger,
	})

	// Setup bufconn gRPC server
	grpcServer := grpc.NewServer()
	stockpb.RegisterStockServiceServer(grpcServer, service)

	lis := bufconn.Listen(bufSize)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("Server exited: %v", err)
		}
	}()

	// Create client
	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return &stockTestFixture{
		client: stockpb.NewStockServiceClient(conn),
		cleanup: func() {
			conn.Close()
			lis.Close()
			grpcServer.Stop()
			storage.Close()
		},
	}
}

// StrainBuilder builds test strain requests
type StrainBuilder struct {
	req *stockpb.NewStrain
}

// NewStrainBuilder creates a new strain builder with sensible defaults
func NewStrainBuilder() *StrainBuilder {
	return &StrainBuilder{
		req: &stockpb.NewStrain{
			Data: &stockpb.NewStrain_Data{
				Type: "strain",
				Attributes: &stockpb.NewStrainAttributes{
					CreatedBy: "test@dictybase.org",
					UpdatedBy: "test@dictybase.org",
					Depositor: "Test Depositor",
					Species:   "Dictyostelium discoideum",
				},
			},
		},
	}
}

// WithLabel sets the strain label
func (b *StrainBuilder) WithLabel(label string) *StrainBuilder {
	b.req.Data.Attributes.Label = label
	return b
}

// WithDepositor sets the depositor
func (b *StrainBuilder) WithDepositor(depositor string) *StrainBuilder {
	b.req.Data.Attributes.Depositor = depositor
	return b
}

// WithParent sets the parent strain ID
func (b *StrainBuilder) WithParent(parentID string) *StrainBuilder {
	b.req.Data.Attributes.Parent = parentID
	return b
}

// Build returns the constructed request
func (b *StrainBuilder) Build() *stockpb.NewStrain {
	return b.req
}

// PlasmidBuilder builds test plasmid requests
type PlasmidBuilder struct {
	req *stockpb.NewPlasmid
}

// NewPlasmidBuilder creates a new plasmid builder with sensible defaults
func NewPlasmidBuilder() *PlasmidBuilder {
	return &PlasmidBuilder{
		req: &stockpb.NewPlasmid{
			Data: &stockpb.NewPlasmid_Data{
				Type: "plasmid",
				Attributes: &stockpb.NewPlasmidAttributes{
					CreatedBy: "test@dictybase.org",
					UpdatedBy: "test@dictybase.org",
					Depositor: "Test Depositor",
				},
			},
		},
	}
}

// WithName sets the plasmid name
func (b *PlasmidBuilder) WithName(name string) *PlasmidBuilder {
	b.req.Data.Attributes.Name = name
	return b
}

// Build returns the constructed request
func (b *PlasmidBuilder) Build() *stockpb.NewPlasmid {
	return b.req
}

// TestCreateAndGetStrain tests basic strain creation and retrieval
func TestCreateAndGetStrain(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strain
	req := NewStrainBuilder().
		WithLabel("axeA2 axeB2").
		WithDepositor("Costanza").
		Build()

	created, err := fix.client.CreateStrain(ctx, req)
	require.NoError(t, err)
	require.Regexp(t, `^DBS\d{7}$`, created.Data.Id)
	require.Equal(t, "axeA2 axeB2", created.Data.Attributes.Label)
	require.Equal(t, "Costanza", created.Data.Attributes.Depositor)
	require.Equal(
		t,
		"general strain",
		created.Data.Attributes.DictyStrainProperty,
	)
	require.NotNil(t, created.Data.Attributes.CreatedAt)
	require.NotNil(t, created.Data.Attributes.UpdatedAt)

	// Get strain
	retrieved, err := fix.client.GetStrain(
		ctx,
		&stockpb.StockId{Id: created.Data.Id},
	)
	require.NoError(t, err)
	require.Equal(t, created.Data.Id, retrieved.Data.Id)
	require.Equal(t, "axeA2 axeB2", retrieved.Data.Attributes.Label)
}

// TestUpdateStrain tests strain update functionality
func TestUpdateStrain(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create initial strain
	created, err := fix.client.CreateStrain(ctx,
		NewStrainBuilder().WithLabel("original").Build())
	require.NoError(t, err)

	// Update strain
	updateReq := &stockpb.StrainUpdate{
		Data: &stockpb.StrainUpdate_Data{
			Id:   created.Data.Id,
			Type: "strain",
			Attributes: &stockpb.StrainUpdateAttributes{
				UpdatedBy: "updater@dictybase.org",
				Label:     "updated label",
				Summary:   "Added summary",
			},
		},
	}

	updated, err := fix.client.UpdateStrain(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, created.Data.Id, updated.Data.Id)
	require.Equal(t, "updated label", updated.Data.Attributes.Label)
	require.Equal(t, "Added summary", updated.Data.Attributes.Summary)
	require.True(t, updated.Data.Attributes.UpdatedAt.AsTime().After(
		created.Data.Attributes.UpdatedAt.AsTime()))
}

// TestLoadStrain tests loading a strain with a specific ID
func TestLoadStrain(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Load strain with custom ID
	now := timestamppb.Now()
	loadReq := &stockpb.ExistingStrain{
		Data: &stockpb.ExistingStrain_Data{
			Type: "strain",
			Id:   "DBS9999999",
			Attributes: &stockpb.ExistingStrainAttributes{
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: "test@dictybase.org",
				UpdatedBy: "test@dictybase.org",
				Depositor: "Loaded Depositor",
				Label:     "loaded strain",
				Species:   "Dictyostelium discoideum",
			},
		},
	}

	loaded, err := fix.client.LoadStrain(ctx, loadReq)
	require.NoError(t, err)
	require.Equal(t, "DBS9999999", loaded.Data.Id)
	require.Equal(t, "loaded strain", loaded.Data.Attributes.Label)

	// Verify it can be retrieved
	retrieved, err := fix.client.GetStrain(
		ctx,
		&stockpb.StockId{Id: "DBS9999999"},
	)
	require.NoError(t, err)
	require.Equal(t, "DBS9999999", retrieved.Data.Id)
}

// TestRemoveStrain tests strain deletion
func TestRemoveStrain(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strain
	created, err := fix.client.CreateStrain(ctx,
		NewStrainBuilder().WithLabel("to delete").Build())
	require.NoError(t, err)

	// Remove strain
	_, err = fix.client.RemoveStock(ctx, &stockpb.StockId{Id: created.Data.Id})
	require.NoError(t, err)

	// Verify it's gone
	_, err = fix.client.GetStrain(ctx, &stockpb.StockId{Id: created.Data.Id})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

// TestCreateAndGetPlasmid tests basic plasmid creation and retrieval
func TestCreateAndGetPlasmid(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	req := NewPlasmidBuilder().
		WithName("pDV-CFPC-act15").
		Build()

	created, err := fix.client.CreatePlasmid(ctx, req)
	require.NoError(t, err)
	require.Regexp(t, `^DBP\d{7}$`, created.Data.Id)
	require.Equal(t, "pDV-CFPC-act15", created.Data.Attributes.Name)
	require.Equal(
		t,
		"cloning vector",
		created.Data.Attributes.DictyPlasmidProperty,
	)

	// Get plasmid
	retrieved, err := fix.client.GetPlasmid(
		ctx,
		&stockpb.StockId{Id: created.Data.Id},
	)
	require.NoError(t, err)
	require.Equal(t, created.Data.Id, retrieved.Data.Id)
}

// TestUpdatePlasmid tests plasmid update functionality
func TestUpdatePlasmid(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	created, err := fix.client.CreatePlasmid(ctx,
		NewPlasmidBuilder().WithName("original").Build())
	require.NoError(t, err)

	updateReq := &stockpb.PlasmidUpdate{
		Data: &stockpb.PlasmidUpdate_Data{
			Id:   created.Data.Id,
			Type: "plasmid",
			Attributes: &stockpb.PlasmidUpdateAttributes{
				UpdatedBy: "updater@dictybase.org",
				Name:      "updated plasmid",
				Summary:   "Updated summary",
			},
		},
	}

	updated, err := fix.client.UpdatePlasmid(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, "updated plasmid", updated.Data.Attributes.Name)
}

// TestLoadPlasmid tests loading a plasmid with a specific ID
func TestLoadPlasmid(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	now2 := timestamppb.Now()
	loadReq := &stockpb.ExistingPlasmid{
		Data: &stockpb.ExistingPlasmid_Data{
			Type: "plasmid",
			Id:   "DBP8888888",
			Attributes: &stockpb.ExistingPlasmidAttributes{
				CreatedAt: now2,
				UpdatedAt: now2,
				CreatedBy: "test@dictybase.org",
				UpdatedBy: "test@dictybase.org",
				Depositor: "Plasmid Depositor",
				Name:      "loaded plasmid",
			},
		},
	}

	loaded, err := fix.client.LoadPlasmid(ctx, loadReq)
	require.NoError(t, err)
	require.Equal(t, "DBP8888888", loaded.Data.Id)
}

// TestStrainParentChild tests parent-child strain relationships
func TestStrainParentChild(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create parent strain
	parent, err := fix.client.CreateStrain(ctx,
		NewStrainBuilder().WithLabel("parent strain").Build())
	require.NoError(t, err)

	// Create child strain
	child, err := fix.client.CreateStrain(ctx,
		NewStrainBuilder().
			WithLabel("child strain").
			WithParent(parent.Data.Id).
			Build())
	require.NoError(t, err)
	require.Equal(t, parent.Data.Id, child.Data.Attributes.Parent)

	// Verify parent can be retrieved
	retrievedParent, err := fix.client.GetStrain(ctx,
		&stockpb.StockId{Id: parent.Data.Id})
	require.NoError(t, err)
	require.Equal(t, "parent strain", retrievedParent.Data.Attributes.Label)
}

// TestListStrains_FilterByDepositor tests filtering strains by depositor
func TestListStrains_FilterByDepositor(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strains with different depositors
	strains := []struct {
		depositor string
		label     string
	}{
		{"Costanza", "strain1"},
		{"Costanza", "strain2"},
		{"Benes", "strain3"},
	}

	for _, s := range strains {
		_, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().
				WithLabel(s.label).
				WithDepositor(s.depositor).
				Build())
		require.NoError(t, err)
	}

	// Filter by depositor
	params := &stockpb.StockParameters{
		Filter: "depositor===Costanza",
	}

	result, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, result.Data, 2)

	for _, strain := range result.Data {
		require.Equal(t, "Costanza", strain.Attributes.Depositor)
	}
}

// TestListStrains_FilterBySpecies tests filtering strains by species
func TestListStrains_FilterBySpecies(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strains with different species
	species := []string{
		"Dictyostelium discoideum",
		"Dictyostelium discoideum",
		"Dictyostelium purpureum",
	}

	for i, sp := range species {
		req := NewStrainBuilder().
			WithLabel(fmt.Sprintf("strain%d", i)).
			Build()
		req.Data.Attributes.Species = sp
		_, err := fix.client.CreateStrain(ctx, req)
		require.NoError(t, err)
	}

	// Filter by species (contains)
	params := &stockpb.StockParameters{
		Filter: "species=@=purpureum",
	}

	result, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	require.Contains(t, result.Data[0].Attributes.Species, "purpureum")
}

// TestListStrains_ComplexFilter tests complex AND filter combinations
func TestListStrains_ComplexFilter(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create test data
	testCases := []struct {
		depositor string
		label     string
		species   string
	}{
		{"Costanza", "axe mutant", "Dictyostelium discoideum"},
		{"Costanza", "wild type", "Dictyostelium discoideum"},
		{"Benes", "axe mutant", "Dictyostelium discoideum"},
		{"Benes", "car mutant", "Dictyostelium purpureum"},
	}

	for _, tc := range testCases {
		req := NewStrainBuilder().
			WithLabel(tc.label).
			WithDepositor(tc.depositor).
			Build()
		req.Data.Attributes.Species = tc.species
		_, err := fix.client.CreateStrain(ctx, req)
		require.NoError(t, err)
	}

	// AND filter: depositor===Costanza AND label contains "axe"
	params := &stockpb.StockParameters{
		Filter: "depositor===Costanza;label=@=axe",
	}

	result, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	require.Equal(t, "Costanza", result.Data[0].Attributes.Depositor)
	require.Contains(t, result.Data[0].Attributes.Label, "axe")
}

// TestListStrains_Pagination tests basic pagination functionality
func TestListStrains_Pagination(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create 25 strains
	totalStrains := 25
	for i := 1; i <= totalStrains; i++ {
		_, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().
				WithLabel(fmt.Sprintf("strain%03d", i)).
				Build())
		require.NoError(t, err)
	}

	// First page (limit 10)
	params := &stockpb.StockParameters{
		Limit: 10,
	}

	page1, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, page1.Data, 10)
	require.NotZero(t, page1.Meta.NextCursor)

	// Second page
	params.Cursor = page1.Meta.NextCursor
	page2, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, page2.Data, 10)
	require.NotZero(t, page2.Meta.NextCursor)

	// Third page (last 5)
	params.Cursor = page2.Meta.NextCursor
	page3, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, page3.Data, 5)
	require.Zero(t, page3.Meta.NextCursor) // No more pages

	// Verify no duplicates across pages
	allIDs := make(map[string]bool)
	for _, strain := range append(append(page1.Data, page2.Data...), page3.Data...) {
		require.False(t, allIDs[strain.Id], "Duplicate ID found")
		allIDs[strain.Id] = true
	}
	require.Len(t, allIDs, totalStrains)
}

// TestListStrains_PaginationWithFilter tests pagination combined with filtering
func TestListStrains_PaginationWithFilter(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create 15 Costanza strains and 10 Benes strains
	for i := 1; i <= 15; i++ {
		_, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().
				WithDepositor("Costanza").
				WithLabel(fmt.Sprintf("costanza%02d", i)).
				Build())
		require.NoError(t, err)
	}

	for i := 1; i <= 10; i++ {
		_, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().
				WithDepositor("Benes").
				WithLabel(fmt.Sprintf("benes%02d", i)).
				Build())
		require.NoError(t, err)
	}

	// Paginate through Costanza strains only
	params := &stockpb.StockParameters{
		Filter: "depositor===Costanza",
		Limit:  10,
	}

	page1, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, page1.Data, 10)

	params.Cursor = page1.Meta.NextCursor
	page2, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Len(t, page2.Data, 5)

	// All results should be Costanza
	for _, strain := range append(page1.Data, page2.Data...) {
		require.Equal(t, "Costanza", strain.Attributes.Depositor)
	}
}

// TestListStrainsByIds tests bulk retrieval by IDs
func TestListStrainsByIds(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create 5 strains
	var createdIDs []string
	for i := 1; i <= 5; i++ {
		created, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().WithLabel(fmt.Sprintf("strain%d", i)).Build())
		require.NoError(t, err)
		createdIDs = append(createdIDs, created.Data.Id)
	}

	// Request subset of IDs
	requestIDs := []string{createdIDs[0], createdIDs[2], createdIDs[4]}
	result, err := fix.client.ListStrainsByIds(ctx,
		&stockpb.StockIdList{Id: requestIDs})
	require.NoError(t, err)
	require.Len(t, result.Data, 3)

	// Verify correct strains returned
	returnedIDs := make(map[string]bool)
	for _, strain := range result.Data {
		returnedIDs[strain.Id] = true
	}

	for _, id := range requestIDs {
		require.True(t, returnedIDs[id], "Expected ID not in result: %s", id)
	}
}

// TestListPlasmids tests basic plasmid listing
func TestListPlasmids(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create plasmids
	for i := 1; i <= 5; i++ {
		_, err := fix.client.CreatePlasmid(ctx,
			NewPlasmidBuilder().
				WithName(fmt.Sprintf("plasmid%d", i)).
				Build())
		require.NoError(t, err)
	}

	// List all
	result, err := fix.client.ListPlasmids(ctx, &stockpb.StockParameters{})
	require.NoError(t, err)
	require.Len(t, result.Data, 5)
}

// TestListPlasmids_FilterByDepositor tests filtering plasmids by depositor
func TestListPlasmids_FilterByDepositor(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	depositors := []string{"Alice", "Alice", "Bob"}
	for i, dep := range depositors {
		req := NewPlasmidBuilder().WithName(fmt.Sprintf("plasmid%d", i)).Build()
		req.Data.Attributes.Depositor = dep
		_, err := fix.client.CreatePlasmid(ctx, req)
		require.NoError(t, err)
	}

	params := &stockpb.StockParameters{
		Filter: "depositor===Alice",
	}

	result, err := fix.client.ListPlasmids(ctx, params)
	require.NoError(t, err)
	require.Len(t, result.Data, 2)
}

// TestCreateStrain_ValidationErrors tests validation error handling
func TestCreateStrain_ValidationErrors(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	tests := []struct {
		name      string
		mutate    func(*stockpb.NewStrain)
		errString string
	}{
		{
			name: "invalid email - created_by",
			mutate: func(req *stockpb.NewStrain) {
				req.Data.Attributes.CreatedBy = "not-an-email"
			},
			errString: "email",
		},
		{
			name: "invalid email - updated_by",
			mutate: func(req *stockpb.NewStrain) {
				req.Data.Attributes.UpdatedBy = "invalid"
			},
			errString: "email",
		},
		{
			name: "missing required field - depositor",
			mutate: func(req *stockpb.NewStrain) {
				req.Data.Attributes.Depositor = ""
			},
			errString: "required",
		},
		{
			name: "missing required field - species",
			mutate: func(req *stockpb.NewStrain) {
				req.Data.Attributes.Species = ""
			},
			errString: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewStrainBuilder().WithLabel("test").Build()
			tt.mutate(req)

			_, err := fix.client.CreateStrain(ctx, req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errString)
		})
	}
}

// TestCreatePlasmid_ValidationErrors tests plasmid validation errors
func TestCreatePlasmid_ValidationErrors(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	tests := []struct {
		name      string
		mutate    func(*stockpb.NewPlasmid)
		errString string
	}{
		{
			name: "invalid email - created_by",
			mutate: func(req *stockpb.NewPlasmid) {
				req.Data.Attributes.CreatedBy = "bad-email"
			},
			errString: "email",
		},
		{
			name: "missing depositor",
			mutate: func(req *stockpb.NewPlasmid) {
				req.Data.Attributes.Depositor = ""
			},
			errString: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewPlasmidBuilder().WithName("test").Build()
			tt.mutate(req)

			_, err := fix.client.CreatePlasmid(ctx, req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errString)
		})
	}
}

// TestGetStrain_NotFound tests NotFound error for strain
func TestGetStrain_NotFound(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	_, err := fix.client.GetStrain(ctx, &stockpb.StockId{Id: "DBS9999999"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")

	// Verify gRPC status code
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}

// TestGetPlasmid_NotFound tests NotFound error for plasmid
func TestGetPlasmid_NotFound(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	_, err := fix.client.GetPlasmid(ctx, &stockpb.StockId{Id: "DBP9999999"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}

// TestUpdateStrain_NotFound tests update on non-existent strain
func TestUpdateStrain_NotFound(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	updateReq := &stockpb.StrainUpdate{
		Data: &stockpb.StrainUpdate_Data{
			Id:   "DBS9999999",
			Type: "strain",
			Attributes: &stockpb.StrainUpdateAttributes{
				UpdatedBy: "test@dictybase.org",
				Label:     "updated",
			},
		},
	}

	_, err := fix.client.UpdateStrain(ctx, updateReq)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

// TestRemoveStock_NotFound tests removal of non-existent stock
func TestRemoveStock_NotFound(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	_, err := fix.client.RemoveStock(ctx, &stockpb.StockId{Id: "DBS9999999"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

// TestCreateStrain_SequentialIDs tests sequential ID generation
func TestCreateStrain_SequentialIDs(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create 3 strains and verify IDs are sequential
	var ids []string
	for i := range 3 {
		created, err := fix.client.CreateStrain(ctx,
			NewStrainBuilder().WithLabel(fmt.Sprintf("strain%d", i)).Build())
		require.NoError(t, err)
		ids = append(ids, created.Data.Id)
	}

	// Verify format and sequence
	require.Equal(t, "DBS0000001", ids[0])
	require.Equal(t, "DBS0000002", ids[1])
	require.Equal(t, "DBS0000003", ids[2])
}

// TestCreatePlasmid_SequentialIDs tests plasmid sequential ID generation
func TestCreatePlasmid_SequentialIDs(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	var ids []string
	for i := range 3 {
		created, err := fix.client.CreatePlasmid(ctx,
			NewPlasmidBuilder().WithName(fmt.Sprintf("plasmid%d", i)).Build())
		require.NoError(t, err)
		ids = append(ids, created.Data.Id)
	}

	require.Equal(t, "DBP0000001", ids[0])
	require.Equal(t, "DBP0000002", ids[1])
	require.Equal(t, "DBP0000003", ids[2])
}

// TestListStrains_EmptyResult tests listing when no strains exist
func TestListStrains_EmptyResult(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// List when no strains exist
	result, err := fix.client.ListStrains(ctx, &stockpb.StockParameters{})
	require.NoError(t, err)
	require.Empty(t, result.Data)
	require.Zero(t, result.Meta.NextCursor)
}

// TestListStrains_FilterNoMatches tests filter with no matches
func TestListStrains_FilterNoMatches(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create some strains
	_, err := fix.client.CreateStrain(ctx,
		NewStrainBuilder().WithDepositor("Costanza").Build())
	require.NoError(t, err)

	// Filter that matches nothing
	params := &stockpb.StockParameters{
		Filter: "depositor===NonExistent",
	}

	result, err := fix.client.ListStrains(ctx, params)
	require.NoError(t, err)
	require.Empty(t, result.Data)
}

// TestListStrainsByIds_EmptyList tests empty ID list
func TestListStrainsByIds_EmptyList(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	result, err := fix.client.ListStrainsByIds(
		ctx,
		&stockpb.StockIdList{Id: []string{}},
	)
	require.NoError(t, err)
	require.Empty(t, result.Data)
}

// TestListStrainsByIds_NonExistentIDs tests retrieval of non-existent IDs
func TestListStrainsByIds_NonExistentIDs(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Request non-existent IDs
	result, err := fix.client.ListStrainsByIds(
		ctx,
		&stockpb.StockIdList{
			Id: []string{"DBS9999997", "DBS9999998", "DBS9999999"},
		},
	)

	// Should return empty results, not error
	require.NoError(t, err)
	require.Empty(t, result.Data)
}

// TestDefaultTermApplication tests default term application
func TestDefaultTermApplication(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strain without specifying dicty_strain_property
	req := NewStrainBuilder().WithLabel("test").Build()
	req.Data.Attributes.DictyStrainProperty = "" // Clear default

	created, err := fix.client.CreateStrain(ctx, req)
	require.NoError(t, err)

	// Verify default term applied
	require.Equal(
		t,
		"general strain",
		created.Data.Attributes.DictyStrainProperty,
	)

	// Same for plasmid
	plasmidReq := NewPlasmidBuilder().WithName("test").Build()
	plasmidReq.Data.Attributes.DictyPlasmidProperty = ""

	createdPlasmid, err := fix.client.CreatePlasmid(ctx, plasmidReq)
	require.NoError(t, err)
	require.Equal(
		t,
		"cloning vector",
		createdPlasmid.Data.Attributes.DictyPlasmidProperty,
	)
}

// TestConcurrentCreates tests concurrent strain creation
func TestConcurrentCreates(t *testing.T) {
	fix := setupStockTest(t)
	defer fix.cleanup()

	ctx := context.Background()

	// Create strains concurrently
	const numGoroutines = 10
	errChan := make(chan error, numGoroutines)
	idChan := make(chan string, numGoroutines)

	for i := range numGoroutines {
		go func(idx int) {
			created, err := fix.client.CreateStrain(ctx,
				NewStrainBuilder().
					WithLabel(fmt.Sprintf("concurrent%d", idx)).
					Build())
			if err != nil {
				errChan <- err
				return
			}
			idChan <- created.Data.Id
			errChan <- nil
		}(i)
	}

	// Collect results
	var ids []string
	for range numGoroutines {
		err := <-errChan
		require.NoError(t, err)
		if err == nil {
			ids = append(ids, <-idChan)
		}
	}

	// Verify all IDs are unique
	idMap := make(map[string]bool)
	for _, id := range ids {
		require.False(t, idMap[id], "Duplicate ID detected: %s", id)
		idMap[id] = true
	}
	require.Len(t, idMap, numGoroutines)
}
