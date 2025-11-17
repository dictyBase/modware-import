package storage

import (
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
)

// FeatureAnnotationStorage defines the interface for feature annotation storage operations
type FeatureAnnotationStorage interface {
	// Create stores a new feature annotation
	Create(annotation *feature.FeatureAnnotation) error

	// GetByID retrieves a feature annotation by its ID
	GetByID(id string) (*feature.FeatureAnnotation, error)

	// GetByName retrieves a feature annotation by its name
	GetByName(name string) (*feature.FeatureAnnotation, error)

	// Update modifies an existing feature annotation
	Update(id string, annotation *feature.FeatureAnnotation) error

	// Delete removes a feature annotation (soft delete if purge=false)
	Delete(id string, purge bool) error

	// AddTag adds a tag property to a feature annotation
	AddTag(id string, tag *feature.TagProperty) error

	// AddTags adds multiple tag properties to a feature annotation
	AddTags(id string, tags []*feature.TagProperty) error

	// SetTags replaces all existing tags with the provided tags (idempotent full-update)
	SetTags(id string, tags []*feature.TagProperty) error

	// UpdateTag modifies an existing tag in a feature annotation
	UpdateTag(id string, tagName string, tag *feature.TagProperty) error

	// RemoveTag removes a tag from a feature annotation
	RemoveTag(id string, tagName string) error

	// RemoveTags removes a tag with specific tag name and value from a feature annotation
	RemoveTags(id string, tag string, value string) error

	// ListByPubmedID retrieves feature annotations by PubMed ID
	ListByPubmedID(pubmedID string) ([]*feature.FeatureAnnotation, error)

	// ListByDOI retrieves feature annotations by DOI
	ListByDOI(doi string) ([]*feature.FeatureAnnotation, error)
}
