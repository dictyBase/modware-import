package strain

import (
	"fmt"
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
	rdr, err := xls.NewReader(file, sheet, date)
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
	anno.descriptor = row[0]
	anno.name = row[1]
	anno.summary = row[2]
	anno.systematicName = row[3]
	anno.characteristic = row[4]
	anno.geneticModification = row[5]
	anno.mutagenesisMethod = row[6]
	anno.plasmid = row[7]
	anno.parentId = row[8]
	anno.genes = row[10]
	anno.genotype = row[11]
	anno.depositor = row[12]
	anno.species = row[13]
	anno.reference = row[14]
	anno.assignedBy = row[15]
	if err := stnr.DataValidator.Struct(anno); err != nil {
		return nil, fmt.Errorf("error in data validation %s", err)
	}

	return anno, nil
}
