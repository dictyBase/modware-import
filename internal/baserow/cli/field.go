package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/client"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/sirupsen/logrus"
)

func updateFieldDefs(
	tbm *database.TableManager,
	defs map[string]map[string]interface{},
	tbl *client.Table,
	logger *logrus.Entry,
) error {
	for fieldName, spec := range defs {
		msg, err := tbm.UpdateField(tbl, fieldName, spec)
		if err != nil {
			return fmt.Errorf(
				"error in updating %s field %s",
				fieldName,
				err,
			)
		}
		logger.Debugf("%s %s", msg, fieldName)
	}

	return nil
}
