package stockcenter

import (
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/sirupsen/logrus"
)

type createAnnoArgs struct {
	ontology string
	tag      string
	value    string
	id       string
	rank     int
	client   pb.TaggedAnnotationServiceClient
}

type createPhenoArgs struct {
	id     string
	rank   int
	pheno  *stockcenter.Phenotype
	client pb.TaggedAnnotationServiceClient
}

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

type annoParams struct {
	client pb.TaggedAnnotationServiceClient
	gc     *pb.TaggedAnnotationGroupCollection
	logger *logrus.Entry
	id     string
	loader string
	err    error
}
