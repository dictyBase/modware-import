package phenotype

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

type PhenotypeAnnotationReader struct {
	rows          *excelize.Rows
	createdOn     time.Time
	dataValidator *validator.Validate
}

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

func (phr *PhenotypeAnnotationReader) Next() bool {
	if phr.rows.Next() {
		return true
	}
	phr.rows.Close()

	return false
}

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
