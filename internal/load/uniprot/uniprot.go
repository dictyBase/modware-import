package uniprot

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	r "github.com/go-redis/redis/v7"
	rds "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// UniprotCacheKey is the key for storing the Redis hash field value
	// for UniprotID -> Gene Name/ID
	UniprotCacheKey = "UNIPROT2NAME/uniprot"
	// GeneCacheKey is the key for storing the Redis has field value
	// for Gene Name/ID -> Uniprot ID
	GeneCacheKey = "GENE2UNIPROT/gene"
)

// Count represents the number of each item in the dataset
type Count struct {
	noMap      int
	geneName   int
	geneID     int
	unresolved int
	isoform    int
}

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
func LoadUniprotMappings(cmd *cobra.Command, args []string) error {
	client := registry.GetRedisClient()
	defer client.Close()
	resp, err := http.Get(viper.GetString("uniprot-url"))
	if err != nil {
		return fmt.Errorf("error in retrieving from uniprot %s", err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	c := &Count{
		noMap:      0,
		geneName:   0,
		geneID:     0,
		unresolved: 0,
		isoform:    0,
	}
	for scanner.Scan() {
		// ignore header
		if strings.HasPrefix(scanner.Text(), "Entry") {
			continue
		}
		s := strings.Split(strings.TrimSpace(scanner.Text()), "\t")
		if err := readLine(s, c, client); err != nil {
			return fmt.Errorf("error in scanning line %s", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error in scanning output %s", err)
	}
	stat := fmt.Sprintf(
		"name:%d\tid:%d\tisoform:%d\tunresolved:%d\tnomap:%d\n",
		c.geneName, c.geneID, c.isoform,
		c.unresolved, c.noMap,
	)
	log.Print(stat)
	return nil
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

func handleIsoforms(s []string, c *Count, client *r.Client) error {
	c.isoform++
	ns := strings.Split(s[2], ";")
	err := client.HSet(UniprotCacheKey, s[0], ns[0]).Err()
	if err != nil {
		return fmt.Errorf("error in setting the value in redis %s %s", s, err)
	}
	err = client.HSet(GeneCacheKey, ns[0], s[0]).Err()
	if err != nil {
		return fmt.Errorf("error in setting the value in redis %s %s", s, err)
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
