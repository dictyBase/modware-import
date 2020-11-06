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
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gwdiDel struct {
	aclient annotation.TaggedAnnotationServiceClient
	sclient pb.StockServiceClient
	logger  *logrus.Entry
}

func (gd *gwdiDel) Execute(id string) error {
	_, err := gd.sclient.RemoveStock(context.Background(), &pb.StockId{Id: id})
	if err != nil {
		return err
	}
	gd.logger.WithFields(logrus.Fields{
		"event": "delete",
		"id":    id,
	}).Debug("remove gwdi strain")
	tac, err := gd.aclient.ListAnnotations(
		context.Background(),
		&annotation.ListParameters{
			Limit:  20,
			Filter: fmt.Sprintf("entry_id===%s", id),
		})
	if err != nil {
		return fmt.Errorf("error in finding any gwdi annotation for %s %s", id, err)
	}
	for _, ta := range tac.Data {
		_, err := gd.aclient.DeleteAnnotation(
			context.Background(),
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

func strainsForDeletion(args *gwdiStrainDelArgs) ([]string, error) {
	cursor := int64(0)
	var ids []string
	for {
		sc, err := args.client.ListStrains(
			context.Background(),
			&pb.StockParameters{
				Cursor: cursor,
				Limit:  20,
				Filter: "descriptor=~GWDI",
			})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return ids, fmt.Errorf("error in searching for gwdi strains %s", err)
			}
		}
		if sc.Meta.NextCursor == 0 {
			return ids, nil
		}
		cursor = sc.Meta.NextCursor
		for _, scData := range sc.Data {
			ids = append(ids, scData.Id)
			args.logger.WithFields(logrus.Fields{
				"event": "queue",
				"id":    scData.Id,
			}).Debug("queued gwdi strain for pruning")
		}
	}
	return ids, nil
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
	if err := gc.createPropAndChar(strain.Data.Id, gwdi); err != nil {
		return err
	}
	gc.logger.WithFields(logrus.Fields{
		"event": "create",
		"id":    strain.Data.Id}).Debug("new gwdi strain record")
	return nil
}

func (gc *gwdiCreate) createPropAndChar(id string, gwdi *stockcenter.GWDIStrain) error {
	for _, char := range gwdi.Characters {
		err := createAnno(&createAnnoArgs{
			user:     gc.user,
			id:       id,
			client:   gc.aclient,
			ontology: gc.strainCharOnto,
			tag:      char,
			value:    gc.value,
		})
		if err != nil {
			return fmt.Errorf("cannot create characteristic %s of gwdi strain %s %s",
				char, id, err,
			)
		}
	}
	for onto, prop := range gwdi.Properties {
		err := createAnno(&createAnnoArgs{
			user:     gc.user,
			id:       id,
			client:   gc.aclient,
			ontology: onto,
			tag:      prop.Property,
			value:    prop.Value,
		})
		if err != nil {
			return fmt.Errorf("cannot create property %s of gwdi strain %s %s",
				prop.Property, id, err,
			)
		}
	}
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
	logger := registry.GetLogger().WithFields(logrus.Fields{
		"type":  "gwdi",
		"stock": "strain",
	})
	if viper.GetBool("gwdi-prune") {
		if err := runStrainDeletion(logger); err != nil {
			return err
		}
	}
	return runConcurrentCreate(logger)
}

func runConcurrentCreate(logger *logrus.Entry) error {
	stclient := regs.GetStockAPIClient()
	annclient := regs.GetAnnotationAPIClient()
	ctx, cancelFn := context.WithCancel(context.Background())
	var errcList []<-chan error
	tasks, errc := createProducer(&gwdiCreateProdArgs{
		ctx:      ctx,
		cancelFn: cancelFn,
		gr:       stockcenter.NewGWDIStrainReader(registry.GetReader(regs.GWDI_READER)),
	})
	errcList = append(errcList, errc)
	errc = createConsumer(&gwdiCreateConsumerArgs{
		concurrency: viper.GetInt("concurrency"),
		tasks:       tasks,
		ctx:         ctx,
		cancelFn:    cancelFn,
		runner: &gwdiCreate{
			aclient:        annclient,
			sclient:        stclient,
			logger:         logger,
			ctx:            ctx,
			user:           regs.DEFAULT_USER,
			value:          val,
			genoTag:        genoTag,
			annoOntology:   regs.DICTY_ANNO_ONTOLOGY,
			strainCharOnto: strainCharOnto,
		},
	})
	errcList = append(errcList, errc)
	return waitForPipeline(errcList...)
}

func runStrainDeletion(logger *logrus.Entry) error {
	stclient := regs.GetStockAPIClient()
	annclient := regs.GetAnnotationAPIClient()
	ids, err := strainsForDeletion(&gwdiStrainDelArgs{
		client: stclient,
		logger: logger,
	})
	if err != nil {
		return err
	}
	del := &gwdiDel{
		aclient: annclient,
		sclient: stclient,
		logger:  logger,
	}
	for _, id := range ids {
		if err := del.Execute(id); err != nil {
			return err
		}
	}
	logger.Infof("deleted %d existed gwdi strains", len(ids))
	return nil
}
