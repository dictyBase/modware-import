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

	doc.Find("table tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")

		// Skip rows that don't have at least 3 cells
		if cells.Length() < 3 {
			return
		}

		// Check if first cell contains DDB_G ID
		firstCellText := strings.TrimSpace(cells.Eq(0).Text())
		if !ddbGeneRegex.MatchString(firstCellText) {
			return
		}

		// Extract the three columns with proper text trimming
		// Based on HTML analysis: gene products in columns 3-4, descriptions in columns 6-7
		record := GeneDataRecord{
			GeneID:      firstCellText,
			GeneProduct: extractTextFromColumns(cells, 3, 4),
			Description: extractTextFromColumns(cells, 6, 7),
		}

		records = append(records, record)
	})

	return records, nil
}

// extractTextFromColumns searches for non-empty text in the specified column range
func extractTextFromColumns(
	cells *goquery.Selection,
	startCol, endCol int,
) string {
	for i := startCol; i <= endCol && i < cells.Length(); i++ {
		text := strings.TrimSpace(cells.Eq(i).Text())
		if text != "" {
			return text
		}
	}
	return ""
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
		if record.GeneProduct == "" {
			continue
		}
		if err := writer.Write([]string{
			record.GeneID,
			record.GeneProduct,
		}); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
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
		if record.Description == "" {
			continue
		}
		if err := writer.Write([]string{
			record.GeneID,
			record.Description,
		}); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
