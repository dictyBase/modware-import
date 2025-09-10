package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v2"
)

// GeneDataRecord represents a single gene data record
type GeneDataRecord struct {
	GeneID      string
	GeneProduct string
	Description string
}

// ParseUnknowmeData processes an HTML file containing gene data table and extracts
// DDB_G entries with their corresponding gene product and description information
func ParseUnknowmeData(cliCtx *cli.Context) error {
	inputFile := cliCtx.String("input")
	geneProductOutput := cliCtx.String("gene-product-output")
	geneDescriptionOutput := cliCtx.String("gene-description-output")

	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Parse HTML file
	records, err := parseHTMLTable(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse HTML table: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no gene data records found in the HTML file")
	}

	// Generate CSV files
	err = generateGeneProductCSV(records, geneProductOutput)
	if err != nil {
		return fmt.Errorf("failed to generate gene product CSV: %w", err)
	}

	err = generateGeneDescriptionCSV(records, geneDescriptionOutput)
	if err != nil {
		return fmt.Errorf("failed to generate gene description CSV: %w", err)
	}

	fmt.Printf("Successfully processed %d gene records\n", len(records))
	fmt.Printf("Gene products written to: %s\n", geneProductOutput)
	fmt.Printf("Gene descriptions written to: %s\n", geneDescriptionOutput)

	return nil
}

// parseHTMLTable extracts gene data from an HTML table
func parseHTMLTable(filename string) ([]GeneDataRecord, error) {
	doc, err := loadHTMLDocument(filename)
	if err != nil {
		return nil, err
	}
	var records []GeneDataRecord
	ddbGeneRegex := regexp.MustCompile(`^DDB_G\d+`)
	doc.Find("table tr").
		Each(extractRowData).
		Each(func(i int, row *goquery.Selection) {
			row.Find("td").
				FilterFunction(func(i int, cell *goquery.Selection) bool {
					if i == 0 {
						return ddbGeneRegex.MatchString(cell.Text())
					}
					return true
				})
		}).
		Each(func(i int, row *goquery.Selection) {
			rec := GeneDataRecord{}
			row.Find("td").Each(func(j int, cell *goquery.Selection) {
				switch j {
				case 0:
					rec.GeneID = cell.Text()
				case 1:
					rec.GeneProduct = cell.Text()
				case 2:
					rec.Description = cell.Text()
				}
			})
			records = append(records, rec)
		})

	return records, nil
}

func extractRowData(i int, row *goquery.Selection) {
	row.Find("td").FilterFunction(NotEmptyCell)
}

func NotEmptyCell(i int, cell *goquery.Selection) bool {
	return len(strings.TrimSpace(cell.Text())) > 0
}

// loadHTMLDocument loads and parses an HTML file
func loadHTMLDocument(filename string) (*goquery.Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// generateGeneProductCSV creates a CSV file with GeneID and gene_product columns
func generateGeneProductCSV(records []GeneDataRecord, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"GeneID", "gene_product"})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows
	for _, record := range records {
		if record.GeneProduct != "" {
			err = writer.Write([]string{record.GeneID, record.GeneProduct})
			if err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
		}
	}

	return nil
}

// generateGeneDescriptionCSV creates a CSV file with GeneID and gene_description columns
func generateGeneDescriptionCSV(
	records []GeneDataRecord,
	filename string,
) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"GeneID", "gene_description"})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows
	for _, record := range records {
		if record.Description != "" {
			err = writer.Write([]string{record.GeneID, record.Description})
			if err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
		}
	}

	return nil
}
