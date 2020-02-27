package stockcenter

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type gwdiDel struct {
	aclient annotation.TaggedAnnotationServiceClient
	sclient pb.StockServiceClient
	logger  *logrus.Entry
	ctx     context.Context
}

func (gd *gwdiDel) Execute(id string) error {
	_, err := gd.sclient.RemoveStock(gd.ctx, &pb.StockId{Id: id})
	if err != nil {
		return err
	}
	gd.logger.WithFields(logrus.Fields{
		"event": "delete",
		"id":    id,
	}).Debug("remove gwdi strain")
	tac, err := gd.aclient.ListAnnotations(
		gd.ctx,
		&annotation.ListParameters{
			Limit:  20,
			Filter: fmt.Sprintf("entry_id===%s", id),
		})
	if err != nil {
		return fmt.Errorf("error in finding any gwdi annotation for %s %s", id, err)
	}
	for _, ta := range tac.Data {
		_, err := gd.aclient.DeleteAnnotation(
			gd.ctx,
			&annotation.DeleteAnnotationRequest{
				Id:    ta.Id,
				Purge: true,
			})
		if err != nil {
			return fmt.Errorf("unable to remove annotation for %s %s", id, err)
		}
	}
	gd.logger.WithFields(logrus.Fields{
		"event": "delete",
		"id":    id,
		"count": len(tac.Data),
	}).Debug("remove gwdi strain annotations")
	return nil
}

func delProducer(args *gwdiDelProdArgs) chan string {
	tasks := make(chan string)
	go func() {
		defer close(tasks)
		for _, data := range args.strains.Data {
			select {
			case <-args.ctx.Done():
				return
			case tasks <- data.Id:
			}
		}
	}()
	return tasks
}

func delConsumer(args *gwdiDelConsumerArgs) chan error {
	errc := make(chan error, 1)
	for i := 0; i < args.concurrency; i++ {
		go func(runner *gwdiDel) {
			defer close(errc)
			for {
				select {
				case <-args.ctx.Done():
					return
				case id, ok := <-args.tasks:
					if !ok {
						return
					}
					if err := runner.Execute(id); err != nil {
						errc <- err
						args.cancelFn()
						return
					}
				}
			}
		}(args.runner)
	}
	return errc
}

type gwdiCreate struct {
	aclient        annotation.TaggedAnnotationServiceClient
	sclient        pb.StockServiceClient
	logger         *logrus.Entry
	ctx            context.Context
	user           string
	value          string
	genoTag        string
	annoOntology   string
	strainCharOnto string
}

func (gc *gwdiCreate) Execute(gwdi *stockcenter.GWDIStrain) error {
	strain, err := createGwdi(gc.sclient, gwdi)
	if err != nil {
		return fmt.Errorf("error in creating new gwdi strain record  %s", err)
	}
	err = createAnno(&createAnnoArgs{
		user:     gc.user,
		id:       strain.Data.Id,
		client:   gc.aclient,
		ontology: gc.annoOntology,
		tag:      gc.genoTag,
		value:    gwdi.Genotype,
	})
	if err != nil {
		return fmt.Errorf("cannot create genotype of gwdi strain %s %s", strain.Data.Id, err)
	}
	for _, char := range gwdi.Characters {
		err = createAnno(&createAnnoArgs{
			user:     gc.user,
			id:       strain.Data.Id,
			client:   gc.aclient,
			ontology: gc.strainCharOnto,
			tag:      char,
			value:    gc.value,
		})
		if err != nil {
			return fmt.Errorf(
				"cannot create characteristic %s of gwdi strain %s %s",
				char, strain.Data.Id, err,
			)
		}
	}
	for onto, prop := range gwdi.Properties {
		err = createAnno(&createAnnoArgs{
			user:     gc.user,
			id:       strain.Data.Id,
			client:   gc.aclient,
			ontology: onto,
			tag:      prop.Property,
			value:    prop.Value,
		})
		if err != nil {
			return fmt.Errorf(
				"cannot create property %s of gwdi strain %s %s",
				prop.Property, strain.Data.Id, err,
			)
		}
	}
	gc.logger.WithFields(logrus.Fields{
		"event": "create",
		"id":    strain.Data.Id,
	}).Debug("new gwdi strain record")
	return nil
}

func createGwdi(client pb.StockServiceClient, gwdi *stockcenter.GWDIStrain) (*pb.Strain, error) {
	attr := &pb.NewStrainAttributes{
		CreatedBy:    regs.DEFAULT_USER,
		UpdatedBy:    regs.DEFAULT_USER,
		Summary:      gwdi.Summary,
		Genes:        gwdi.Genes,
		Depositor:    gwdi.Depositor,
		Label:        gwdi.Label,
		Species:      gwdi.Species,
		Plasmid:      gwdi.Plasmid,
		Parent:       gwdi.Parent,
		Publications: []string{gwdi.Publication},
		Names:        []string{gwdi.Name},
	}
	return client.CreateStrain(
		context.Background(),
		&pb.NewStrain{
			Data: &pb.NewStrain_Data{
				Type:       "strain",
				Attributes: attr,
			},
		},
	)
}

func createConsumer(args *gwdiCreateConsumerArgs) chan error {
	errc := make(chan error, 1)
	for i := 0; i < args.concurrency; i++ {
		go func(runner *gwdiCreate) {
			defer close(errc)
			for {
				select {
				case <-args.ctx.Done():
					return
				case gwdi, ok := <-args.tasks:
					if !ok {
						return
					}
					if err := runner.Execute(gwdi); err != nil {
						errc <- err
						args.cancelFn()
						return
					}
				}
			}
		}(args.runner)
	}
	return errc
}

func createProducer(args *gwdiCreateProdArgs) (chan *stockcenter.GWDIStrain, chan error) {
	tasks := make(chan *stockcenter.GWDIStrain)
	errc := make(chan error, 1)
	go func() {
		defer close(tasks)
		defer close(errc)
		for args.gr.Next() {
			gwdi, err := args.gr.Value()
			select {
			case <-args.ctx.Done():
				return
			default:
				if err != nil {
					errc <- err
					args.cancelFn()
					return
				}
				tasks <- gwdi
			}

		}
	}()
	return tasks, errc
}

func LoadGwdi(cmd *cobra.Command, args []string) error {
	gr := stockcenter.NewGWDIStrainReader(registry.GetReader(regs.GWDI_READER))
	stclient := regs.GetStockAPIClient()
	annclient := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger().WithFields(logrus.Fields{
		"type":  "gwdi",
		"stock": "strain",
	})
	count := 0
	for gr.Next() {
		gwdi, err := gr.Value()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"event": "read",
			}).Errorf("gwdi datasource error %s", err)
			continue
		}
		strain, err := createGwdi(stclient, gwdi)
		if err != nil {
			return fmt.Errorf("error in creating new gwdi strain record  %s", err)
		}
		err = createAnno(&createAnnoArgs{
			user:     regs.DEFAULT_USER,
			id:       strain.Data.Id,
			client:   annclient,
			ontology: regs.DICTY_ANNO_ONTOLOGY,
			tag:      genoTag,
			value:    gwdi.Genotype,
		})
		if err != nil {
			return fmt.Errorf("cannot create genotype of gwdi strain %s %s", strain.Data.Id, err)
		}
		for _, char := range gwdi.Characters {
			err = createAnno(&createAnnoArgs{
				user:     regs.DEFAULT_USER,
				id:       strain.Data.Id,
				client:   annclient,
				ontology: strainCharOnto,
				tag:      char,
				value:    val,
			})
			if err != nil {
				return fmt.Errorf(
					"cannot create characteristic %s of gwdi strain %s %s",
					char, strain.Data.Id, err,
				)
			}
		}
		for onto, prop := range gwdi.Properties {
			err = createAnno(&createAnnoArgs{
				user:     regs.DEFAULT_USER,
				id:       strain.Data.Id,
				client:   annclient,
				ontology: onto,
				tag:      prop.Property,
				value:    prop.Value,
			})
			if err != nil {
				return fmt.Errorf(
					"cannot create property %s of gwdi strain %s %s",
					prop.Property, strain.Data.Id, err,
				)
			}
		}
		logger.WithFields(logrus.Fields{
			"event": "create",
			"id":    strain.Data.Id,
		}).Debug("new gwdi strain record")
		count++
	}
	logger.WithFields(logrus.Fields{
		"event": "load",
		"count": count,
	}).Info("all gwdi records")
	return nil
}
