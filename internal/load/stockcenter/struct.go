package stockcenter

import (
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/sirupsen/logrus"
)

type getPhenoArgs struct {
	ontology string
	id       string
	client   pb.TaggedAnnotationServiceClient
}

type strainPhenoArgs struct {
	id         string
	client     pb.TaggedAnnotationServiceClient
	phenoSlice []*stockcenter.Phenotype
}

type processPhenoArgs struct {
	client pb.TaggedAnnotationServiceClient
	pr     stockcenter.PhenotypeReader
	logger *logrus.Entry
}

type strainInvArgs struct {
	id       string
	client   pb.TaggedAnnotationServiceClient
	invSlice []*stockcenter.StrainInventory
	found    bool
}

type plasmidInvArgs struct {
	id       string
	client   pb.TaggedAnnotationServiceClient
	invSlice []*stockcenter.PlasmidInventory
	found    bool
}

type validateTagArgs struct {
	client   pb.TaggedAnnotationServiceClient
	tag      string
	ontology string
	id       string
	stock    string
	loader   string
	logger   *logrus.Entry
}
