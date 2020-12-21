package uniprot

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// IDCacheKey is the key for storing redis hash field value
	IDCacheKey = "UNIPROT2NAME/uniprot"
)

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
	nc := 0
	gnc := 0
	gic := 0
	urc := 0
	sc := 0
	for scanner.Scan() {
		// ignore header
		if strings.HasPrefix(scanner.Text(), "Entry") {
			continue
		}
		s := strings.Split(strings.TrimSpace(scanner.Text()), "\t")
		sl := len(s)
		switch {
		// if there is no mapping
		case sl == 1:
			nc++
		// only gene ids
		case sl == 2:
			gic++
			gs := strings.Split(s[1], ";")
			if len(gs) > 3 {
				log.Printf("unresolved line %s\t%s\n", s[0], s[1])
				urc++
			} else {
				// store in redis
				err := client.HSet(IDCacheKey, s[0], gs[0]).Err()
				if err != nil {
					return fmt.Errorf("error in setting the value in redis %s %s", s, err)
				}
			}
		// gene name
		case sl == 3:
			gnc++
			if strings.Contains(s[2], ";") {
				sc++
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
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error in scanning output %s", err)
	}
	stat := fmt.Sprintf("name:%d\tid:%d\tisoform:%d\tunresolved:%d\tnomap:%d\n", gnc, gic, sc, urc, nc)
	log.Print(stat)
	return nil
}
