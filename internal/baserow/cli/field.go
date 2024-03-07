package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func updateFieldDefs(
	tbm *database.TableManager,
	defs map[string]map[string]interface{},
	tbl *client.Table,
	logger *logrus.Entry,
) error {
	grp := &errgroup.Group{}
	for fieldName, spec := range defs {
		grp.Go(updateFieldFunc(tbm, tbl, fieldName, spec, logger))
	}

	return grp.Wait() // Wait for all goroutines to complete before returning
}

// updateFieldFunc returns a function that conforms to the signature expected by errgroup.Go
func updateFieldFunc(
	tbm *database.TableManager,
	tbl *client.Table,
	fieldName string,
	spec map[string]interface{},
	logger *logrus.Entry,
) func() error {
	return func() error {
		msg, err := tbm.UpdateField(tbl, fieldName, spec)
		if err != nil {
			return fmt.Errorf(
				"error in updating %s field %s",
				fieldName,
				err,
			)
		}
		logger.Info(msg)
		return nil
	}
}

func mergeFieldDefs(
	m1, m2 map[string]map[string]interface{},
) map[string]map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
