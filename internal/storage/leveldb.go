package storage

import (
	"fmt"
	"slices"
	"strings"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type LevelDBStorage struct {
	db        *leveldb.DB
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewLevelDBStorage(
	logger *logrus.Logger,
) (FeatureAnnotationStorage, error) {
	// Use in-memory storage following the pattern in cache.go
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open in-memory LevelDB database: %w",
			err,
		)
	}
	leveldbStorage := &LevelDBStorage{
		db:        db,
		logger:    logger,
		validator: validator.New(),
	}
	return leveldbStorage, nil
}

func (s *LevelDBStorage) Close() error {
	return s.db.Close()
}

func (s *LevelDBStorage) Create(annotation *feature.FeatureAnnotation) error {
	key := fmt.Sprintf("annotation:%s", annotation.Id)
	exists, err := s.db.Has([]byte(key), nil)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to check if annotation exists: %v",
			err,
		)
	}
	if exists {
		return status.Errorf(
			codes.AlreadyExists,
			"annotation with ID %s already exists",
			annotation.Id,
		)
	}
	data, err := proto.Marshal(annotation)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to serialize annotation: %v",
			err,
		)
	}
	txn, err := s.db.OpenTransaction()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to start transaction: %v",
			err,
		)
	}
	if err := txn.Put([]byte(key), data, nil); err != nil {
		txn.Discard()
		return status.Errorf(
			codes.Internal,
			"failed to store annotation: %v",
			err,
		)
	}
	// Update indexes
	if err := s.updateIndexes(txn, annotation, true); err != nil {
		txn.Discard()
		return status.Errorf(
			codes.Internal,
			"failed to update indexes: %v",
			err,
		)
	}
	if err := txn.Commit(); err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to commit transaction: %v",
			err,
		)
	}
	s.logger.WithField("annotation_id", annotation.Id).
		Debug("Created annotation in LevelDB")
	return nil
}

func (s *LevelDBStorage) GetByID(
	id string,
) (*feature.FeatureAnnotation, error) {
	if err := s.validator.Var(id, "required"); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	key := fmt.Sprintf("annotation:%s", id)
	data, err := s.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, status.Errorf(
				codes.NotFound,
				"annotation with ID %s not found",
				id,
			)
		}
		return nil, status.Errorf(
			codes.Internal,
			"failed to retrieve annotation: %v",
			err,
		)
	}

	annotation := &feature.FeatureAnnotation{}
	if err := proto.Unmarshal(data, annotation); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to deserialize annotation: %v",
			err,
		)
	}

	if annotation.IsObsolete {
		return nil, status.Errorf(
			codes.NotFound,
			"annotation with ID %s is obsolete",
			id,
		)
	}

	return annotation, nil
}

func (s *LevelDBStorage) GetByName(
	name string,
) (*feature.FeatureAnnotation, error) {
	if err := s.validator.Var(name, "required"); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid name: %v", err),
		)
	}
	indexKey := fmt.Sprintf("name_index:%s", name)
	idData, err := s.db.Get([]byte(indexKey), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, status.Errorf(
				codes.NotFound,
				"annotation with name %s not found",
				name,
			)
		}
		return nil, status.Errorf(
			codes.Internal,
			"failed to lookup name index: %v",
			err,
		)
	}

	id := string(idData)
	return s.getByIDInternal(id)
}

func (s *LevelDBStorage) Update(
	id string,
	annotation *feature.FeatureAnnotation,
) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	existing, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	annotation.Id = id
	data, err := proto.Marshal(annotation)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to serialize annotation: %v",
			err,
		)
	}

	// Start transaction for atomic operations
	txn, err := s.db.OpenTransaction()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to start transaction: %v",
			err,
		)
	}

	// Store updated annotation
	key := fmt.Sprintf("annotation:%s", id)
	if err := txn.Put([]byte(key), data, nil); err != nil {
		txn.Discard()
		return status.Errorf(
			codes.Internal,
			"failed to update annotation: %v",
			err,
		)
	}

	// Update indexes (remove old, add new)
	if err := s.updateIndexes(txn, existing, false); err != nil {
		txn.Discard()
		return status.Errorf(
			codes.Internal,
			"failed to remove old indexes: %v",
			err,
		)
	}
	if err := s.updateIndexes(txn, annotation, true); err != nil {
		txn.Discard()
		return status.Errorf(
			codes.Internal,
			"failed to update indexes: %v",
			err,
		)
	}

	// Commit transaction
	if err := txn.Commit(); err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to commit transaction: %v",
			err,
		)
	}

	s.logger.WithField("annotation_id", id).
		Debug("Updated annotation in LevelDB")
	return nil
}

func (s *LevelDBStorage) Delete(id string, purge bool) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	// Get annotation (outside transaction)
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("annotation:%s", id)
	if purge {
		// Hard delete - use transaction for atomicity
		txn, err := s.db.OpenTransaction()
		if err != nil {
			return status.Errorf(
				codes.Internal,
				"failed to start transaction: %v",
				err,
			)
		}

		// Delete annotation
		if err := txn.Delete([]byte(key), nil); err != nil {
			txn.Discard()
			return status.Errorf(
				codes.Internal,
				"failed to delete annotation: %v",
				err,
			)
		}

		// Remove from indexes
		if err := s.updateIndexes(txn, annotation, false); err != nil {
			txn.Discard()
			return status.Errorf(
				codes.Internal,
				"failed to remove indexes: %v",
				err,
			)
		}

		// Commit transaction
		if err := txn.Commit(); err != nil {
			return status.Errorf(
				codes.Internal,
				"failed to commit transaction: %v",
				err,
			)
		}

		s.logger.WithField("annotation_id", id).
			Debug("Purged annotation from LevelDB")
	} else {
		// Soft delete - simple single operation
		annotation.IsObsolete = true
		data, err := proto.Marshal(annotation)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to serialize annotation: %v", err)
		}

		if err := s.db.Put([]byte(key), data, nil); err != nil {
			return status.Errorf(codes.Internal, "failed to update annotation: %v", err)
		}

		s.logger.WithField("annotation_id", id).Debug("Marked annotation as obsolete in LevelDB")
	}

	return nil
}

func (s *LevelDBStorage) AddTag(id string, tag *feature.TagProperty) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}
	// Add tag
	annotation.Attributes.Properties = append(
		annotation.Attributes.Properties,
		tag,
	)
	// Save updated annotation (single operation, no transaction needed)
	return s.updateAnnotation(annotation)
}

func (s *LevelDBStorage) AddTags(id string, tags []*feature.TagProperty) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}
	annotation.Attributes.Properties = slices.Concat(
		annotation.Attributes.Properties,
		tags,
	)
	return s.updateAnnotation(annotation)
}

func (s *LevelDBStorage) SetTags(id string, tags []*feature.TagProperty) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}
	annotation.Attributes.Properties = tags
	return s.updateAnnotation(annotation)
}

func (s *LevelDBStorage) UpdateTag(
	id string,
	tagName string,
	tag *feature.TagProperty,
) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	if err := s.validator.Var(tagName, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid tagName: %v", err),
		)
	}
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}
	// Find and update tag
	found := false
	for i, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag == tagName {
			annotation.Attributes.Properties[i] = tag
			found = true
			break
		}
	}

	if !found {
		return status.Errorf(
			codes.NotFound,
			"tag with name %s not found",
			tagName,
		)
	}

	// Save updated annotation (single operation, no transaction needed)
	return s.updateAnnotation(annotation)
}

func (s *LevelDBStorage) RemoveTag(id string, tagName string) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	if err := s.validator.Var(tagName, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid tagName: %v", err),
		)
	}
	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}
	if annotation.Attributes == nil {
		annotation.Attributes = &feature.FeatureAnnotationAttributes{}
	}
	found := false
	newTags := make(
		[]*feature.TagProperty,
		0,
		len(annotation.Attributes.Properties),
	)
	for _, existingTag := range annotation.Attributes.Properties {
		if existingTag.Tag != tagName {
			newTags = append(newTags, existingTag)
		} else {
			found = true
		}
	}

	if !found {
		return status.Errorf(
			codes.NotFound,
			"tag with name %s not found",
			tagName,
		)
	}

	annotation.Attributes.Properties = newTags
	return s.updateAnnotation(annotation)
}

// RemoveTags removes a tag with specific tag name and value from a feature
// annotation.
func (s *LevelDBStorage) RemoveTags(id string, tag string, value string) error {
	if err := s.validator.Var(id, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid id: %v", err),
		)
	}
	if err := s.validator.Var(tag, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid tag: %v", err),
		)
	}
	if err := s.validator.Var(value, "required"); err != nil {
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid value: %v", err),
		)
	}

	annotation, err := s.getByIDInternal(id)
	if err != nil {
		return err
	}

	annotation.Attributes.Properties = slices.DeleteFunc(
		annotation.Attributes.Properties,
		func(existingTag *feature.TagProperty) bool {
			return existingTag.Tag == tag &&
				existingTag.Value == value
		},
	)

	return s.updateAnnotation(annotation)
}

func (s *LevelDBStorage) ListByPubmedID(
	pubmedId string,
) ([]*feature.FeatureAnnotation, error) {
	if err := s.validator.Var(pubmedId, "required"); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid pubmedId: %v", err),
		)
	}
	indexKey := fmt.Sprintf("pubmed_index:%s", pubmedId)
	idsData, err := s.db.Get([]byte(indexKey), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []*feature.FeatureAnnotation{}, nil
		}
		return nil, status.Errorf(
			codes.Internal,
			"failed to lookup PubMed index: %v",
			err,
		)
	}

	ids := strings.Split(string(idsData), ",")
	annotations := make([]*feature.FeatureAnnotation, 0, len(ids))

	for _, id := range ids {
		if id == "" {
			continue
		}
		annotation, err := s.getByIDInternal(id)
		if err != nil {
			// Skip missing annotations (might have been deleted)
			if status.Code(err) == codes.NotFound {
				continue
			}
			return nil, err
		}
		annotations = append(annotations, annotation)
	}

	return annotations, nil
}

func (s *LevelDBStorage) ListByDOI(
	doi string,
) ([]*feature.FeatureAnnotation, error) {
	if err := s.validator.Var(doi, "required"); err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("invalid doi: %v", err),
		)
	}
	indexKey := fmt.Sprintf("doi_index:%s", doi)
	idsData, err := s.db.Get([]byte(indexKey), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []*feature.FeatureAnnotation{}, nil
		}
		return nil, status.Errorf(
			codes.Internal,
			"failed to lookup DOI index: %v",
			err,
		)
	}

	ids := strings.Split(string(idsData), ",")
	annotations := make([]*feature.FeatureAnnotation, 0, len(ids))

	for _, id := range ids {
		if id == "" {
			continue
		}
		annotation, err := s.getByIDInternal(id)
		if err != nil {
			// Skip missing annotations (might have been deleted)
			if status.Code(err) == codes.NotFound {
				continue
			}
			return nil, err
		}
		annotations = append(annotations, annotation)
	}

	return annotations, nil
}

// Internal helper methods

func (s *LevelDBStorage) getByIDInternal(
	id string,
) (*feature.FeatureAnnotation, error) {
	key := fmt.Sprintf("annotation:%s", id)
	data, err := s.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, status.Errorf(
				codes.NotFound,
				"annotation with ID %s not found",
				id,
			)
		}
		return nil, status.Errorf(
			codes.Internal,
			"failed to retrieve annotation: %v",
			err,
		)
	}

	annotation := &feature.FeatureAnnotation{}
	if err := proto.Unmarshal(data, annotation); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to deserialize annotation: %v",
			err,
		)
	}

	if annotation.IsObsolete {
		return nil, status.Errorf(
			codes.NotFound,
			"annotation with ID %s is obsolete",
			id,
		)
	}

	return annotation, nil
}

func (s *LevelDBStorage) updateAnnotation(
	annotation *feature.FeatureAnnotation,
) error {
	data, err := proto.Marshal(annotation)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to serialize annotation: %v",
			err,
		)
	}
	key := fmt.Sprintf("annotation:%s", annotation.Id)
	if err := s.db.Put([]byte(key), data, nil); err != nil {
		return status.Errorf(
			codes.Internal,
			"failed to update annotation: %v",
			err,
		)
	}

	return nil
}

func (s *LevelDBStorage) updateIndexes(
	txn *leveldb.Transaction,
	annotation *feature.FeatureAnnotation,
	add bool,
) error {
	if err := s.updateNameIndex(txn, annotation, add); err != nil {
		return err
	}
	if err := s.updatePubmedIndexes(txn, annotation, add); err != nil {
		return err
	}
	return s.updateDOIIndexes(txn, annotation, add)
}

func (s *LevelDBStorage) updateNameIndex(
	txn *leveldb.Transaction,
	annotation *feature.FeatureAnnotation,
	add bool,
) error {
	nameKey := fmt.Sprintf("name_index:%s", annotation.Attributes.Name)
	if add {
		return txn.Put([]byte(nameKey), []byte(annotation.Id), nil)
	}
	return txn.Delete([]byte(nameKey), nil)
}

func (s *LevelDBStorage) updatePubmedIndexes(
	txn *leveldb.Transaction,
	annotation *feature.FeatureAnnotation,
	add bool,
) error {
	for _, pubmedId := range annotation.Attributes.Pubmed {
		if err := s.updateListIndexWithTxn(txn, "pubmed_index", pubmedId, annotation.Id, add); err != nil {
			return err
		}
	}
	return nil
}

func (s *LevelDBStorage) updateDOIIndexes(
	txn *leveldb.Transaction,
	annotation *feature.FeatureAnnotation,
	add bool,
) error {
	for _, publication := range annotation.Attributes.Publications {
		// Check if it's a DOI (basic check for DOI format)
		if strings.HasPrefix(publication, "10.") {
			if err := s.updateListIndexWithTxn(txn, "doi_index", publication, annotation.Id, add); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *LevelDBStorage) updateListIndexWithTxn(
	txn *leveldb.Transaction,
	indexType, indexValue, annotationID string,
	add bool,
) error {
	indexKey := fmt.Sprintf("%s:%s", indexType, indexValue)
	existingIDs, err := s.getIndexIDs(txn, indexKey)
	if err != nil {
		return err
	}

	var updatedIDs []string
	if add {
		updatedIDs = s.addIDToList(existingIDs, annotationID)
	} else {
		updatedIDs = s.removeIDFromList(existingIDs, annotationID)
	}

	return s.saveIndexIDs(txn, indexKey, updatedIDs)
}

func (s *LevelDBStorage) getIndexIDs(
	txn *leveldb.Transaction,
	indexKey string,
) ([]string, error) {
	data, err := txn.Get([]byte(indexKey), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	if err == leveldb.ErrNotFound {
		return []string{}, nil
	}
	return strings.Split(string(data), ","), nil
}

func (s *LevelDBStorage) addIDToList(
	existingIDs []string,
	annotationID string,
) []string {
	for _, id := range existingIDs {
		if id == annotationID {
			return existingIDs // Already exists
		}
	}
	return append(existingIDs, annotationID)
}

func (s *LevelDBStorage) removeIDFromList(
	existingIDs []string,
	annotationID string,
) []string {
	newIDs := make([]string, 0, len(existingIDs))
	for _, id := range existingIDs {
		if id != annotationID && id != "" {
			newIDs = append(newIDs, id)
		}
	}
	return newIDs
}

func (s *LevelDBStorage) saveIndexIDs(
	txn *leveldb.Transaction,
	indexKey string,
	ids []string,
) error {
	if len(ids) == 0 || (len(ids) == 1 && ids[0] == "") {
		return txn.Delete([]byte(indexKey), nil)
	}
	updatedData := strings.Join(ids, ",")
	return txn.Put([]byte(indexKey), []byte(updatedData), nil)
}
