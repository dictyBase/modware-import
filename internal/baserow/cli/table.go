package cli

import (
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/urfave/cli/v2"
)

func allTableIDs(
	tbm *database.TableManager,
	flagNames []string,
	cltx *cli.Context,
) (map[string]int, error) {
	idMaps := make(map[string]int)
	for _, name := range flagNames {
		id, err := tbm.TableNameToId(cltx.String(name))
		if err != nil {
			return idMaps, err
		}
		idMaps[name] = id
	}

	return idMaps, nil
}
