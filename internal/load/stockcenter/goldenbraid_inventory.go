package stockcenter

import (
	"context"
	"fmt"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	S "github.com/IBM/fp-go/semigroup"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

type InventoryRecord struct {
	PlasmidName string
	Location    string
}

type InventoryProcessingSummary struct {
	SuccessCount int
	ErrorCount   int
	Errors       []string
}

func InventorySummarySemigroup() S.Semigroup[InventoryProcessingSummary] {
	return S.MakeSemigroup(func(a, b InventoryProcessingSummary) InventoryProcessingSummary {
		return InventoryProcessingSummary{
			SuccessCount: a.SuccessCount + b.SuccessCount,
			ErrorCount:   a.ErrorCount + b.ErrorCount,
			Errors:       append(a.Errors, b.Errors...),
		}
	})
}

// resolvePlasmidID finds the plasmid ID for a given name
func resolvePlasmidID(name string) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		client := regsc.GetStockAPIClient()
		// Filter query: name==<Name>
		filter := fmt.Sprintf("name==%s", name)
		resp, err := client.ListPlasmids(context.Background(), &stock.StockParameters{
			Filter: filter,
			Limit:  2, // We only need 1, but requesting 2 helps detect duplicates
		})
		if err != nil {
			return "", fmt.Errorf("error listing plasmids for name %s: %w", name, err)
		}

		if len(resp.Data) == 0 {
			return "", fmt.Errorf("plasmid not found: %s", name)
		}
		if len(resp.Data) > 1 {
			return "", fmt.Errorf("multiple plasmids found for name: %s", name)
		}

		return resp.Data[0].Id, nil
	})
}

// syncInventory handles the check-delete-create logic for inventory
func syncInventory(plasmidID, location string) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		client := regsc.GetAnnotationAPIClient()
		// Check if inventory exists
		gc, err := getInventory(plasmidID, client, regsc.PlasmidInvOntO)
		if err != nil {
			return "", fmt.Errorf("error checking inventory: %w", err)
		}

		// If found (collection has data), delete it
		if len(gc.Data) > 0 {
			if err := delAnnotationGroup(client, gc); err != nil {
				return "", fmt.Errorf("error deleting existing inventory: %w", err)
			}
		}

		// Create new inventory
		return createInventory(plasmidID, location, client)
	})
}

func createInventory(
	plasmidID, location string,
	client pb.TaggedAnnotationServiceClient,
) (string, error) {
	// Prepare attributes map
	attributes := map[string]string{
		regsc.InvLocationTag:    location,
		regsc.InvStorageDateTag: time.Now().Format(time.RFC3339Nano),
	}

	var annoIds []string
	for tag, value := range attributes {
		anno, err := createAnnoWithRank(&createAnnoArgs{
			ontology: regsc.PlasmidInvOntO,
			client:   client,
			id:       plasmidID,
			value:    value,
			tag:      tag,
			rank:     0,
		})
		if err != nil {
			return "", err
		}
		annoIds = append(annoIds, anno.Data.Id)
	}

	// Create annotation group
	_, err := client.CreateAnnotationGroup(
		context.Background(),
		&pb.AnnotationIdList{Ids: annoIds},
	)
	if err != nil {
		return "", err
	}

	// Mark existence
	err = createAnno(&createAnnoArgs{
		client:   client,
		tag:      regsc.PlasmidInvTag,
		id:       plasmidID,
		ontology: regsc.PlasmidInvOntO,
		value:    regsc.InvExistValue,
	})
	if err != nil {
		return "", err
	}

	return plasmidID, nil
}

func processRow(record InventoryRecord) IOE.IOEither[error, string] {
	return F.Pipe2(
		resolvePlasmidID(record.PlasmidName),
		IOE.Chain(func(id string) IOE.IOEither[error, string] {
			return syncInventory(id, record.Location)
		}),
		IOE.Map[error](func(id string) string { return id }),
	)
}

func ProcessInventory(
	config InventoryLoaderConfig,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return IOE.TryCatchError(func() (InventoryProcessingSummary, error) {
		rows, err := config.DB.Query("SELECT Name, Location FROM inventory")
		if err != nil {
			return InventoryProcessingSummary{}, fmt.Errorf("error querying inventory: %w", err)
		}
		defer rows.Close()

		summary := InventoryProcessingSummary{}

		for rows.Next() {
			var name, location string
			if err := rows.Scan(&name, &location); err != nil {
				config.Logger.Error("failed to scan row", "error", err)
				summary.ErrorCount++
				summary.Errors = append(summary.Errors, err.Error())
				continue
			}

			record := InventoryRecord{
				PlasmidName: name,
				Location:    location,
			}

			// Process record
			result := processRow(record)()
			if E.IsRight(result) {
				summary.SuccessCount++
				config.Logger.Info("processed inventory", "plasmid", name, "location", location)
			} else {
				summary.ErrorCount++
				// Extract error from Left
				E.Fold(
					func(err error) string {
						summary.Errors = append(summary.Errors, fmt.Sprintf("%s: %v", name, err))
						config.Logger.Error("failed to process inventory", "plasmid", name, "error", err)
						return ""
					},
					func(id string) string { return "" },
				)(result)
			}
		}

		return summary, nil
	})
}
