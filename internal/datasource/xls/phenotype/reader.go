// Package phenotype provides functionality to read phenotype annotations
// from an Excel file. It defines the PhenotypeAnnotationReader struct with methods
// to create a new reader, iterate over rows of data, and retrieve phenotype annotations
// as structured data with validation. The reader is initialized with a file path,
// a specific sheet name, and a creation date, and it includes error handling for
// file and row reading. It uses the third-party libraries excelize for Excel file
// manipulation and go-playground/validator for data validation.
package phenotype

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

// PhenotypeAnnotationReader is responsible for reading phenotype annotations
// from an Excel file
type PhenotypeAnnotationReader struct {
	rows          *excelize.Rows
	createdOn     time.Time
	dataValidator *validator.Validate
}

// NewPhenotypeAnnotationReader creates a new reader for phenotype annotations from an Excel file.
// It initializes the reader for the specified sheet in the file and sets the creation date for the annotations.
// The function also sets up a data validator for structural validation of the annotations.
// If the function encounters an error while opening the file or reading the rows, it returns the reader object
// created up to that point along with the error.
//
// Parameters:
// - file: The path to the Excel file to be read.
// - sheet: The name of the sheet within the Excel file which contains the phenotype annotations.
// - date: The creation date to be associated with the annotations being read.
func NewPhenotypeAnnotationReader(
	file, sheet string, date time.Time,
) (*PhenotypeAnnotationReader, error) {
	phenoReader := &PhenotypeAnnotationReader{createdOn: date}
	reader, err := excelize.OpenFile(file)
	if err != nil {
		return phenoReader, fmt.Errorf("error in reading file %s %s", file, err)
	}
	defer reader.Close()
	rows, err := reader.Rows(sheet)
	if err != nil {
		return phenoReader, fmt.Errorf("error in reading rows %s", err)
	}
	phenoReader.rows = rows
	phenoReader.dataValidator = validator.New(
		validator.WithPrivateFieldValidation(),
	)

	return phenoReader, nil
}

// Next advances the reader to the next row of phenotype annotations.
// If there are no more rows to read or an error occurs, it will close the reader and return false.
// This method should be called before each call to Value to advance the reader to the next phenotype annotation.
func (phr *PhenotypeAnnotationReader) Next() bool {
	if phr.rows.Next() {
		return true
	}
	phr.rows.Close()

	return false
}

// Value retrieves the current phenotype annotation from the reader.
// Before calling Value, Next should be used to advance the reader to the desired row.
// Value decodes the current row into a PhenotypeAnnotation struct and performs data validation.
// If the validation fails or an error occurs while reading the columns, it returns an error.
func (phr *PhenotypeAnnotationReader) Value() (*PhenotypeAnnotation, error) {
	anno := &PhenotypeAnnotation{}
	row, err := phr.rows.Columns()
	if err != nil {
		return anno, fmt.Errorf("error in reading column %s", err)
	}
	if len(row) == 0 {
		anno.empty = true
		return anno, nil
	}
	anno.strainId = row[0]
	anno.strainDescriptor = row[1]
	anno.phenotypeId = row[2]
	anno.notes = row[4]
	anno.assayId = row[5]
	anno.environmentId = row[7]
	anno.reference = row[8]
	anno.assignedBy = row[10]
	anno.deleted = false
	if err := phr.dataValidator.Struct(anno); err != nil {
		return nil, fmt.Errorf("error in data validation %s", err)
	}

	return anno, nil
}
