package xls

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

type XlsReader struct {
	Rows          *excelize.Rows
	CreatedOn     time.Time
	DataValidator *validator.Validate
}

func NewReaderFromStream(
	fdr io.Reader, sheet string, date time.Time, skipHeader bool,
) (*XlsReader, error) {
	xlsr := &XlsReader{CreatedOn: date}
	reader, err := excelize.OpenReader(fdr)
	if err != nil {
		return xlsr, fmt.Errorf("error in reading %s", err)
	}
	rows, err := reader.Rows(sheet)
	if err != nil {
		return xlsr, fmt.Errorf("error in reading rows %s", err)
	}
	if skipHeader {
		if rows.Next() {
			_, err := rows.Columns()
			if err != nil {
				return xlsr, fmt.Errorf("error in skipping header row %s", err)
			}

		}
	}
	xlsr.Rows = rows
	xlsr.DataValidator = validator.New(
		validator.WithPrivateFieldValidation(),
	)

	return xlsr, nil
}

func NewReader(
	file, sheet string, date time.Time, skipHeader bool,
) (*XlsReader, error) {
	xlsr := &XlsReader{CreatedOn: date}
	reader, err := os.Open(file)
	if err != nil {
		return xlsr, fmt.Errorf(
			"error in reading file %s %s",
			file,
			err,
		)
	}
	return NewReaderFromStream(reader, sheet, date, skipHeader)
}

func (xlsr *XlsReader) Next() bool {
	if xlsr.Rows.Next() {
		return true
	}
	xlsr.Rows.Close()

	return false
}
