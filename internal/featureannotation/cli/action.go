package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
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

// processGeneEntryParams holds the parameters for the processGeneEntry function.
type processGeneEntryParams struct {
	entry  *Gene
	dbh    *arangomanager.Database
	client feature.FeatureAnnotationServiceClient
	logger *logrus.Entry
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

	for result.Scan() {
		entry := &Gene{}
		if err := result.Read(&entry); err != nil {
			return cli.Exit(
				fmt.Sprintf("error reading query result: %s", err),
				2,
			)
		}
		// Call helper function to process the entry
		params := &processGeneEntryParams{
			entry:  entry,
			dbh:    dbh,
			client: client,
			logger: logger,
		}
		if err := processGeneEntry(params); err != nil {
			// Exit if the helper function encounters an error
			return cli.Exit(err.Error(), 2)
		}
	}

	return nil
}

// fetchPubmedIDs queries and retrieves PubMed IDs associated with a given gene feature.
func fetchPubmedIDs(
	entry *Gene,
	dbh *arangomanager.Database,
	logger *logrus.Entry,
) ([]string, error) {
	pubmedResult, err := dbh.SearchRows(
		ListPubmedsByFeature,
		map[string]interface{}{
			"feature_id": entry.FeatureID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error querying PubMed IDs for feature %s: %w",
			entry.FeatureID,
			err,
		)
	}
	defer pubmedResult.Close()

	if pubmedResult.IsEmpty() {
		logger.Infof(
			"No PubMed references found for feature %s",
			entry.FeatureID,
		)
		return nil, nil // No IDs found, not an error
	}

	pubmedIDs := make([]string, 0)
	for pubmedResult.Scan() {
		var pubmed string
		if err := pubmedResult.Read(&pubmed); err != nil {
			return nil, fmt.Errorf(
				"error reading PubMed ID for feature %s: %w",
				entry.FeatureID,
				err,
			)
		}
		pubmedIDs = append(pubmedIDs, pubmed)
	}
	logger.Infof("Feature %s has PubMed reference: %s",
		entry.FeatureID,
		pubmedIDs,
	)
	return pubmedIDs, nil
}

// processGeneEntry handles fetching PubMed IDs and creating the annotation for a single gene entry.
func processGeneEntry(params *processGeneEntryParams) error {
	logger := params.logger
	entry := params.entry
	dbh := params.dbh
	client := params.client

	logger.Debugf("Feature has geneid %s and name %s",
		entry.GeneID,
		entry.Name,
	)

	// Fetch PubMed IDs
	pubmedIDs, err := fetchPubmedIDs(entry, dbh, logger)
	if err != nil {
		return err // Propagate error from fetching IDs
	}

	// If no PubMed IDs were found, skip creating the annotation for this entry
	if len(pubmedIDs) == 0 {
		return nil
	}

	createdBy := DefaultUserName
	if val, ok := annMap[entry.CreatedBy]; ok {
		createdBy = val
	}

	// Set up gRPC call
	res, err := client.CreateFeatureAnnotation(
		context.Background(),
		&feature.NewFeatureAnnotation{
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
		})
	if err != nil {
		return fmt.Errorf(
			"failed to create feature annotation for gene %s: %w",
			entry.GeneID,
			err,
		)
	}

	logger.Infof(
		"Created new feature annotation record %s for feature name %s",
		res.Attributes.Name,
		res.Id,
	)
	return nil
}
