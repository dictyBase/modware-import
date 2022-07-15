package stockcenter

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	sreg "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type parentStrain struct {
	aclient annotation.TaggedAnnotationServiceClient
	sclient pb.StockServiceClient
	logger  *logrus.Entry
}

func (p *parentStrain) isPresent(id string) (bool, error) {
	_, err := p.sclient.GetStrain(
		context.Background(),
		&pb.StockId{Id: id},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			p.logger.Infof("parent strain %s is absent", id)
			return false, nil
		}
		return false,
			fmt.Errorf("error in finding parent strain %s", err)
	}
	p.logger.Infof("parent strain %s is present", id)
	return true, nil
}

func (p *parentStrain) createAX4() error {
	_, err := p.sclient.LoadStrain(
		context.Background(),
		sreg.AX4ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX4 Parent strain %s %s", sreg.AX3ParentID, err,
		)
	}
	if err := p.insertAX4Props(); err != nil {
		return err
	}
	return p.createExtraProps(sreg.AX4ParentID)
}

func (p *parentStrain) createAX3() error {
	_, err := p.sclient.LoadStrain(
		context.Background(),
		sreg.AX3ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX3 Parent strain %s %s", sreg.AX3ParentID, err,
		)
	}
	if err := p.insertAX3Props(); err != nil {
		return err
	}
	return p.createExtraProps(sreg.AX3ParentID)
}

func (p *parentStrain) insertAX4Props() error {
	// systematic name
	err := createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyAnnoOntology,
		id:       sreg.AX4ParentID,
		tag:      sysnameTag,
		user:     sreg.DefaultUser,
		value:    sreg.AX4ParentID,
	})
	if err != nil {
		return err
	}
	// mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyMutagenesisOntology,
		id:       sreg.AX4ParentID,
		tag:      mutmethodTag,
		user:     sreg.DefaultUser,
		value:    "Spontaneous",
	})
	if err != nil {
		return err
	}
	// genetic modification
	return createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyAnnoOntology,
		id:       sreg.AX4ParentID,
		user:     sreg.DefaultUser,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func (p *parentStrain) findOrCreateAX4() error {
	ok, err := p.isPresent(sreg.AX4ParentID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return p.createAX4()
}

func (p *parentStrain) findOrCreateAX3() error {
	ok, err := p.isPresent(sreg.AX3ParentID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return p.createAX3()
}

func (p *parentStrain) insertAX3Props() error {
	// systematic name
	err := createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyAnnoOntology,
		id:       sreg.AX3ParentID,
		tag:      sysnameTag,
		user:     sreg.DefaultUser,
		value:    "AX3",
	})
	if err != nil {
		return err
	}
	// mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyMutagenesisOntology,
		id:       sreg.AX3ParentID,
		tag:      mutmethodTag,
		user:     sreg.DefaultUser,
		value:    "N-Methyl-N-Nitro-N-Nitrosoguanidine",
	})
	if err != nil {
		return err
	}
	// genetic modification
	return createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DictyAnnoOntology,
		id:       sreg.AX3ParentID,
		user:     sreg.DefaultUser,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func (p *parentStrain) createExtraProps(id string) error {
	// genotype
	_, err := NewOrReloadGeno(p.aclient, &genoArgs{
		ontology: sreg.DictyAnnoOntology,
		user:     sreg.DefaultUser,
		value:    "axeA1,axeB1,axeC1",
		tag:      genoTag,
		id:       id,
	})
	if err != nil {
		return err
	}
	// strain characteristics
	return createAnno(&createAnnoArgs{
		ontology: sreg.DictyStraincharOntology,
		user:     sreg.DefaultUser,
		client:   p.aclient,
		tag:      "axenic",
		id:       id,
		value:    val,
	})
}
