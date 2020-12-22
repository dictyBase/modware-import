package uniprot

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	r "github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// IDCacheKey is the key for storing redis hash field value
	IDCacheKey = "UNIPROT2NAME/uniprot"
)

// Count represents the number of each item in the dataset
type Count struct {
	noMap      int
	geneName   int
	geneID     int
	unresolved int
	isoform    int
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

func readLine(s []string, c *Count, client *r.Client) error {
	sl := len(s)
	switch {
	// if there is no mapping
	case sl == 1:
		c.noMap++
	// only gene ids
	case sl == 2:
		c.geneID++
		gs := strings.Split(s[1], ";")
		if len(gs) > 3 {
			log.Printf("unresolved line %s\t%s\n", s[0], s[1])
			c.unresolved++
		} else {
			// store in redis
			err := client.HSet(IDCacheKey, s[0], gs[0]).Err()
			if err != nil {
				return fmt.Errorf("error in setting the value in redis %s %s", s, err)
			}
		}
	// gene name
	case sl == 3:
		c.geneName++
		if strings.Contains(s[2], ";") {
			c.isoform++
			ns := strings.Split(s[2], ";")
			// store in redis
			err := client.HSet(IDCacheKey, s[0], ns[0]).Err()
			if err != nil {
				return fmt.Errorf("error in setting the value in redis %s %s", s, err)
			}
		} else {
			// store in redis
			err := client.HSet(IDCacheKey, s[0], s[2]).Err()
			if err != nil {
				return fmt.Errorf("error in setting the value in redis %s %s", s, err)
			}
		}
	default:
		log.Printf("something seriously wrong with this line %s\n", s)
	}
	return nil
}
