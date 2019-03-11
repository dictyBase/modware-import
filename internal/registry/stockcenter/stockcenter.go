package stockcenter

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/spf13/viper"
)

const (
	ORDER_CLIENT_KEY      = "order-client"
	PLASMID_ID_MAP_READER = "plasmid-id-map-input"
	ORDER_READER          = "order-input"
)

var sv = viper.New()

func SetOrderAPIClient(oc order.OrderServiceClient) {
	sv.Set(ORDER_CLIENT_KEY, oc)
}

func GetOrderAPIClient() order.OrderServiceClient {
	oc, _ := sv.Get(ORDER_CLIENT_KEY).(order.OrderServiceClient)
	return oc
}
