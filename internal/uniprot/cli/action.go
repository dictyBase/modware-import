package cli

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	rdloader "github.com/dictyBase/modware-import/internal/uniprot/redis"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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

// LoadUniprotMappings stores uniprot and gene name or identifier mapping in redis
func LoadUniprotMappings(cltx *cli.Context) error {
	redisClient := registry.GetRedisClient()
	if redisClient == nil {
		return cli.Exit("Redis client is not set", 1)
	}
	defer redisClient.Close()

	logger := registry.GetLogger()
	url := cltx.String("uniprot-url")
	totalEntries := 0

	loader := rdloader.NewRedisUniprotLoader(redisClient)

	for len(url) > 0 {
		logger.Debugf("Processing Uniprot page: %s", url)
		idMaps, nextURL, err := processUniprotPage(url)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		if err := loader.Load(idMaps); err != nil {
			return cli.Exit(err.Error(), 1)
		}

		totalEntries += len(idMaps)
		logger.Infof(
			"Loaded %d Uniprot entries (Total: %d)",
			len(idMaps),
			totalEntries,
		)
		for _, entry := range idMaps {
			logger.WithFields(logrus.Fields{
				"UniprotID": entry.UniprotID,
				"GeneID":    entry.GeneID,
				"GeneSyms":  strings.Join(entry.GeneSym, ", "),
			}).Debug("Loaded Uniprot entry")
		}

		url = nextURL
	}

	logger.Infof(
		"Completed loading Uniprot mappings. Total entries: %d",
		totalEntries,
	)
	return nil
}

func processUniprotPage(url string) ([]rdloader.UniprotMap, string, error) {
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

func extractUniprotMaps(uniProtResp UniProtResponse) []rdloader.UniprotMap {
	idMaps := make([]rdloader.UniprotMap, 0)
	for _, entry := range uniProtResp.Results {
		dictyID, geneNames := extractCrossReferenceInfo(entry)
		if len(dictyID) == 0 {
			continue
		}
		idMaps = append(idMaps, rdloader.UniprotMap{
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
