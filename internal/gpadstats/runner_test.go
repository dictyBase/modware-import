package gpadstats

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestRunURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := `!gpad-version 1.2
!date 2023-01-01
!comments
DDB_G0285425		involved_in	GO:0005575	PMID:26484320	ECO:0000314	dictyBase:DDB_G0285425	10090	20191111	dictyBase		
DDB_G0285427		involved_in	GO:0005576	PMID:26484321	ECO:0000315	dictyBase:DDB_G0285427	10090	20191111	dictyBase		
DDB_G0285425		involved_in	GO:0005575	PMID:26484320	ECO:0000314	dictyBase:DDB_G0285425	10090	20191111	dictyBase		
`
		zw := gzip.NewWriter(w)
		_, _ = zw.Write([]byte(raw))
		zw.Close()
	}))
	defer ts.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "gene-count-url",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "url"},
				},
				Action: RunURL,
			},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := app.Run([]string{"app", "gene-count-url", "--url", ts.URL})

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	assert.NoError(t, err)
	output := buf.String()
	// 2 unique IDs: DDB_G0285425, DDB_G0285427
	assert.Contains(t, output, "Unique Gene Count: 2")
	// Both ECO codes 314 and 315 are in the list.
	assert.Contains(t, output, "Unique Gene with ECO code Count: 2")
}
