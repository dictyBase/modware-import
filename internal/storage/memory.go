package storage

import (
	"fmt"
	"strings"
	"sync"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MemoryStorage implements FeatureAnnotationStorage using in-memory maps
type MemoryStorage struct {
	mu sync.RWMutex

	// Primary storage
	annotations map[string]*feature.FeatureAnnotation

	// Indexes for efficient lookups
	nameIndex   map[string]string   // name -> id
	pubmedIndex map[string][]string // pubmed_id -> []id
	doiIndex    map[string][]string // doi -> []id

	logger *logrus.Logger
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage(logger *logrus.Logger) *MemoryStorage {
	return &MemoryStorage{
		annotations: make(map[string]*feature.FeatureAnnotation),
		nameIndex:   make(map[string]string),
		pubmedIndex: make(map[string][]string),
		doiIndex:    make(map[string][]string),
		logger:      logger,
	}
}

// Create stores a new feature annotation
func (m *MemoryStorage) Create(annotation *feature.FeatureAnnotation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if annotation.Id == "" {
		return status.Error(codes.InvalidArgument, "annotation ID is required")
	}

	// Check if annotation already exists
	if _, exists := m.annotations[annotation.Id]; exists {
		return status.Error(
			codes.AlreadyExists,
			fmt.Sprintf("annotation with ID %s already exists", annotation.Id),
		)
	}

	// Store annotation
	m.annotations[annotation.Id] = annotation

	// Update indexes
	m.updateIndexes(annotation)

	m.logger.WithField("id", annotation.Id).Debug("Created feature annotation")
	return nil
}

// GetByID retrieves a feature annotation by its ID
func (m *MemoryStorage) GetByID(id string) (*feature.FeatureAnnotation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	// Check if obsolete
	if annotation.IsObsolete {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s is obsolete", id),
		)
	}

	return annotation, nil
}

// GetByName retrieves a feature annotation by its name
func (m *MemoryStorage) GetByName(
	name string,
) (*feature.FeatureAnnotation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, exists := m.nameIndex[name]
	if !exists {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with name %s not found", name),
		)
	}

	annotation := m.annotations[id]
	if annotation.IsObsolete {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with name %s is obsolete", name),
		)
	}

	return annotation, nil
}

// Update modifies an existing feature annotation
func (m *MemoryStorage) Update(
	id string,
	annotation *feature.FeatureAnnotation,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	// Remove old indexes
	m.removeIndexes(existing)

	// Update annotation
	annotation.Id = id // Ensure ID stays the same
	m.annotations[id] = annotation

	// Update indexes
	m.updateIndexes(annotation)

	m.logger.WithField("id", id).Debug("Updated feature annotation")
	return nil
}

// Delete removes a feature annotation
func (m *MemoryStorage) Delete(id string, purge bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if purge {
		// Completely remove annotation and indexes
		m.removeIndexes(annotation)
		delete(m.annotations, id)
		m.logger.WithField("id", id).Debug("Purged feature annotation")
	} else {
		// Soft delete - mark as obsolete
		annotation.IsObsolete = true
		m.logger.WithField("id", id).Debug("Soft deleted feature annotation")
	}

	return nil
}

// AddTag adds a tag property to a feature annotation
func (m *MemoryStorage) AddTag(id string, tag *feature.TagProperty) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}

	// Check if tag already exists
	for _, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag == tag.Tag {
			return status.Error(
				codes.AlreadyExists,
				fmt.Sprintf("tag %s already exists", tag.Tag),
			)
		}
	}

	// Add tag
	annotation.Attributes.Properties = append(
		annotation.Attributes.Properties,
		tag,
	)

	m.logger.WithFields(logrus.Fields{
		"id":  id,
		"tag": tag.Tag,
	}).Debug("Added tag to feature annotation")

	return nil
}

// AddTags adds multiple tag properties to a feature annotation
func (m *MemoryStorage) AddTags(id string, tags []*feature.TagProperty) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}

	for _, tag := range tags {
		for _, existingTag := range annotation.Attributes.Properties {
			if existingTag.Tag == tag.Tag {
				return status.Error(
					codes.AlreadyExists,
					fmt.Sprintf("tag %s already exists", tag.Tag),
				)
			}
		}

		annotation.Attributes.Properties = append(
			annotation.Attributes.Properties,
			tag,
		)
	}

	m.logger.WithFields(logrus.Fields{
		"id":        id,
		"tag_count": len(tags),
	}).Debug("Added tags to feature annotation")

	return nil
}

// SetTags replaces all existing tags with the provided tags (idempotent full-update)
func (m *MemoryStorage) SetTags(id string, tags []*feature.TagProperty) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}

	annotation.Attributes.Properties = tags

	m.logger.WithFields(logrus.Fields{
		"id":        id,
		"tag_count": len(tags),
	}).Debug("Set tags for feature annotation")

	return nil
}

// UpdateTag modifies an existing tag in a feature annotation
func (m *MemoryStorage) UpdateTag(
	id string,
	tagName string,
	tag *feature.TagProperty,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("tag %s not found", tagName),
		)
	}

	// Find and update tag
	for i, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag == tagName {
			annotation.Attributes.Properties[i] = tag
			m.logger.WithFields(logrus.Fields{
				"id":  id,
				"tag": tagName,
			}).Debug("Updated tag in feature annotation")
			return nil
		}
	}

	return status.Error(
		codes.NotFound,
		fmt.Sprintf("tag %s not found", tagName),
	)
}

// RemoveTag removes a tag from a feature annotation
func (m *MemoryStorage) RemoveTag(id string, tagName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("tag %s not found", tagName),
		)
	}

	// Find and remove tag
	for i, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag == tagName {
			// Remove tag by slicing
			annotation.Attributes.Properties = append(
				annotation.Attributes.Properties[:i],
				annotation.Attributes.Properties[i+1:]...,
			)
			m.logger.WithFields(logrus.Fields{
				"id":  id,
				"tag": tagName,
			}).Debug("Removed tag from feature annotation")
			return nil
		}
	}

	return status.Error(
		codes.NotFound,
		fmt.Sprintf("tag %s not found", tagName),
	)
}

// RemoveTags removes a tag with specific tag name and value from a feature annotation
func (m *MemoryStorage) RemoveTags(id string, tag string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	annotation, exists := m.annotations[id]
	if !exists {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("annotation with ID %s not found", id),
		)
	}

	if annotation.Attributes == nil {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("tag %s with value %s not found", tag, value),
		)
	}

	// Find and remove matching tag(s) with both tag name and value
	var removed bool
	var newProperties []*feature.TagProperty
	for _, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag == tag && existingTag.Value == value {
			removed = true
			// Skip this tag (don't add to newProperties)
		} else {
			newProperties = append(newProperties, existingTag)
		}
	}

	if !removed {
		return status.Error(
			codes.NotFound,
			fmt.Sprintf("tag %s with value %s not found", tag, value),
		)
	}

	annotation.Attributes.Properties = newProperties

	m.logger.WithFields(logrus.Fields{
		"id":    id,
		"tag":   tag,
		"value": value,
	}).Debug("Removed tag from feature annotation")

	return nil
}

// ListByPubmedID retrieves feature annotations by PubMed ID
func (m *MemoryStorage) ListByPubmedID(
	pubmedId string,
) ([]*feature.FeatureAnnotation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids, exists := m.pubmedIndex[pubmedId]
	if !exists {
		return []*feature.FeatureAnnotation{}, nil
	}

	var annotations []*feature.FeatureAnnotation
	for _, id := range ids {
		if annotation, exists := m.annotations[id]; exists &&
			!annotation.IsObsolete {
			annotations = append(annotations, annotation)
		}
	}

	return annotations, nil
}

// ListByDOI retrieves feature annotations by DOI
func (m *MemoryStorage) ListByDOI(
	doi string,
) ([]*feature.FeatureAnnotation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids, exists := m.doiIndex[doi]
	if !exists {
		return []*feature.FeatureAnnotation{}, nil
	}

	var annotations []*feature.FeatureAnnotation
	for _, id := range ids {
		if annotation, exists := m.annotations[id]; exists &&
			!annotation.IsObsolete {
			annotations = append(annotations, annotation)
		}
	}

	return annotations, nil
}

// updateIndexes adds/updates all indexes for an annotation
func (m *MemoryStorage) updateIndexes(annotation *feature.FeatureAnnotation) {
	// Name index
	if annotation.Attributes != nil && annotation.Attributes.Name != "" {
		m.nameIndex[annotation.Attributes.Name] = annotation.Id
	}

	// PubMed index
	if annotation.Attributes != nil {
		for _, pubmedId := range annotation.Attributes.Pubmed {
			if !m.containsString(m.pubmedIndex[pubmedId], annotation.Id) {
				m.pubmedIndex[pubmedId] = append(
					m.pubmedIndex[pubmedId],
					annotation.Id,
				)
			}
		}
	}

	// DOI index
	if annotation.Attributes != nil {
		for _, publication := range annotation.Attributes.Publications {
			if strings.HasPrefix(publication, "10.") { // Basic DOI format check
				if !m.containsString(m.doiIndex[publication], annotation.Id) {
					m.doiIndex[publication] = append(
						m.doiIndex[publication],
						annotation.Id,
					)
				}
			}
		}
	}
}

// removeIndexes removes all indexes for an annotation
func (m *MemoryStorage) removeIndexes(annotation *feature.FeatureAnnotation) {
	// Name index
	if annotation.Attributes != nil && annotation.Attributes.Name != "" {
		delete(m.nameIndex, annotation.Attributes.Name)
	}

	// PubMed index
	if annotation.Attributes != nil {
		for _, pubmedId := range annotation.Attributes.Pubmed {
			m.pubmedIndex[pubmedId] = m.removeString(
				m.pubmedIndex[pubmedId],
				annotation.Id,
			)
			if len(m.pubmedIndex[pubmedId]) == 0 {
				delete(m.pubmedIndex, pubmedId)
			}
		}
	}

	// DOI index
	if annotation.Attributes != nil {
		for _, publication := range annotation.Attributes.Publications {
			if strings.HasPrefix(publication, "10.") {
				m.doiIndex[publication] = m.removeString(
					m.doiIndex[publication],
					annotation.Id,
				)
				if len(m.doiIndex[publication]) == 0 {
					delete(m.doiIndex, publication)
				}
			}
		}
	}
}

// Helper functions
func (m *MemoryStorage) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (m *MemoryStorage) removeString(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}
