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

	// Create iterator for processing gene data
	geneDataIter, err := parseHTMLTableIter(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create iterator: %w", err)
	}

	// Create CSV writers
	geneProductWriter, err := NewGeneProductCSVWriter(geneProductOutput)
	if err != nil {
		return fmt.Errorf("failed to create gene product writer: %w", err)
	}
	defer geneProductWriter.Close()

	geneDescriptionWriter, err := NewGeneDescriptionCSVWriter(
		geneDescriptionOutput,
	)
	if err != nil {
		return fmt.Errorf("failed to create gene description writer: %w", err)
	}
	defer geneDescriptionWriter.Close()

	// Count records while processing for reporting
	var totalRecords int
	countingIter := func(yield func(GeneDataRecord) bool) {
		for record := range geneDataIter {
			totalRecords++
			if !yield(record) {
				return
			}
		}
	}

	// Process records with both writers using single iteration
	err = processRecordsWithWriters(
		countingIter,
		geneProductWriter,
		geneDescriptionWriter,
	)
	if err != nil {
		return fmt.Errorf("failed to process records: %w", err)
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
			clen := cells.Length()

			// Skip rows that don't have at least 3 cells
			if clen < 3 {
				return
			}

			// Check if first cell contains DDB_G ID
			firstCellText := strings.TrimSpace(cells.Eq(0).Text())
			if !ddbGeneRegex.MatchString(firstCellText) {
				return
			}

			// Extract gene product first to determine description scanning range
			geneProduct, geneProductCol := extractTextFromColumnsWithIndex(
				cells,
				3,
				7,
			)
			// Determine description scanning range based on gene product location
			var description string
			if geneProductCol != -1 {
				// Gene product found: start description scanning after gene product column
				description = extractTextFromColumns(
					cells,
					geneProductCol+1,
					clen-1,
				)
			} else {
				// Gene product not found: scan all remaining cells for description
				description = extractTextFromColumns(cells, 1, clen-1)
			}

			record := GeneDataRecord{
				GeneID:      firstCellText,
				GeneProduct: geneProduct,
				Description: description,
			}

			// Yield the record and check if iteration should continue
			if !yield(record) {
				return // Early termination requested by consumer
			}
		})
	}, nil
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

func skipGeneProduct(gp GeneDataRecord) bool {
	if gp.GeneProduct == "" ||
		strings.Contains(
			gp.GeneProduct,
			"no gp",
		) || strings.Contains(gp.GeneProduct, "unknown") {
		return true
	}
	return false
}

// CSVRecordWriter defines a common interface for writing gene data records to CSV
type CSVRecordWriter interface {
	WriteRecord(record GeneDataRecord) error
	ShouldSkip(record GeneDataRecord) bool
	Close() error
}

// GeneProductCSVWriter writes gene product data to CSV format
type GeneProductCSVWriter struct {
	writer *csv.Writer
	file   *os.File
}

// NewGeneProductCSVWriter creates a new CSV writer for gene products
func NewGeneProductCSVWriter(filename string) (*GeneProductCSVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	writer := csv.NewWriter(file)

	// Write header
	err = writer.Write([]string{"GeneID", "gene_product"})
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write header: %w", err)
	}

	return &GeneProductCSVWriter{
		writer: writer,
		file:   file,
	}, nil
}

func (w *GeneProductCSVWriter) WriteRecord(record GeneDataRecord) error {
	return w.writer.Write([]string{record.GeneID, record.GeneProduct})
}

func (w *GeneProductCSVWriter) ShouldSkip(record GeneDataRecord) bool {
	return skipGeneProduct(record)
}

func (w *GeneProductCSVWriter) Close() error {
	w.writer.Flush()
	return w.file.Close()
}

// GeneDescriptionCSVWriter writes gene description data to CSV format
type GeneDescriptionCSVWriter struct {
	writer *csv.Writer
	file   *os.File
}

// NewGeneDescriptionCSVWriter creates a new CSV writer for gene descriptions
func NewGeneDescriptionCSVWriter(
	filename string,
) (*GeneDescriptionCSVWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	writer := csv.NewWriter(file)

	// Write header
	err = writer.Write([]string{"GeneID", "gene_description"})
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write header: %w", err)
	}

	return &GeneDescriptionCSVWriter{
		writer: writer,
		file:   file,
	}, nil
}

func (w *GeneDescriptionCSVWriter) WriteRecord(record GeneDataRecord) error {
	return w.writer.Write([]string{record.GeneID, record.Description})
}

func (w *GeneDescriptionCSVWriter) ShouldSkip(record GeneDataRecord) bool {
	return record.Description == ""
}

func (w *GeneDescriptionCSVWriter) Close() error {
	w.writer.Flush()
	return w.file.Close()
}

// processRecordsWithWriters iterates through records once and writes to multiple CSV writers
func processRecordsWithWriters(
	iter iter.Seq[GeneDataRecord],
	writers ...CSVRecordWriter,
) error {
	for record := range iter {
		for _, writer := range writers {
			if !writer.ShouldSkip(record) {
				if err := writer.WriteRecord(record); err != nil {
					return fmt.Errorf(
						"failed to write record: %w",
						err,
					)
				}
			}
		}
	}
	return nil
}
