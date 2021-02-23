package stockcenter

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	stockcenter "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter/gwdi"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoadGwdi(cmd *cobra.Command, args []string) error {
	logger := registry.GetLogger().WithFields(logrus.Fields{
		"type":  "gwdi",
		"stock": "strain",
	})
	if viper.GetBool("gwdi-prune") {
		gd := &gwdiDel{
			aclient: regs.GetAnnotationAPIClient(),
			sclient: regs.GetStockAPIClient(),
			logger:  logger,
		}
		if err := gd.runStrainDeletion(); err != nil {
			return err
		}
	}
	mr, err := mutantGroups(registry.GetReader(regs.GWDI_READER))
	if err != nil {
		return err
	}
	for _, g := range mr {
		err := runConcurrentCreate(logger, g)
		if err != nil {
			return err
		}
	}
	return nil
}

func mutantGroups(r io.Reader) ([]stockcenter.GWDIMutantReader, error) {
	var mr []stockcenter.GWDIMutantReader
	gw, err := stockcenter.NewGWDI(r)
	if err != nil {
		return mr, err
	}
	if err := gw.AnnotateMutant(); err != nil {
		return mr, err
	}
	groups, err := gw.AllGroups()
	if err != nil {
		return mr, err
	}
	for _, g := range groups {
		mr = append(mr, gw.MutantReader(g))
	}
	return mr, nil
}

func runConcurrentCreate(logger *logrus.Entry, gr stockcenter.GWDIMutantReader) error {
	stclient := regs.GetStockAPIClient()
	annclient := regs.GetAnnotationAPIClient()
	ps := &parentStrain{
		aclient: annclient,
		sclient: stclient,
	}
	if err := ps.findOrCreateAX3(); err != nil {
		return err
	}
	if err := ps.findOrCreateAX4(); err != nil {
		return err
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	var errcList []<-chan error
	tasks, errc := createProducer(&gwdiCreateProdArgs{
		ctx:      ctx,
		cancelFn: cancelFn,
		gr:       gr,
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
			strainCharOnto: regs.DICTY_STRAINCHAR_ONTOLOGY,
		},
	})
	errcList = append(errcList, errc)
	return waitForPipeline(errcList...)
}

func createConsumer(args *gwdiCreateConsumerArgs) chan error {
	errc := make(chan error, 1)
	counter := make(chan int, 1)
	wg := new(sync.WaitGroup)
	wg.Add(args.concurrency)
	for i := 0; i < args.concurrency; i++ {
		go func(args *gwdiCreateConsumerArgs) {
			lc := 0
			defer func() { counter <- lc }()
			defer wg.Done()
			for {
				select {
				case <-args.ctx.Done():
					return
				case gwdi, ok := <-args.tasks:
					if !ok {
						return
					}
					if err := args.runner.execute(gwdi); err != nil {
						errc <- err
						args.cancelFn()
						return
					}
				}
				lc++
			}
		}(args)
	}
	go loadingCount(args.runner.logger, counter)
	go syncLoader(wg, counter, errc)
	return errc
}

func loadingCount(logger *logrus.Entry, counter chan int) {
	c := 0
	for v := range counter {
		c += v
	}
	logger.WithFields(logrus.Fields{
		"type":  "counter",
		"count": c,
	}).Infof("loaded gwdi strains")
}

func syncLoader(wg *sync.WaitGroup, counter chan int, errc chan error) {
	wg.Wait()
	close(counter)
	close(errc)
}

func createProducer(args *gwdiCreateProdArgs) (chan *stockcenter.GWDIStrain, chan error) {
	tasks := make(chan *stockcenter.GWDIStrain)
	errc := make(chan error, 1)
	go func(args *gwdiCreateProdArgs) {
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
	}(args)
	return tasks, errc
}

type gwdiDel struct {
	aclient annotation.TaggedAnnotationServiceClient
	sclient pb.StockServiceClient
	logger  *logrus.Entry
}

func (gd *gwdiDel) runStrainDeletion() error {
	ids, err := gd.strainsForDeletion()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := gd.execute(id); err != nil {
			return err
		}
	}
	gd.logger.Infof("deleted %d existed gwdi strains", len(ids))
	return nil
}

func (gd *gwdiDel) strainsForDeletion() ([]string, error) {
	cursor := int64(0)
	var ids []string
	for {
		sc, err := gd.sclient.ListStrains(
			context.Background(),
			&pb.StockParameters{
				Cursor: cursor,
				Limit:  20,
				Filter: "name=~GWDI_",
			})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				break
			}
			return ids, fmt.Errorf("error getting list of strains %s", err)
		}
		if sc.Meta.NextCursor == 0 {
			ids = append(ids, gd.queueIds(sc)...)
			break
		}
		cursor = sc.Meta.NextCursor
		ids = append(ids, gd.queueIds(sc)...)
	}
	return ids, nil
}

func (gd *gwdiDel) queueIds(sc *pb.StrainCollection) []string {
	var ids []string
	for _, scData := range sc.Data {
		ids = append(ids, scData.Id)
		gd.logger.WithFields(logrus.Fields{
			"event": "queue",
			"id":    scData.Id,
		}).Debug("queued gwdi strain for pruning")
	}
	return ids
}

func (gd *gwdiDel) deleteAnno(id string) error {
	tac, err := gd.aclient.ListAnnotations(
		context.Background(),
		&annotation.ListParameters{
			Limit:  20,
			Filter: fmt.Sprintf("entry_id===%s", id),
		})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
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
			if status.Code(err) == codes.NotFound {
				continue
			}
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

func (gd *gwdiDel) execute(id string) error {
	_, err := gd.sclient.RemoveStock(context.Background(), &pb.StockId{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
		if strings.Contains(err.Error(), "document not found") {
			return nil
		}
		return fmt.Errorf("error in removing gwdi strain with id %s %s", id, err)
	}
	gd.logger.WithFields(logrus.Fields{
		"event": "delete",
		"id":    id,
	}).Debug("remove gwdi strain")
	return gd.deleteAnno(id)
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

func (gc *gwdiCreate) execute(gwdi *stockcenter.GWDIStrain) error {
	strain, err := gc.createGwdi(gwdi)
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
	// create presence of inventory annotation
	return createAnno(&createAnnoArgs{
		client:   gc.aclient,
		tag:      regs.StrainInvTag,
		ontology: regs.StrainInvOnto,
		value:    regs.InvExistValue,
		user:     regs.DEFAULT_USER,
		id:       id,
	})
}

func (gc *gwdiCreate) createGwdi(gwdi *stockcenter.GWDIStrain) (*pb.Strain, error) {
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
	strain, err := gc.sclient.CreateStrain(
		context.Background(),
		&pb.NewStrain{
			Data: &pb.NewStrain_Data{
				Type:       "strain",
				Attributes: attr,
			},
		},
	)
	if err != nil {
		return strain, fmt.Errorf("error in creating gwdi strain %s", err)
	}
	return strain, nil
}
