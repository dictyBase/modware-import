package stockcenter

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	sreg "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func createAX4Parent(args *parentArgs) error {
	_, err := args.sclient.GetStrain(
		context.Background(),
		&stock.StockId{Id: sreg.AX4ParentId},
	)
	if err == nil { //AX4 parent exists
		return nil
	}
	if grpc.Code(err) != codes.NotFound {
		return err
	}
	_, err = args.sclient.LoadStrain(
		context.Background(),
		sreg.AX4ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX4 Parent strain %s %s", sreg.AX4ParentId, err,
		)
	}
	if err := loadAX4ParentProps(args); err != nil {
		return err

	}
	return loadParentMoreProps(sreg.AX4ParentId, args)
}

func createAX3Parent(args *parentArgs) error {
	_, err := args.sclient.GetStrain(
		context.Background(),
		&stock.StockId{Id: sreg.AX3ParentId},
	)
	if err == nil { //AX3 parent exists
		return nil
	}
	if grpc.Code(err) != codes.NotFound {
		return err
	}
	_, err = args.sclient.LoadStrain(
		context.Background(),
		sreg.AX3ParentStrain(),
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating AX3 Parent strain %s %s", sreg.AX3ParentId, err,
		)
	}
	if err := loadAX3ParentProps(args); err != nil {
		return err

	}
	return loadParentMoreProps(sreg.AX3ParentId, args)
}

func loadAX3ParentProps(args *parentArgs) error {
	//systematic name
	err := createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX3ParentId,
		tag:      sysnameTag,
		value:    "AX3",
	})
	if err != nil {
		return err
	}
	//mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_MUTAGENESIS_ONTOLOGY,
		id:       sreg.AX3ParentId,
		tag:      mutmethodTag,
		value:    "N-Methyl-N-Nitro-N-Nitrosoguanidine",
	})
	if err != nil {
		return err
	}
	//genetic modification
	return createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX3ParentId,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func loadAX4ParentProps(args *parentArgs) error {
	//systematic name
	err := createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX4ParentId,
		tag:      sysnameTag,
		value:    sreg.AX4ParentId,
	})
	if err != nil {
		return err
	}
	//mutagenesis method
	err = createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_MUTAGENESIS_ONTOLOGY,
		id:       sreg.AX4ParentId,
		tag:      mutmethodTag,
		value:    "Spontaneous",
	})
	if err != nil {
		return err
	}
	//genetic modification
	return createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		id:       sreg.AX4ParentId,
		tag:      muttypeTag,
		value:    "endogenous mutation",
	})
}

func loadParentMoreProps(id string, args *parentArgs) error {
	//genotype
	_, err := NewOrReloadGeno(args.aclient, &genoArgs{
		ontology: sreg.DICTY_ANNO_ONTOLOGY,
		user:     sreg.DEFAULT_USER,
		id:       id,
		tag:      genoTag,
		value:    "axeA1,axeB1,axeC1",
	})
	if err != nil {
		return err
	}
	//strain characteristics
	return createAnno(&createAnnoArgs{
		client:   args.aclient,
		ontology: sreg.DICTY_STRAINCHAR_ONTOLOGY,
		id:       id,
		tag:      "axenic",
		value:    val,
	})
}
