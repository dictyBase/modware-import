package stockcenter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
	"github.com/dictyBase/modware-import/internal/regexp"
)

//StrainInventory is the container for strain inventory
type StrainInventory struct {
	StrainId         string
	PrivateComment   string
	PublicComment    string
	StoredOn         time.Time
	StoredAs         string
	VialsCount       string
	VialColor        string
	PhysicalLocation string
	RecordLine       string
}

//StrainInventoryReader is the defined interface for reading the data
type StrainInventoryReader interface {
	datasource.IteratorWithoutValue
	Value() (*StrainInventory, error)
}

type tsvStrainInventoryReader struct {
	*tsource.TsvReader
}

//NewTsvStrainInventoryReader is to get an instance of StrainInventoryReader
func NewTsvStrainInventoryReader(r io.Reader) StrainInventoryReader {
	tr := bufio.NewScanner(r)
	return &tsvStrainInventoryReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new StrainInventory instance
func (sir *tsvStrainInventoryReader) Value() (*StrainInventory, error) {
	inv := new(StrainInventory)
	if sir.Err != nil {
		return inv, sir.Err
	}
	inv.StrainId = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	inv.VialsCount = sir.Record[3]
	if len(sir.Record[5]) > 0 {
		inv.StoredAs = sir.Record[5]
	}
	if len(sir.Record[6]) > 0 {
		storedOn, err := parseInvDate(sir.Record[6])
		if err != nil {
			return inv, err
		}
		inv.StoredOn = storedOn
	}
	inv.PrivateComment = sir.Record[7]
	if len(sir.Record) >= 9 {
		inv.PublicComment = sir.Record[8]
	}
	inv.RecordLine = strings.Join(sir.Record, "\t")
	return inv, nil
}

func parseInvDate(date string) (time.Time, error) {
	m := regexp.DateRegxp.FindStringSubmatch(date)
	if m == nil {
		return time.Time{}, errors.New("error in parsing date string")
	}
	return time.Parse(
		registry.STOCK_DATE_LAYOUT,
		fmt.Sprintf("%s-%s-%s", m[1], ucFirstAllLower(m[2]), m[3]),
	)
}

func ucFirstAllLower(s string) string {
	return fmt.Sprintf("%s%s", string(s[0]), strings.ToLower(s[1:]))
}
