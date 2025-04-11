package cli

import (
	"context"
	"fmt"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DefaultUserName is the default creator/updater for annotations
const DefaultUserName = "dcr@dictycr.org"

var annMap = map[string]string{
	"CGM_DDB_PASC": "pgaudet@northwestern.edu",
	"CGM_DDB_PFEY": "pfey@northwestern.edu",
	"CGM_DDB_BOBD": "robert-dodson@northwestern.edu",
	"CGM_DDB_KPIL": "kpilchar@northwestern.edu",
	"CGM_DDB":      "dictybase@northwestern.edu",
}

type Gene struct {
	FeatureID string `json:"feature_id"`
	GeneID    string `json:"gene_id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

func LoadFeatureAnnotation(cltx *cli.Context) error {
	logger := registry.GetLogger()
	dbh := registry.GetArangodbConnection()
	client := registry.GetFeatureAnnotationAPIClient()

	result, err := dbh.Search(ListActiveGenesQ)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	defer result.Close()

	if result.IsEmpty() {
		logger.Error("No PubMed references found in feature annotations")
		return nil
	}

	entry := &Gene{}
	for result.Scan() {
		if err := result.Read(&entry); err != nil {
			return cli.Exit(
				fmt.Sprintf("error reading query result: %s", err),
				2,
			)
		}

		logger.Debugf("Feature has geneid %s and name %s",
			entry.GeneID,
			entry.Name,
		)

		// Query for PubMed IDs per feature
		pubmedResult, err := dbh.SearchRows(
			ListPubmedsByFeature,
			map[string]interface{}{
				"feature_id": entry.FeatureID,
			},
		)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("error querying PubMed IDs: %s", err),
				2,
			)
		}
		defer pubmedResult.Close()
		if pubmedResult.IsEmpty() {
			logger.Infof(
				"No PubMed references found for feature %s",
				entry.FeatureID,
			)
			continue
		}

		pubmedIDs := make([]string, 0)
		for pubmedResult.Scan() {
			var pubmed string
			if err := pubmedResult.Read(&pubmed); err != nil {
				return cli.Exit(
					fmt.Sprintf("error reading PubMed ID: %s", err),
					2,
				)
			}
			pubmedIDs = append(pubmedIDs, pubmed)
		}
		logger.Infof("Feature %s has PubMed reference: %s",
			entry.FeatureID,
			pubmedIDs,
		)

		createdBy := DefaultUserName
		if val, ok := annMap[entry.CreatedBy]; ok {
			createdBy = val
		}

		// Create new feature annotation record
		annotation := &feature.NewFeatureAnnotation{
			Type:       "gene",
			Id:         entry.GeneID,
			IsObsolete: false,
			CreatedBy:  createdBy,
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:   entry.Name,
				Pubmed: pubmedIDs,
			},
		}

		// Set up gRPC call with timeout
		res, err := client.CreateFeatureAnnotation(
			context.Background(),
			annotation,
		)
		if err != nil {
			return fmt.Errorf("failed to create feature annotation: %v", err)
		}

		logger.Infof(
			"Created new feature annotation record %s for feature name %s",
			res.Attributes.Name,
			res.Id,
		)
	}

	return nil
}
