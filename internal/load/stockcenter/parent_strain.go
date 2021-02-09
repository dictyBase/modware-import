package stockcenter

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	sreg "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type parentStrain struct {
	aclient annotation.TaggedAnnotationServiceClient
	sclient pb.StockServiceClient
}

func (p *parentStrain) isPresent(id string) (bool, error) {
	_, err := p.sclient.GetStrain(
		context.Background(),
		&pb.StockId{Id: id},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false,
			fmt.Errorf("error in finding parent strain %s", err)
	}
	return true, nil
}

func (p *parentStrain) createAX4() error {
	_, err := p.sclient.LoadStrain(
		context.Background(),
		sreg.AX4ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX4 Parent strain %s %s", sreg.AX3ParentId, err,
		)
	}
	if err := p.insertAX4Props(); err != nil {
		return err
	}
	return p.createExtraProps(sreg.AX4ParentId)
}

func (p *parentStrain) createAX3() error {
	_, err := p.sclient.LoadStrain(
		context.Background(),
		sreg.AX3ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX3 Parent strain %s %s", sreg.AX3ParentId, err,
		)
	}
	if err := p.insertAX3Props(); err != nil {
		return err
	}
	return p.createExtraProps(sreg.AX3ParentId)
}

func (p *parentStrain) insertAX4Props() error {
	//systematic name
	err := createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX4ParentId,
		tag:      sysnameTag,
		user:     sreg.DEFAULT_USER,
		value:    sreg.AX4ParentId,
	})
	if err != nil {
		return err
	}
	//mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_MUTAGENESIS_ONTOLOGY,
		id:       sreg.AX4ParentId,
		tag:      mutmethodTag,
		user:     sreg.DEFAULT_USER,
		value:    "Spontaneous",
	})
	if err != nil {
		return err
	}
	//genetic modification
	return createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX4ParentId,
		user:     sreg.DEFAULT_USER,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func (p *parentStrain) findOrCreateAX4() error {
	ok, err := p.isPresent(sreg.AX4ParentId)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return p.createAX4()
}

func (p *parentStrain) findOrCreateAX3() error {
	ok, err := p.isPresent(sreg.AX3ParentId)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return p.createAX3()
}

func (p *parentStrain) insertAX3Props() error {
	//systematic name
	err := createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX3ParentId,
		tag:      sysnameTag,
		user:     sreg.DEFAULT_USER,
		value:    "AX3",
	})
	if err != nil {
		return err
	}
	//mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_MUTAGENESIS_ONTOLOGY,
		id:       sreg.AX3ParentId,
		tag:      mutmethodTag,
		user:     sreg.DEFAULT_USER,
		value:    "N-Methyl-N-Nitro-N-Nitrosoguanidine",
	})
	if err != nil {
		return err
	}
	//genetic modification
	return createAnno(&createAnnoArgs{
		client:   p.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX3ParentId,
		user:     sreg.DEFAULT_USER,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func (p *parentStrain) createExtraProps(id string) error {
	//genotype
	_, err := NewOrReloadGeno(p.aclient, &genoArgs{
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		user:     sreg.DEFAULT_USER,
		value:    "axeA1,axeB1,axeC1",
		tag:      genoTag,
		id:       id,
	})
	if err != nil {
		return err
	}
	//strain characteristics
	return createAnno(&createAnnoArgs{
		ontology: sreg.DICTY_STRAINCHAR_ONTOLOGY,
		user:     sreg.DEFAULT_USER,
		client:   p.aclient,
		tag:      "axenic",
		id:       id,
		value:    val,
	})
}
