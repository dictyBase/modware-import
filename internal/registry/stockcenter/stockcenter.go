package stockcenter

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/spf13/viper"
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
