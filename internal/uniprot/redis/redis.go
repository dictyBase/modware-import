package redis

import (
	"context"
	"fmt"

	rds "github.com/redis/go-redis/v9"
)

const (
	// UniprotCacheKey is the key for storing the Redis hash field value
	// for UniprotID -> Gene Name/ID
	UniprotCacheKey = "UNIPROT2NAME/uniprot"
	// GeneCacheKey is the key for storing the Redis has field value
	// for Gene Name/ID -> Uniprot ID
	GeneCacheKey = "GENE2UNIPROT/gene"
)

type UniprotMap struct {
	UniprotID string
	GeneID    string
	GeneSym   []string
}

// UniprotLoader defines the interface for loading Uniprot mappings
type UniprotLoader interface {
	Load(maps []UniprotMap) error
}

// RedisUniprotLoader implements UniprotLoader for Redis
type RedisUniprotLoader struct {
	client *rds.Client
}

func NewRedisUniprotLoader(client *rds.Client) UniprotLoader {
	return &RedisUniprotLoader{client: client}
}

func (r *RedisUniprotLoader) Load(maps []UniprotMap) error {
	ctx := context.Background()
	pipe := r.client.Pipeline()
	for _, umap := range maps {
		// Store UniprotID -> GeneID
		pipe.HSet(ctx, UniprotCacheKey, umap.UniprotID, umap.GeneID)
		// Store GeneID -> UniprotID
		pipe.HSet(ctx, GeneCacheKey, umap.GeneID, umap.UniprotID)
		// Store GeneSym -> UniprotID for each gene symbol
		for _, symbol := range umap.GeneSym {
			pipe.HSet(ctx, GeneCacheKey, symbol, umap.UniprotID)
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error loading UniprotMaps to Redis: %v", err)
	}
	return nil
}
