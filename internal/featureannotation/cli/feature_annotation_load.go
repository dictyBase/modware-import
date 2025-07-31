package cli

import (
	"context"
	"fmt"

	"github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// processGeneEntryParams holds the parameters for the processGeneEntry
// function.
type processGeneEntryParams struct {
	entry  *Gene
	dbh    *arangomanager.Database
	client feature.FeatureAnnotationServiceClient
	logger *logrus.Entry
}

type Gene struct {
	FeatureID int    `json:"feature_id"`
	GeneID    string `json:"gene_id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
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
			"error querying PubMed IDs for feature %d: %w",
			entry.FeatureID,
			err,
		)
	}
	defer pubmedResult.Close()

	if pubmedResult.IsEmpty() {
		logger.Infof(
			"No PubMed references found for feature %d",
			entry.FeatureID,
		)
		return []string{}, nil // Return empty slice instead of nil
	}

	pubmedIDs := make([]string, 0)
	for pubmedResult.Scan() {
		var pubmed string
		if err := pubmedResult.Read(&pubmed); err != nil {
			return nil, fmt.Errorf(
				"error reading PubMed ID for feature %d: %w",
				entry.FeatureID,
				err,
			)
		}
		pubmedIDs = append(pubmedIDs, pubmed)
	}
	logger.Infof("Feature %d has PubMed reference: %v",
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

	logger.Debugf("Feature has geneid %s", entry.GeneID)

	// Fetch PubMed IDs
	pubmedIDs, err := fetchPubmedIDs(entry, dbh, logger)
	if err != nil {
		return err // Propagate error from fetching IDs
	}

	// If no PubMed IDs were found or pubmedIDs is nil, skip creating the annotation for this entry
	if len(pubmedIDs) == 0 {
		logger.Infof("Skipping feature %s with no PubMed IDs", entry.GeneID)
		return nil
	}

	createdBy := resolveCreator(entry)

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

func resolveCreator(entry *Gene) string {
	createdBy := DefaultUserName
	if val, ok := AnnMap[entry.CreatedBy]; ok {
		createdBy = val
	}
	return createdBy
}
