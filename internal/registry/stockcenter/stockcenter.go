package stockcenter

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/spf13/viper"
)

const (
	ORDER_CLIENT_KEY           = "order-client"
	STOCK_CLIENT_KEY           = "stock-client"
	ANNOTATION_CLIENT_KEY      = "annotation-client"
	PLASMID_ID_MAP_READER      = "plasmid-id-map-input"
	ORDER_READER               = "order-input"
	STRAIN_READER              = "strain-input"
	STRAIN_ANNOTATOR_READER    = "strain-annotator-input"
	STRAIN_PUB_READER          = "strain-pub-input"
	STRAIN_GENE_READER         = "strain-gene-input"
	STRAIN_CHAR_READER         = "strain-char-input"
	PHENO_READER               = "pheno-input"
	GENO_READER                = "geno-input"
	INV_READER                 = "inv-input"
	STRAINPROP_READER          = "strainprop-input"
	DEFAULT_USER               = "dictybase@northwestern.edu"
	DICTY_ANNO_ONTOLOGY        = "dicty_annotation"
	DICTY_MUTAGENESIS_ONTOLOGY = "Dd Mutagenesis Method"
	DICTY_GENETICMOD_ONTOLOGY  = "genetic modification"
)

var sv = viper.New()

func SetOrderAPIClient(oc order.OrderServiceClient) {
	sv.Set(ORDER_CLIENT_KEY, oc)
}

func SetStockAPIClient(sc stock.StockServiceClient) {
	sv.Set(STOCK_CLIENT_KEY, sc)
}

func SetAnnotationAPIClient(ac annotation.TaggedAnnotationServiceClient) {
	sv.Set(ANNOTATION_CLIENT_KEY, ac)
}

func GetOrderAPIClient() order.OrderServiceClient {
	oc, _ := sv.Get(ORDER_CLIENT_KEY).(order.OrderServiceClient)
	return oc
}

func GetStockAPIClient() stock.StockServiceClient {
	sc, _ := sv.Get(STOCK_CLIENT_KEY).(stock.StockServiceClient)
	return sc
}

func GetAnnotationAPIClient() annotation.TaggedAnnotationServiceClient {
	ac, _ := sv.Get(ANNOTATION_CLIENT_KEY).(annotation.TaggedAnnotationServiceClient)
	return ac
}
