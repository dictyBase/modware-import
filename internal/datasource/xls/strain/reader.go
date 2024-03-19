package strain

import (
	"fmt"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource/xls"
)

type StrainAnnotationReader struct {
	*xls.XlsReader
}

func NewStrainAnnotationReader(
	file, sheet string, date time.Time,
) (*StrainAnnotationReader, error) {
	strainReader := &StrainAnnotationReader{}
	rdr, err := xls.NewReader(file, sheet, date, true)
	if err != nil {
		return strainReader, err
	}
	strainReader.XlsReader = rdr
	return strainReader, nil
}

func (stnr *StrainAnnotationReader) Value() (*StrainAnnotation, error) {
	anno := &StrainAnnotation{}
	row, err := stnr.Rows.Columns()
	if err != nil {
		return anno, fmt.Errorf("error in reading column %s", err)
	}
	if len(row) == 0 {
		anno.empty = true
		return anno, nil
	}
	anno.descriptor = strings.TrimSpace(row[0])
	anno.name = strings.TrimSpace(row[1])
	anno.summary = strings.TrimSpace(row[2])
	anno.systematicName = strings.TrimSpace(row[3])
	anno.characteristic = strings.TrimSpace(row[4])
	anno.geneticModification = strings.TrimSpace(row[5])
	anno.mutagenesisMethod = strings.TrimSpace(row[6])
	anno.plasmid = strings.TrimSpace(row[7])
	anno.parentId = strings.TrimSpace(row[8])
	anno.genes = strings.TrimSpace(row[10])
	anno.genotype = strings.TrimSpace(row[11])
	anno.depositor = strings.TrimSpace(row[12])
	anno.species = strings.TrimSpace(row[13])
	anno.reference = strings.TrimSpace(row[14])
	anno.assignedBy = strings.TrimSpace(row[15])
	if err := stnr.DataValidator.Struct(anno); err != nil {
		return anno, fmt.Errorf("error in data validation %s", err)
	}

	return anno, nil
}
