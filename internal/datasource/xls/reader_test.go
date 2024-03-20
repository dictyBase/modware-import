package xls

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func randomInt(num int) (int, error) {
	randomValue, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		return 0, err
	}

	return int(randomValue.Int64()), nil
}

// fixedLenRandomString generates a random string of fixed length.
func fixedLenRandomString(length int) string {
	alphanum := []byte(
		"123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	)
	byt := make([]byte, 0)
	alen := len(alphanum)
	for i := 0; i < length; i++ {
		pos, _ := randomInt(alen)
		byt = append(byt, alphanum[pos])
	}

	return string(byt)
}

func createTempExcelFile(
	t *testing.T,
	sheetName string,
	createRow bool,
) (string, error) {
	var empty string
	xfh := excelize.NewFile()
	index, err := xfh.NewSheet(sheetName)
	if err != nil {
		return empty, fmt.Errorf(
			"error in creating sheet %s %s",
			sheetName,
			err,
		)
	}
	if createRow {
		xfh.SetCellValue(sheetName, "A1", "Hello")
		xfh.SetCellValue(sheetName, "A2", "No Hello for you")
	}
	xfh.SetActiveSheet(index)
	tmpFileName := filepath.Join(
		t.TempDir(),
		fmt.Sprintf("%s.xlsx", fixedLenRandomString(9)),
	)
	err = xfh.SaveAs(tmpFileName)
	if err != nil {
		return empty, fmt.Errorf(
			"error in saving excel file %s %s",
			tmpFileName,
			err,
		)
	}

	return tmpFileName, nil
}

func TestNewReader(t *testing.T) {
	t.Run("should return error if file does not exist", func(t *testing.T) {
		t.Parallel()
		_, err := NewReader(
			"non-existent-file.xlsx",
			"Sheet1",
			time.Now(),
			false,
		)
		assert.Error(t, err)
	})

	t.Run("should return error if sheet does not exist", func(t *testing.T) {
		t.Parallel()
		xlsFile, err := createTempExcelFile(t, "Sheet34", false)
		assert.NoError(t, err)
		_, err = NewReader(xlsFile, "non-existent-sheet", time.Now(), false)
		assert.Error(t, err)
	})

	t.Run(
		"should return XlsReader if file and sheet exist",
		func(t *testing.T) {
			t.Parallel()
			xlsFile, err := createTempExcelFile(t, "Sheet74", false)
			assert.NoError(t, err)
			reader, err := NewReader(xlsFile, "Sheet74", time.Now(), false)
			assert.NoError(t, err)
			assert.NotNil(t, reader)
		},
	)
}

func TestXlsReader_Next(t *testing.T) {
	t.Run("should return true if there are rows", func(t *testing.T) {
		t.Parallel()
		xlsFile, err := createTempExcelFile(t, "Sheet72", true)
		assert.NoError(t, err)
		reader, err := NewReader(xlsFile, "Sheet72", time.Now(), false)
		assert.NoError(t, err)
		assert.True(t, reader.Next())
	})
	t.Run("should return false if there are no rows", func(t *testing.T) {
		t.Parallel()
		xlsFile, err := createTempExcelFile(t, "Sheet79", false)
		assert.NoError(t, err)
		reader, err := NewReader(xlsFile, "Sheet79", time.Now(), false)
		assert.NoError(t, err)
		assert.False(t, reader.Next())
	})
	t.Run("should return false if there are no more rows", func(t *testing.T) {
		t.Parallel()
		xlsFile, err := createTempExcelFile(t, "Sheet89", true)
		assert.NoError(t, err)
		reader, err := NewReader(xlsFile, "Sheet89", time.Now(), false)
		assert.NoError(t, err)
		assert.True(t, reader.Next())
		assert.True(t, reader.Next())
		assert.False(t, reader.Next())
	})
	t.Run("should skip headers or first row", func(t *testing.T) {
		t.Parallel()
		xlsFile, err := createTempExcelFile(t, "Sheet799", true)
		assert.NoError(t, err)
		reader, err := NewReader(xlsFile, "Sheet799", time.Now(), true)
		assert.NoError(t, err)
		assert.True(t, reader.Next())
		assert.False(t, reader.Next())
	})
}
