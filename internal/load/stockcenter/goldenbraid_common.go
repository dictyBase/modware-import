package stockcenter

import (
	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

// collectionToOption converts PlasmidCollection to Option[*Plasmid].
// Returns None when Data is empty (plasmid not found), Some(first) otherwise.
func collectionToOption(
	collection *stock.PlasmidCollection,
) O.Option[*stock.Plasmid] {
	return F.Pipe2(
		collection.Data,
		A.Head[*stock.PlasmidCollection_Data],
		O.Map(convertCollectionDataToPlasmid),
	)
}

// convertCollectionDataToPlasmid reshapes a collection item to a full Plasmid message.
func convertCollectionDataToPlasmid(
	data *stock.PlasmidCollection_Data,
) *stock.Plasmid {
	return &stock.Plasmid{
		Data: &stock.Plasmid_Data{
			Type:       data.Type,
			Id:         data.Id,
			Attributes: data.Attributes,
		},
	}
}
