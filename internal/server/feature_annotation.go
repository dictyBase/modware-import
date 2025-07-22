package server

import (
	"context"
	"fmt"
	"time"

	"buf.build/go/protovalidate"
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
	return &FeatureAnnotationServer{
		storage:   storage,
		logger:    logger,
		validator: validator.New(),
	}
}

// CreateFeatureAnnotation creates a new feature annotation
func (s *FeatureAnnotationServer) CreateFeatureAnnotation(
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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

// ListFeatureAnnotationsByPubmedId lists feature annotations by PubMed ID
func (s *FeatureAnnotationServer) ListFeatureAnnotationsByPubmedId(
	ctx context.Context,
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
	ctx context.Context,
	req *feature.DOI,
) (*feature.FeatureAnnotationCollection, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("validation failed: %v", err),
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
