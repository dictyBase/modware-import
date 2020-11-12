package stockcenter

import (
	"context"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	reg "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func createAX3Parent(args *parentArgs) error {
	_, err := args.sclient.GetStrain(
		context.Background(),
		&stock.StockId{Id: reg.AX3ParentId},
	)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return err
		}
		_, err := args.sclient.LoadStrain(
			context.Background(),
			reg.AX3ParentStrain(),
		)
		if err != nil {
			return fmt.Errorf(
				"error in creating AX3 Parent strain %s %s", reg.AX3ParentId, err,
			)
		}
	}
	return nil
}
