package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupMiniredis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestRedisUniprotLoader_Load(t *testing.T) {
	mr, client := setupMiniredis(t)
	defer mr.Close()

	loader := NewRedisUniprotLoader(client)

	testCases := []struct {
		name string
		maps []UniprotMap
	}{
		{
			name: "Single UniprotMap",
			maps: []UniprotMap{
				{
					UniprotID: "P12345",
					GeneID:    "DDB_G0123456",
					GeneSym:   []string{"geneA", "geneB"},
				},
			},
		},
		{
			name: "Multiple UniprotMaps",
			maps: []UniprotMap{
				{
					UniprotID: "P12345",
					GeneID:    "DDB_G0123456",
					GeneSym:   []string{"geneA", "geneB"},
				},
				{
					UniprotID: "Q67890",
					GeneID:    "DDB_G0789012",
					GeneSym:   []string{"geneC"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.Load(tc.maps)
			assert.NoError(t, err)

			ctx := context.Background()

			// Check UniprotID -> GeneID mapping
			for _, umap := range tc.maps {
				geneID, err := client.HGet(ctx, UniprotCacheKey, umap.UniprotID).Result()
				assert.NoError(t, err)
				assert.Equal(t, umap.GeneID, geneID)

				// Check GeneID -> UniprotID mapping
				uniprotID, err := client.HGet(ctx, GeneCacheKey, umap.GeneID).Result()
				assert.NoError(t, err)
				assert.Equal(t, umap.UniprotID, uniprotID)

				// Check GeneSym -> UniprotID mappings
				for _, sym := range umap.GeneSym {
					uniprotID, err := client.HGet(ctx, GeneCacheKey, sym).Result()
					assert.NoError(t, err)
					assert.Equal(t, umap.UniprotID, uniprotID)
				}
			}
		})
	}
}

func TestRedisUniprotLoader_Load_Error(t *testing.T) {
	mr, client := setupMiniredis(t)
	defer mr.Close()

	loader := NewRedisUniprotLoader(client)

	// Simulate a Redis error by closing the connection
	mr.Close()

	err := loader.Load([]UniprotMap{
		{
			UniprotID: "P12345",
			GeneID:    "DDB_G0123456",
			GeneSym:   []string{"geneA"},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error loading UniprotMaps to Redis")
}

func TestRedisUniprotLoader_Load_EmptyInput(t *testing.T) {
	mr, client := setupMiniredis(t)
	defer mr.Close()

	loader := NewRedisUniprotLoader(client)

	err := loader.Load([]UniprotMap{})
	assert.NoError(t, err)
}

func TestRedisUniprotLoader_Load_DuplicateEntries(t *testing.T) {
	mr, client := setupMiniredis(t)
	defer mr.Close()

	loader := NewRedisUniprotLoader(client)

	maps := []UniprotMap{
		{
			UniprotID: "P12345",
			GeneID:    "DDB_G0123456",
			GeneSym:   []string{"geneA", "geneB"},
		},
		{
			UniprotID: "P12345",
			GeneID:    "DDB_G0123456",
			GeneSym:   []string{"geneA", "geneC"},
		},
	}

	err := loader.Load(maps)
	assert.NoError(t, err)

	ctx := context.Background()

	// Check that the last entry overwrites the previous one
	geneID, err := client.HGet(ctx, UniprotCacheKey, "P12345").Result()
	assert.NoError(t, err)
	assert.Equal(t, "DDB_G0123456", geneID)

	uniprotID, err := client.HGet(ctx, GeneCacheKey, "geneA").Result()
	assert.NoError(t, err)
	assert.Equal(t, "P12345", uniprotID)

	uniprotID, err = client.HGet(ctx, GeneCacheKey, "geneC").Result()
	assert.NoError(t, err)
	assert.Equal(t, "P12345", uniprotID)

	// geneB should not exist in the GeneCacheKey
	_, err = client.HGet(ctx, GeneCacheKey, "geneB").Result()
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}
