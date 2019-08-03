package stockcenter

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
)

const (
	synTag       = "synonym"
	sysnameTag   = "systematic name"
	mutmethodTag = "mutagenesis method"
	muttypeTag   = "mutant type"
)

func LoadStrainProp(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewCsvStockPropReader(registry.GetReader(regs.STRAINPROP_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	pcount := 0
	for pr.Next() {
		prop, err := pr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in reading property for strain %s",
				err,
			)
		}
		var onto string
		switch prop.Property {
		case synTag:
		case sysnameTag:
		case muttypeTag:
			onto = regs.DICTY_ANNO_ONTOLOGY
		case mutmethodTag:
			onto = regs.DICTY_MUTAGENESIS_ONTOLOGY

		default:
			logger.Warnf(
				"property %s is not recognized, record is not loaded",
				prop.Property,
			)
			continue
		}
		logger.Debugf("going with onto: %s id: %s", onto, prop.Id)
		_, err = findOrCreateAnno(client, prop.Property, prop.Id, onto, prop.Value)
		if err != nil {
			return err
		}
		logger.Debugf(
			"loaded strain %s property with prop %s and value %s",
			prop.Id, prop.Property, prop.Value,
		)
		pcount++
	}
	return nil
}
