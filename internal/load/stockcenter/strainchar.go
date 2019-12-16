package stockcenter

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	strainCharOnto = "strain_characteristics"
	val            = "novalue"
)

func LoadStrainChar(cmd *cobra.Command, args []string) error {
	scr := stockcenter.NewTsvCharacterReader(registry.GetReader(regs.STRAINCHAR_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	count := 0
	for scr.Next() {
		chs, err := scr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in reading characteristics for strain %s",
				err,
			)
		}
		_, err = findOrCreateAnno(client, chs.Character, chs.Id, strainCharOnto, val)
		if err != nil {
			return err
		}
		logger.Debugf(
			"loaded strain %s characteristics with prop %s and value %s",
			chs.Id, chs.Character, val,
		)
		count++
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "characteristic",
			"stock": "strain",
			"event": "load",
			"count": count,
		}).Infof("loaded strain characteristics")
	return nil
}
