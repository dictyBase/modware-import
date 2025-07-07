package server

import (
	"context"
	"fmt"
	"regexp"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
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
	validator := validator.New()

	// Register custom validation for DOI format
	validator.RegisterValidation("doi", validateDOI)

	return &FeatureAnnotationServer{
		storage:   storage,
		logger:    logger,
		validator: validator,
	}
}

// CreateFeatureAnnotation creates a new feature annotation
func (s *FeatureAnnotationServer) CreateFeatureAnnotation(
	ctx context.Context,
	req *feature.NewFeatureAnnotation,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithField("id", req.Id).Debug("CreateFeatureAnnotation called")

	// Validate request
	if err := s.validateNewFeatureAnnotation(req); err != nil {
		return nil, err
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
		s.logger.WithError(err).WithField("id", req.Id).Error("Failed to create annotation")
		return nil, err
	}

	s.logger.WithField("id", req.Id).Info("Created feature annotation")
	return annotation, nil
}

// GetFeatureAnnotation retrieves a feature annotation by ID
func (s *FeatureAnnotationServer) GetFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationId,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithField("id", req.Id).Debug("GetFeatureAnnotation called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		s.logger.WithError(err).WithField("id", req.Id).Error("Failed to get annotation")
		return nil, err
	}

	return annotation, nil
}

// GetFeatureAnnotationByName is not implemented in the current protobuf interface
// This method would retrieve a feature annotation by name if it was in the interface

// UpdateFeatureAnnotation updates an existing feature annotation
func (s *FeatureAnnotationServer) UpdateFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationUpdate,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithField("id", req.Id).Debug("UpdateFeatureAnnotation called")

	// Validate request
	if err := s.validateFeatureAnnotationUpdate(req); err != nil {
		return nil, err
	}

	// Get existing annotation
	existing, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Create updated annotation
	now := timestamppb.New(time.Now())
	updated := &feature.FeatureAnnotation{
		Type:       req.Type,
		Id:         req.Id,
		Attributes: req.Attributes,
		CreatedBy:  existing.CreatedBy, // Preserve original creator
		UpdatedBy:  req.UpdatedBy,
		CreatedAt:  existing.CreatedAt, // Preserve creation time
		UpdatedAt:  now,
		IsObsolete: req.IsObsolete,
	}

	// Update storage
	if err := s.storage.Update(req.Id, updated); err != nil {
		s.logger.WithError(err).WithField("id", req.Id).Error("Failed to update annotation")
		return nil, err
	}

	s.logger.WithField("id", req.Id).Info("Updated feature annotation")
	return updated, nil
}

// DeleteFeatureAnnotation deletes a feature annotation
func (s *FeatureAnnotationServer) DeleteFeatureAnnotation(
	ctx context.Context,
	req *feature.DeleteFeatureAnnotationRequest,
) (*emptypb.Empty, error) {
	s.logger.WithFields(logrus.Fields{
		"id":    req.Id,
		"purge": req.Purge,
	}).Debug("DeleteFeatureAnnotation called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if err := s.storage.Delete(req.Id, req.Purge); err != nil {
		s.logger.WithError(err).WithField("id", req.Id).Error("Failed to delete annotation")
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
	ctx context.Context,
	req *feature.AddTagRequest,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,
		"tag": req.Tag.Tag,
	}).Debug("AddTag called")

	// Validate request
	if err := s.validateAddTagRequest(req); err != nil {
		return nil, err
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
		s.logger.WithError(err).WithFields(logrus.Fields{
			"id":  req.Id,
			"tag": req.Tag.Tag,
		}).Error("Failed to add tag")
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

// UpdateTag updates an existing tag in a feature annotation
func (s *FeatureAnnotationServer) UpdateTag(
	ctx context.Context,
	req *feature.UpdateTagRequest,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,
		"tag": req.Tag.Tag,
	}).Debug("UpdateTag called")

	// Validate request
	if err := s.validateUpdateTagRequest(req); err != nil {
		return nil, err
	}

	// Create TagProperty from TagPropertyUpdate
	now := timestamppb.New(time.Now())
	tag := &feature.TagProperty{
		Tag:       req.Tag.Tag,
		Value:     req.Tag.Value,
		UpdatedBy: req.Tag.UpdatedBy,
		UpdatedAt: now,
	}

	// Update tag in storage
	if err := s.storage.UpdateTag(req.Id, req.Tag.Tag, tag); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"id":  req.Id,
			"tag": req.Tag.Tag,
		}).Error("Failed to update tag")
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
	}).Info("Updated tag in feature annotation")

	return annotation, nil
}

// RemoveTag removes a tag from a feature annotation
func (s *FeatureAnnotationServer) RemoveTag(
	ctx context.Context,
	req *feature.RemoveTagRequest,
) (*feature.FeatureAnnotation, error) {
	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,
		"tag": req.Tag,
	}).Debug("RemoveTag called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "annotation ID is required")
	}
	if req.Tag == "" {
		return nil, status.Error(codes.InvalidArgument, "tag name is required")
	}

	// Remove tag from storage
	if err := s.storage.RemoveTag(req.Id, req.Tag); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"id":  req.Id,
			"tag": req.Tag,
		}).Error("Failed to remove tag")
		return nil, err
	}

	// Return updated annotation
	annotation, err := s.storage.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,
		"tag": req.Tag,
	}).Info("Removed tag from feature annotation")

	return annotation, nil
}

// ListFeatureAnnotationsByPubmedId lists feature annotations by PubMed ID
func (s *FeatureAnnotationServer) ListFeatureAnnotationsByPubmedId(
	ctx context.Context,
	req *feature.PubmedId,
) (*feature.FeatureAnnotationCollection, error) {
	s.logger.WithField("pubmed_id", req.Id).Debug("ListFeatureAnnotationsByPubmedId called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "PubMed ID is required")
	}

	annotations, err := s.storage.ListByPubmedID(req.Id)
	if err != nil {
		s.logger.WithError(err).WithField("pubmed_id", req.Id).Error("Failed to list annotations by PubMed ID")
		return nil, err
	}

	return &feature.FeatureAnnotationCollection{
		Data: annotations,
	}, nil
}

// ListFeatureAnnotationsByDOI lists feature annotations by DOI
func (s *FeatureAnnotationServer) ListFeatureAnnotationsByDOI(
	ctx context.Context,
	req *feature.DOI,
) (*feature.FeatureAnnotationCollection, error) {
	s.logger.WithField("doi", req.Id).Debug("ListFeatureAnnotationsByDOI called")

	// Validate DOI format
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "DOI is required")
	}

	if !isValidDOI(req.Id) {
		return nil, status.Error(codes.InvalidArgument, "invalid DOI format")
	}

	annotations, err := s.storage.ListByDOI(req.Id)
	if err != nil {
		s.logger.WithError(err).WithField("doi", req.Id).Error("Failed to list annotations by DOI")
		return nil, err
	}

	return &feature.FeatureAnnotationCollection{
		Data: annotations,
	}, nil
}

// Validation functions

func (s *FeatureAnnotationServer) validateNewFeatureAnnotation(req *feature.NewFeatureAnnotation) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if req.Attributes == nil {
		return status.Error(codes.InvalidArgument, "attributes are required")
	}

	if req.Attributes.Name == "" {
		return status.Error(codes.InvalidArgument, "feature name is required")
	}

	if req.CreatedBy == "" {
		return status.Error(codes.InvalidArgument, "created_by is required")
	}

	if !isValidEmail(req.CreatedBy) {
		return status.Error(codes.InvalidArgument, "created_by must be a valid email")
	}

	// Validate DOIs in publications
	for _, doi := range req.Attributes.Publications {
		if !isValidDOI(doi) {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid DOI format: %s", doi))
		}
	}

	return nil
}

func (s *FeatureAnnotationServer) validateFeatureAnnotationUpdate(req *feature.FeatureAnnotationUpdate) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if req.UpdatedBy == "" {
		return status.Error(codes.InvalidArgument, "updated_by is required")
	}

	if !isValidEmail(req.UpdatedBy) {
		return status.Error(codes.InvalidArgument, "updated_by must be a valid email")
	}

	// Validate DOIs in publications if attributes provided
	if req.Attributes != nil {
		for _, doi := range req.Attributes.Publications {
			if !isValidDOI(doi) {
				return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid DOI format: %s", doi))
			}
		}
	}

	return nil
}

func (s *FeatureAnnotationServer) validateAddTagRequest(req *feature.AddTagRequest) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if req.Tag == nil {
		return status.Error(codes.InvalidArgument, "tag is required")
	}

	if req.Tag.Tag == "" {
		return status.Error(codes.InvalidArgument, "tag name is required")
	}

	if req.Tag.Value == "" {
		return status.Error(codes.InvalidArgument, "tag value is required")
	}

	if req.Tag.CreatedBy == "" {
		return status.Error(codes.InvalidArgument, "created_by is required")
	}

	if !isValidEmail(req.Tag.CreatedBy) {
		return status.Error(codes.InvalidArgument, "created_by must be a valid email")
	}

	return nil
}

func (s *FeatureAnnotationServer) validateUpdateTagRequest(req *feature.UpdateTagRequest) error {
	if req.Id == "" {
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if req.Tag == nil {
		return status.Error(codes.InvalidArgument, "tag is required")
	}

	if req.Tag.Tag == "" {
		return status.Error(codes.InvalidArgument, "tag name is required")
	}

	if req.Tag.Value == "" {
		return status.Error(codes.InvalidArgument, "tag value is required")
	}

	if req.Tag.UpdatedBy == "" {
		return status.Error(codes.InvalidArgument, "updated_by is required")
	}

	if !isValidEmail(req.Tag.UpdatedBy) {
		return status.Error(codes.InvalidArgument, "updated_by must be a valid email")
	}

	return nil
}

// Helper validation functions

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidDOI(doi string) bool {
	doiRegex := regexp.MustCompile(`^10\.[0-9]{4,}(\.[0-9]+)*\/[-._;()\/:a-zA-Z0-9]+$`)
	return doiRegex.MatchString(doi)
}

func validateDOI(fl validator.FieldLevel) bool {
	return isValidDOI(fl.Field().String())
}
