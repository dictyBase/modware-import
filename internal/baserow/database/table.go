package database

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/sirupsen/logrus"
)

type TableManager struct {
	Logger     *logrus.Entry
	Client     *client.APIClient
	Ctx        context.Context
	Token      string
	DatabaseId int32
}

func (tbm *TableManager) CreateTable(
	table string,
) (*client.Table, error) {
	var row []interface{}
	row = append(row, []string{""})
	tbl, resp, err := tbm.Client.
		DatabaseTablesApi.
		CreateDatabaseTable(tbm.Ctx, tbm.DatabaseId).
		TableCreate(client.TableCreate{Name: table, Data: row}).
		Execute()
	if err != nil {
		return tbl, fmt.Errorf(
			"error in creating table %s %s",
			table, err,
		)
	}
	defer resp.Body.Close()

	return tbl, nil
}
