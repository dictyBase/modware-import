package server

import (
	"context"
	"fmt"
	"time"

	"buf.build/go/protovalidate"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FeatureAnnotationServer implements the FeatureAnnotationService gRPC interface
type FeatureAnnotationServer struct {
	feature.UnimplementedFeatureAnnotationServiceServer
	storage   storage.FeatureAnnotationStorage
	logger    *logrus.Logger
	validator *validator.Validate
}

// NewFeatureAnnotationServer creates a new FeatureAnnotationServer instance
func NewFeatureAnnotationServer(
	storage storage.FeatureAnnotationStorage,
	logger *logrus.Logger,
) *FeatureAnnotationServer {
	v := validator.New()
	return &FeatureAnnotationServer{
		storage:   storage,
		logger:    logger,
		validator: v,
	}
}

// CreateFeatureAnnotation creates a new feature annotation
func (s *FeatureAnnotationServer) CreateFeatureAnnotation(
	_ context.Context,
	req *feature.NewFeatureAnnotation,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	// Create FeatureAnnotation from NewFeatureAnnotation
	now := timestamppb.New(time.Now())
	annotation := &feature.FeatureAnnotation{
		Type:       req.Type,
		Id:         req.Id,
		Attributes: req.Attributes,
		CreatedBy:  req.CreatedBy,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: req.IsObsolete,
	}

	// Store annotation
	if err := s.storage.Create(annotation); err != nil {
		return nil, err
	}

	s.logger.WithField("id", req.Id).Info("Created feature annotation")
	return annotation, nil
}

// GetFeatureAnnotation retrieves a feature annotation by ID
func (s *FeatureAnnotationServer) GetFeatureAnnotation(
	_ context.Context,
	req *feature.FeatureAnnotationId,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return annotation, nil
}

// GetFeatureAnnotationByName retrieves a feature annotation by name
func (s *FeatureAnnotationServer) GetFeatureAnnotationByName(
	_ context.Context,
	req *feature.FeatureName,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	annotation, err := s.storage.GetByName(req.Name)
	if err != nil {
		return nil, err
	}

	return annotation, nil
}

// UpdateFeatureAnnotation updates an existing feature annotation
func (s *FeatureAnnotationServer) UpdateFeatureAnnotation(
	_ context.Context,
	req *feature.FeatureAnnotationUpdate,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	// Get existing annotation
	existing, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Determine which attributes to use
	// When both are provided, update_attributes takes precedence (as per proto spec)
	var attributes *feature.FeatureAnnotationAttributes
	switch {
	case req.UpdateAttributes != nil:
		// Handle partial updates with update_attributes
		attributes = applyPartialUpdate(existing.Attributes, req.UpdateAttributes)
	case req.Attributes != nil: //nolint:staticcheck // Intentional use of deprecated field for backward compatibility
		// Fall back to deprecated attributes field for backward compatibility
		attributes = req.Attributes //nolint:staticcheck // Intentional use of deprecated field for backward compatibility
	default:
		// Keep existing attributes if neither is provided
		attributes = existing.Attributes
	}

	// Create updated annotation
	now := timestamppb.New(time.Now())
	updated := &feature.FeatureAnnotation{
		Type:       req.Type,
		Id:         req.Id,
		Attributes: attributes,
		CreatedBy:  existing.CreatedBy, // Preserve original creator
		UpdatedBy:  req.UpdatedBy,
		CreatedAt:  existing.CreatedAt, // Preserve creation time
		UpdatedAt:  now,
		IsObsolete: req.IsObsolete,
	}

	// Update storage
	if err := s.storage.Update(req.Id, updated); err != nil {
		s.logger.WithError(err).
			WithField("id", req.Id).
			Error("Failed to update annotation")
		return nil, err
	}

	s.logger.WithField("id", req.Id).Info("Updated feature annotation")
	return updated, nil
}

// DeleteFeatureAnnotation deletes a feature annotation
func (s *FeatureAnnotationServer) DeleteFeatureAnnotation(
	_ context.Context,
	req *feature.DeleteFeatureAnnotationRequest,
) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	if err := s.storage.Delete(req.Id, req.Purge); err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":    req.Id,
		"purge": req.Purge,
	}).Info("Deleted feature annotation")

	return &emptypb.Empty{}, nil
}

// AddTag adds a tag to a feature annotation
func (s *FeatureAnnotationServer) AddTag(
	_ context.Context,
	req *feature.AddTagRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	// Create TagProperty from TagPropertyCreate
	now := timestamppb.New(time.Now())
	tag := &feature.TagProperty{
		Tag:       req.Tag.Tag,
		Value:     req.Tag.Value,
		CreatedBy: req.Tag.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Add tag to storage
	if err := s.storage.AddTag(req.Id, tag); err != nil {
		return nil, err
	}

	// Return updated annotation
	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,
		"tag": req.Tag.Tag,
	}).Info("Added tag to feature annotation")

	return annotation, nil
}

// AddTags adds multiple tags to a feature annotation
func (s *FeatureAnnotationServer) AddTags(
	_ context.Context,
	req *feature.AddTagsRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}
	tags := collection.Map(req.Tags, newTagPropertyConverter)
	if err := s.storage.AddTags(req.Id, tags); err != nil {
		return nil, err
	}
	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	s.logger.WithFields(logrus.Fields{
		"id":        req.Id,
		"tag_count": len(req.Tags),
	}).Info("Added tags to feature annotation")

	return annotation, nil
}

// SetTags replaces all existing tags with the provided tags (idempotent full-update)
func (s *FeatureAnnotationServer) SetTags(
	_ context.Context,
	req *feature.SetTagsRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}
	tags := collection.Map(req.Tags, newTagPropertyConverter)
	if err := s.storage.SetTags(req.Id, tags); err != nil {
		return nil, err
	}
	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	s.logger.WithFields(logrus.Fields{
		"id":        req.Id,
		"tag_count": len(req.Tags),
	}).Info("Set tags for feature annotation")

	return annotation, nil
}

// RemoveTags removes a tag from a feature annotation
func (s *FeatureAnnotationServer) RemoveTags(
	_ context.Context,
	req *feature.RemoveTagsRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	// Remove tag from storage
	if err := s.storage.RemoveTags(req.Id, req.Tag, req.Value); err != nil {
		return nil, err
	}

	// Return updated annotation
	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":    req.Id,
		"tag":   req.Tag,
		"value": req.Value,
	}).Info("Removed tag from feature annotation")

	return annotation, nil
}

// ListFeatureAnnotationsByPubmedId lists feature annotations by PubMed ID
//
//nolint:revive,stylecheck // Method name must match protobuf-generated interface
func (s *FeatureAnnotationServer) ListFeatureAnnotationsByPubmedId(
	_ context.Context,
	req *feature.PubmedId,
) (*feature.FeatureAnnotationCollection, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
		)
	}

	annotations, err := s.storage.ListByPubmedID(req.Id)
	if err != nil {
		s.logger.WithError(err).
			WithField("pubmed_id", req.Id).
			Error("Failed to list annotations by PubMed ID")
		return nil, err
	}

	return &feature.FeatureAnnotationCollection{
		Data: annotations,
	}, nil
}

// ListFeatureAnnotationsByDOI lists feature annotations by DOI
func (s *FeatureAnnotationServer) ListFeatureAnnotationsByDOI(
	_ context.Context,
	req *feature.DOI,
) (*feature.FeatureAnnotationCollection, error) {
	// Existing protovalidate check for protobuf validation rules
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("protobuf validation failed: %v", err),
		)
	}

	annotations, err := s.storage.ListByDOI(req.Id)
	if err != nil {
		return nil, err
	}

	return &feature.FeatureAnnotationCollection{
		Data: annotations,
	}, nil
}

func newTagPropertyConverter(
	tagCreate *feature.TagPropertyCreate,
) *feature.TagProperty {
	now := timestamppb.New(time.Now())
	tag := &feature.TagProperty{
		Tag:       tagCreate.Tag,
		Value:     tagCreate.Value,
		CreatedBy: tagCreate.CreatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if tagCreate.CreatedAt != nil &&
		tagCreate.CreatedAt.CheckValid() == nil {
		tag.CreatedAt = tagCreate.CreatedAt
		tag.UpdatedAt = tagCreate.CreatedAt
	}
	return tag
}

// applyPartialUpdate merges update attributes with existing attributes
// Only non-empty fields from updateAttrs are applied to existing attributes
func applyPartialUpdate(
	existing *feature.FeatureAnnotationAttributes,
	updateAttrs *feature.FeatureAnnotationUpdateAttributes,
) *feature.FeatureAnnotationAttributes {
	// Start with a copy of existing attributes
	result := &feature.FeatureAnnotationAttributes{
		Name:         existing.Name,
		Synonyms:     existing.Synonyms,
		Publications: existing.Publications,
		Pubmed:       existing.Pubmed,
		Dblinks:      existing.Dblinks,
		Properties:   existing.Properties,
	}

	// Apply partial updates - only update fields that are provided
	if updateAttrs.Name != "" {
		result.Name = updateAttrs.Name
	}
	if len(updateAttrs.Synonyms) > 0 {
		result.Synonyms = updateAttrs.Synonyms
	}
	if len(updateAttrs.Publications) > 0 {
		result.Publications = updateAttrs.Publications
	}
	if len(updateAttrs.Pubmed) > 0 {
		result.Pubmed = updateAttrs.Pubmed
	}
	if len(updateAttrs.Dblinks) > 0 {
		result.Dblinks = updateAttrs.Dblinks
	}
	if len(updateAttrs.Properties) > 0 {
		result.Properties = updateAttrs.Properties
	}

	return result
}
