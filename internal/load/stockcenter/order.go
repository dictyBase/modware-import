package stockcenter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/cobra"
)

func newPlasmidMap(r io.Reader) (*plasmidIdMap, error) {
	m := hashmap.New()
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &plasmidIdMap{}, err
		}
		m.Put(row[1], row[0])
	}
	return &plasmidIdMap{idmap: m}, nil
}

func (pm *plasmidIdMap) name2Id(name string) (string, bool) {
	v, ok := pm.idmap.Get(name)
	if !ok {
		return "", false
	}
	id, _ := v.(string)
	return id, true
}

func LoadOrder(cmd *cobra.Command, args []string) error {
	m, err := newPlasmidMap(registry.GetReader(regs.PLASMID_ID_MAP_READER))
	if err != nil {
		return fmt.Errorf("error in making plasmid map %s", err)
	}
	or := stockcenter.NewCsvStockOrderReader(registry.GetReader(regs.ORDER_READER))
	client := regs.GetOrderAPIClient()
	_, err = client.PrepareForOrder(context.Background(), &empty.Empty{})
	if err != nil {
		return fmt.Errorf("error in preparing for loading error %s", err)
	}
	logger := registry.GetLogger()
OUTER:
	for or.Next() {
		order, err := or.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading order for items %s and user %s",
				strings.Join(order.Items, " "),
				order.User,
			)
		}
		var items []string
	INNER:
		for _, e := range order.Items {
			if stRegex.MatchString(e) {
				items = append(items, e)
				continue INNER
			}
			pid, ok := m.name2Id(e)
			if !ok {
				logger.Warnf(
					"could not map plasmid %s from items %s and user %s",
					e,
					strings.Join(order.Items, " "),
					order.User,
				)
				continue OUTER
			}
			items = append(items, pid)
		}
		norder, err := client.LoadOrder(
			context.Background(),
			&pb.ExistingOrder{
				Data: &pb.ExistingOrder_Data{
					Type: "order",
					Attributes: &pb.ExistingOrderAttributes{
						CreatedAt: TimestampProto(order.CreatedAt),
						UpdatedAt: TimestampProto(order.CreatedAt),
						Purchaser: order.User,
						Items:     items,
					},
				},
			})
		if err != nil {
			return fmt.Errorf(
				"error in loading order for items %s and user %s",
				strings.Join(order.Items, " "),
				order.User,
			)
		}
		logger.Infof(
			"loaded order for items %s and user %s with id %s",
			strings.Join(norder.Data.Attributes.Items, " "),
			norder.Data.Attributes.Purchaser,
			norder.Data.Id,
		)
	}
	return nil
}
