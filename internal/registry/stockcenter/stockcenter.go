package stockcenter

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/spf13/viper"
)

var sv = viper.New()

func SetOrderAPIClient(oc order.OrderServiceClient) {
	sv.Set(OrderClientKey, oc)
}

func SetStockAPIClient(sc stock.StockServiceClient) {
	sv.Set(StockClientKey, sc)
}

func SetAnnotationAPIClient(ac annotation.TaggedAnnotationServiceClient) {
	sv.Set(AnnotationClientKey, ac)
}

func GetOrderAPIClient() order.OrderServiceClient {
	oc, _ := sv.Get(OrderClientKey).(order.OrderServiceClient)
	return oc
}

func GetStockAPIClient() stock.StockServiceClient {
	sc, _ := sv.Get(StockClientKey).(stock.StockServiceClient)
	return sc
}

func GetAnnotationAPIClient() annotation.TaggedAnnotationServiceClient {
	ac, _ := sv.Get(AnnotationClientKey).(annotation.TaggedAnnotationServiceClient)
	return ac
}
