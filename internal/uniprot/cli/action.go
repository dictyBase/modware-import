package cli

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/dictyBase/modware-import/internal/uniprot/client"
	rds "github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)

const (
	// UniprotCacheKey is the key for storing the Redis hash field value
	// for UniprotID -> Gene Name/ID
	UniprotCacheKey = "UNIPROT2NAME/uniprot"
	// GeneCacheKey is the key for storing the Redis has field value
	// for Gene Name/ID -> Uniprot ID
	GeneCacheKey = "GENE2UNIPROT/gene"
)

type UniProtResponse struct {
	Results []UniProtEntry `json:"results"`
}

type UniProtEntry struct {
	PrimaryAccession string                  `json:"primaryAccession"`
	CrossReferences  []UniProtCrossReference `json:"uniProtKBCrossReferences"`
}

type UniProtCrossReference struct {
	Database   string                    `json:"database"`
	ID         string                    `json:"id"`
	Properties []UniProtCrossRefProperty `json:"properties"`
}

type UniProtCrossRefProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UniprotMap struct {
	UniprotID string
	GeneID    string
	GeneSym   []string
}

// LoadUniprotMappings stores uniprot and gene name or identifier mapping in redis
func LoadUniprotMappings(cltx *cli.Context) error {
	if err := client.SetRedisClient(cltx); err != nil {
		return fmt.Errorf("error setting up Redis client: %w", err)
	}
	redisClient := registry.GetRedisClient()
	defer redisClient.Close()

	url := cltx.String("uniprot-url")
	for len(url) > 0 {
		idMaps, nextURL, err := processUniprotPage(url)
		if err != nil {
			return err
		}
		if err := loadUniprotMapsToRedis(idMaps, redisClient); err != nil {
			return err
		}
		url = nextURL
	}
	return nil
}

func processUniprotPage(url string) ([]UniprotMap, string, error) {
	resp, err := fetchUniprotData(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	uniProtResp, err := decodeUniprotResponse(resp)
	if err != nil {
		return nil, "", err
	}

	idMaps := extractUniprotMaps(uniProtResp)
	nextURL := extractNextPageURL(resp.Header.Get("Link"))

	return idMaps, nextURL, nil
}

func fetchUniprotData(urlStr string) (*http.Response, error) {
	// Validate the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Ensure the URL uses HTTPS and is from the expected domain
	if parsedURL.Scheme != "https" ||
		!strings.HasSuffix(parsedURL.Host, "uniprot.org") {
		return nil, fmt.Errorf("invalid URL scheme or host: %s", urlStr)
	}

	// Make the HTTP request
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"HTTP request failed with status code: %d",
			resp.StatusCode,
		)
	}
	return resp, nil
}

func decodeUniprotResponse(resp *http.Response) (UniProtResponse, error) {
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return UniProtResponse{}, fmt.Errorf(
			"Error creating gzip reader: %v",
			err,
		)
	}
	defer gzReader.Close()

	var uniProtResp UniProtResponse
	if err := json.NewDecoder(gzReader).Decode(&uniProtResp); err != nil {
		return UniProtResponse{}, fmt.Errorf("Error decoding JSON: %v", err)
	}
	return uniProtResp, nil
}

func extractUniprotMaps(uniProtResp UniProtResponse) []UniprotMap {
	idMaps := make([]UniprotMap, 0)
	for _, entry := range uniProtResp.Results {
		dictyID, geneNames := extractCrossReferenceInfo(entry)
		if len(dictyID) == 0 {
			continue
		}
		idMaps = append(idMaps, UniprotMap{
			UniprotID: entry.PrimaryAccession,
			GeneID:    dictyID,
			GeneSym:   geneNames,
		})
	}
	return idMaps
}

func extractCrossReferenceInfo(entry UniProtEntry) (string, []string) {
	var dictyID string
	geneNames := make([]string, 0)
	for _, ref := range entry.CrossReferences {
		if ref.Database == "dictyBase" {
			dictyID = ref.ID
		}
		for _, prop := range ref.Properties {
			if prop.Key == "GeneName" {
				geneNames = append(geneNames, prop.Value)
			}
		}
	}
	return dictyID, geneNames
}

func loadUniprotMapsToRedis(maps []UniprotMap, client *rds.Client) error {
	ctx := context.Background()
	pipe := client.Pipeline()
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

func extractNextPageURL(linkHeader string) string {
	if len(linkHeader) == 0 {
		return ""
	}
	parts := strings.Split(linkHeader, ";")
	if len(parts) != 2 {
		return ""
	}
	if strings.Contains(parts[1], `rel="next"`) {
		nextURL := strings.Trim(parts[0], " <>")
		return nextURL
	}
	return ""
}
