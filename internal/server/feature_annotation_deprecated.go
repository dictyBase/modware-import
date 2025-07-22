//nolint:staticcheck // SA1019: This file contains deprecated protobuf APIs maintained for backward compatibility
package server

import (
	"context"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// This file contains deprecated gRPC methods that are maintained for backward compatibility.
// These methods use deprecated protobuf types and fields as marked in the protobuf definitions.
// New development should use the modern alternatives:
// - Instead of UpdateTag, use SetTags or AddTags
// - Instead of RemoveTag, use RemoveTags

// UpdateTag updates an existing tag in a feature annotation
// Deprecated: Use SetTags or AddTags instead
func (s *FeatureAnnotationServer) UpdateTag(
	ctx context.Context,
	req *feature.UpdateTagRequest, //nolint:staticcheck // SA1019: deprecated but maintained for backward compatibility
) (*feature.FeatureAnnotation, error) {
	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,      //nolint:staticcheck // SA1019: deprecated field access
		"tag": req.Tag.Tag, //nolint:staticcheck // SA1019: deprecated field access
	}).Debug("UpdateTag called")

	// Validate request
	if err := s.validateUpdateTagRequest(req); err != nil {
		return nil, err
	}

	// Create TagProperty from TagPropertyUpdate
	now := timestamppb.New(time.Now())
	tag := &feature.TagProperty{
		Tag:       req.Tag.Tag,       //nolint:staticcheck // SA1019: deprecated field access
		Value:     req.Tag.Value,     //nolint:staticcheck // SA1019: deprecated field access
		UpdatedBy: req.Tag.UpdatedBy, //nolint:staticcheck // SA1019: deprecated field access
		UpdatedAt: now,
	}

	// Update tag in storage
	//nolint:staticcheck // SA1019: deprecated field access
	if err := s.storage.UpdateTag(req.Id, req.Tag.Tag, tag); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"id":  req.Id,      //nolint:staticcheck // SA1019: deprecated field access
			"tag": req.Tag.Tag, //nolint:staticcheck // SA1019: deprecated field access
		}).Error("Failed to update tag")
		return nil, err
	}

	// Return updated annotation
	annotation, err := s.storage.GetByID(req.Id) //nolint:staticcheck // SA1019: deprecated field access
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,      //nolint:staticcheck // SA1019: deprecated field access
		"tag": req.Tag.Tag, //nolint:staticcheck // SA1019: deprecated field access
	}).Info("Updated tag in feature annotation")

	return annotation, nil
}

// RemoveTag removes a tag from a feature annotation
// Deprecated: Use RemoveTags instead
func (s *FeatureAnnotationServer) RemoveTag(
	ctx context.Context,
	req *feature.RemoveTagRequest, //nolint:staticcheck // SA1019: deprecated but maintained for backward compatibility
) (*feature.FeatureAnnotation, error) {
	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,  //nolint:staticcheck // SA1019: deprecated field access
		"tag": req.Tag, //nolint:staticcheck // SA1019: deprecated field access
	}).Debug("RemoveTag called")

	if req.Id == "" { //nolint:staticcheck // SA1019: deprecated field access
		return nil, status.Error(codes.InvalidArgument, "annotation ID is required")
	}
	if req.Tag == "" { //nolint:staticcheck // SA1019: deprecated field access
		return nil, status.Error(codes.InvalidArgument, "tag name is required")
	}

	// Remove tag from storage
	if err := s.storage.RemoveTag(req.Id, req.Tag); err != nil { //nolint:staticcheck // SA1019: deprecated field access
		s.logger.WithError(err).WithFields(logrus.Fields{
			"id":  req.Id,  //nolint:staticcheck // SA1019: deprecated field access
			"tag": req.Tag, //nolint:staticcheck // SA1019: deprecated field access
		}).Error("Failed to remove tag")
		return nil, err
	}

	// Return updated annotation
	annotation, err := s.storage.GetByID(req.Id) //nolint:staticcheck // SA1019: deprecated field access
	if err != nil {
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"id":  req.Id,  //nolint:staticcheck // SA1019: deprecated field access
		"tag": req.Tag, //nolint:staticcheck // SA1019: deprecated field access
	}).Info("Removed tag from feature annotation")

	return annotation, nil
}

// validateUpdateTagRequest validates the deprecated UpdateTagRequest
// Deprecated: Used only by the deprecated UpdateTag method
//
//nolint:staticcheck // SA1019: deprecated but maintained for backward compatibility
func (s *FeatureAnnotationServer) validateUpdateTagRequest(req *feature.UpdateTagRequest) error {
	if req.Id == "" { //nolint:staticcheck // SA1019: deprecated field access
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	if req.Tag == nil {
		return status.Error(codes.InvalidArgument, "tag is required")
	}

	if req.Tag.Tag == "" { //nolint:staticcheck // SA1019: deprecated field access
		return status.Error(codes.InvalidArgument, "tag name is required")
	}

	if req.Tag.Value == "" { //nolint:staticcheck // SA1019: deprecated field access
		return status.Error(codes.InvalidArgument, "tag value is required")
	}

	if req.Tag.UpdatedBy == "" { //nolint:staticcheck // SA1019: deprecated field access
		return status.Error(codes.InvalidArgument, "updated_by is required")
	}

	if !isValidEmail(req.Tag.UpdatedBy) { //nolint:staticcheck // SA1019: deprecated field access
		return status.Error(codes.InvalidArgument, "updated_by must be a valid email")
	}

	return nil
}
