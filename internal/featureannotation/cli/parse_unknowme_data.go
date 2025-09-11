package cli

import (
	"encoding/csv"
	"fmt"
	"iter"
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

	// Create iterator for gene products CSV
	geneProductIter, err := parseHTMLTableIter(inputFile)
	if err != nil {
		return fmt.Errorf(
			"failed to create iterator for gene products: %w",
			err,
		)
	}

	// Count records while generating gene product CSV for reporting
	var totalRecords int
	countingIter := func(yield func(GeneDataRecord) bool) {
		for record := range geneProductIter {
			totalRecords++
			if !yield(record) {
				return
			}
		}
	}

	// Generate gene product CSV
	err = generateGeneProductCSVFromIter(countingIter, geneProductOutput)
	if err != nil {
		return fmt.Errorf("failed to generate gene product CSV: %w", err)
	}

	// Create a second iterator for gene descriptions CSV
	geneDescriptionIter, err := parseHTMLTableIter(inputFile)
	if err != nil {
		return fmt.Errorf(
			"failed to create iterator for gene descriptions: %w",
			err,
		)
	}

	// Generate gene description CSV
	err = generateGeneDescriptionCSVFromIter(
		geneDescriptionIter,
		geneDescriptionOutput,
	)
	if err != nil {
		return fmt.Errorf("failed to generate gene description CSV: %w", err)
	}

	// Validate we found records
	if totalRecords == 0 {
		return fmt.Errorf("no gene data records found in the HTML file")
	}

	fmt.Printf("Successfully processed %d gene records\n", totalRecords)
	fmt.Printf("Gene products written to: %s\n", geneProductOutput)
	fmt.Printf("Gene descriptions written to: %s\n", geneDescriptionOutput)

	return nil
}

// parseHTMLTableIter extracts gene data from an HTML table using an iterator pattern.
// Returns an iter.Seq that yields GeneDataRecord items lazily for memory-efficient processing.
func parseHTMLTableIter(filename string) (iter.Seq[GeneDataRecord], error) {
	doc, err := loadHTMLDocument(filename)
	if err != nil {
		return nil, err
	}

	ddbGeneRegex := regexp.MustCompile(`^DDB_G\d+`)

	// Return an iterator function that yields GeneDataRecord items
	return func(yield func(GeneDataRecord) bool) {
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

			// Yield the record and check if iteration should continue
			if !yield(record) {
				return // Early termination requested by consumer
			}
		})
	}, nil
}

// extractTextFromColumns searches for non-empty text in the specified column range
// It first checks for h2 elements within td cells, then falls back to direct text content
func extractTextFromColumns(
	cells *goquery.Selection,
	startCol, endCol int,
) string {
	for i := startCol; i <= endCol && i < cells.Length(); i++ {
		cell := cells.Eq(i)

		// First try to extract text from h2 elements within the cell
		h2Text := strings.TrimSpace(cell.Find("h2").Text())
		if h2Text != "" {
			return h2Text
		}

		// Fall back to direct text content of the cell
		text := strings.TrimSpace(cell.Text())
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

// generateGeneProductCSVFromIter creates a CSV file with GeneID and gene_product columns using an iterator
func generateGeneProductCSVFromIter(
	iter iter.Seq[GeneDataRecord],
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
	err = writer.Write([]string{"GeneID", "gene_product"})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data rows using iterator
	for record := range iter {
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

// generateGeneDescriptionCSVFromIter creates a CSV file with GeneID and gene_description columns using an iterator
func generateGeneDescriptionCSVFromIter(
	iter iter.Seq[GeneDataRecord],
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

	// Write data rows using iterator
	for record := range iter {
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
