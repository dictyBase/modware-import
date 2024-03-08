package xls

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

type XlsReader struct {
	Rows          *excelize.Rows
	CreatedOn     time.Time
	DataValidator *validator.Validate
}

func NewReader(
	file, sheet string, date time.Time,
) (*XlsReader, error) {
	xlsr := &XlsReader{CreatedOn: date}
	reader, err := excelize.OpenFile(file)
	if err != nil {
		return xlsr, fmt.Errorf(
			"error in reading file %s %s",
			file,
			err,
		)
	}
	defer reader.Close()
	rows, err := reader.Rows(sheet)
	if err != nil {
		return xlsr, fmt.Errorf("error in reading rows %s", err)
	}
	xlsr.Rows = rows
	xlsr.DataValidator = validator.New(
		validator.WithPrivateFieldValidation(),
	)

	return xlsr, nil
}

func (xlsr *XlsReader) Next() bool {
	if xlsr.Rows.Next() {
		return true
	}
	xlsr.Rows.Close()

	return false
}
