package stockcenter

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrainProp(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewTsvStockPropReader(registry.GetReader(regs.StrainpropReader))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	pcount := 0
	for pr.Next() {
		prop, err := pr.Value()
		if err != nil {
			return fmt.Errorf("error in reading property for strain %s", err)
		}
		var onto string
		switch prop.Property {
		case muttypeTag:
			onto = regs.DictyAnnoOntology
		case sysnameTag:
			onto = regs.DictyAnnoOntology
		case mutmethodTag:
			onto = regs.DictyMutagenesisOntology
		case synTag:
			// it is loaded by the synonym loader
			continue
		default:
			logger.Warnf("property %s is not recognized, record is not loaded",
				prop.Property,
			)
			continue
		}
		_, err = findOrCreateAnno(&createAnnoArgs{
			client:   client,
			ontology: onto,
			id:       prop.Id,
			value:    prop.Value,
			tag:      prop.Property,
		})
		if err != nil {
			return err
		}
		logger.Debugf("loaded strain %s property with prop %s and value %s",
			prop.Id, prop.Property, prop.Value,
		)
		pcount++
	}
	logger.WithFields(logrus.Fields{
		"type":  "property",
		"stock": "strain",
		"event": "load",
		"count": pcount,
	}).Infof("loaded strain properties")
	return nil
}
