package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var noncharReg = regexp.MustCompile("[^a-z0-9]+")

const updatedBy = "pfey@northwestern.edu"

func Slugify(name string) string {
	return strings.Trim(
		noncharReg.ReplaceAllString(strings.ToLower(name), "-"),
		"-",
	)
}

// sourceItem is a validated S3 object ready for upsert.
type sourceItem struct {
	key       string
	name      string
	namespace string
	slug      string
	payload   string
}

// objectSource abstracts the S3 operations needed by the loader so preflight
// can be tested without a live MinIO client.
type objectSource interface {
	list(bucket, prefix string, doneCh <-chan struct{}) <-chan minio.ObjectInfo
	get(bucket, key string) (io.ReadCloser, error)
}

type minioObjectSource struct {
	client *minio.Client
}

func (s minioObjectSource) list(
	bucket, prefix string, doneCh <-chan struct{},
) <-chan minio.ObjectInfo {
	return s.client.ListObjects(bucket, prefix, true, doneCh)
}

func (s minioObjectSource) get(bucket, key string) (io.ReadCloser, error) {
	return s.client.GetObject(bucket, key, minio.GetObjectOptions{})
}

func LoadContent(cltx *cli.Context) error {
	logger := registry.GetLogger()
	s3Client := registry.GetS3Client()
	client := regsc.GetContentAPIClient()

	src := minioObjectSource{client: s3Client}
	err := loadContent(
		client,
		logger,
		src,
		cltx.String("s3-bucket"),
		cltx.String("s3-bucket-path"),
	)
	if err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
	}
	return nil
}

// loadContent runs preflight (no API mutations) then mutates sequentially in
// source listing order.
func loadContent(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	src objectSource,
	bucket, prefix string,
) error {
	items, err := preflight(src, logger, bucket, prefix)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := upsertContent(client, logger, item); err != nil {
			return err
		}
	}
	return nil
}

// preflight streams the S3 listing, keeps the first item per slug in listing
// order, and retrieves plus validates only those survivors. It performs no
// content API mutation and returns before any mutation when a non-duplicate
// check fails. Duplicates are warned about and skipped, never fetched.
func preflight(
	src objectSource, logger *logrus.Entry, bucket, prefix string,
) ([]sourceItem, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	var items []sourceItem
	seen := mapset.NewSet[string]()
	for info := range src.list(bucket, prefix, doneCh) {
		if info.Err != nil {
			return nil, fmt.Errorf("error listing object %s: %w", info.Key, info.Err)
		}
		item, err := sourceItemFromKey(info.Key)
		if err != nil {
			return nil, err
		}
		if seen.Contains(item.slug) {
			logger.Warnf(
				"duplicate slug %q from key %s; skipping",
				item.slug,
				item.key,
			)
			continue
		}
		seen.Add(item.slug)
		read, err := readSourceItem(src, bucket, item)
		if err != nil {
			return nil, err
		}
		items = append(items, read)
	}

	return items, nil
}

// sourceItemFromKey derives item identity from the object key alone; no S3 I/O
// and no body is required, which is why dedupe can precede retrieval.
func sourceItemFromKey(key string) (sourceItem, error) {
	name, namespace, err := parseSourceKey(key)
	if err != nil {
		return sourceItem{}, err
	}

	slug := Slugify(fmt.Sprintf("%s %s", namespace, name))
	if slug == "" {
		return sourceItem{}, fmt.Errorf("empty slug derived from object %s", key)
	}

	return sourceItem{
		key:       key,
		name:      name,
		namespace: namespace,
		slug:      slug,
	}, nil
}

// readSourceItem retrieves and validates a single survivor item's body.
func readSourceItem(
	src objectSource, bucket string, item sourceItem,
) (sourceItem, error) {
	reader, err := src.get(bucket, item.key)
	if err != nil {
		return sourceItem{}, fmt.Errorf("error getting object %s: %w", item.key, err)
	}
	payload, err := readAllAndClose(reader, item.key)
	if err != nil {
		return sourceItem{}, err
	}
	if err := validatePayload(item.key, payload); err != nil {
		return sourceItem{}, err
	}

	item.payload = string(payload)

	return item, nil
}

// readAllAndClose reads the object body and closes it on every path, surfacing
// a close error only when no earlier error occurred.
func readAllAndClose(reader io.ReadCloser, key string) (payload []byte, err error) {
	defer func() {
		if cerr := reader.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing object %s: %w", key, cerr)
		}
	}()
	payload, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading content for object %s: %w", key, err)
	}
	return payload, nil
}

// validatePayload requires a non-empty payload that is valid JSON of any shape.
func validatePayload(key string, payload []byte) error {
	if strings.TrimSpace(string(payload)) == "" {
		return fmt.Errorf("empty content payload for object %s", key)
	}
	if !json.Valid(payload) {
		return fmt.Errorf("invalid JSON payload for object %s", key)
	}
	return nil
}

// parseSourceKey validates the basename grammar
// `<dfp|dsc|news>-<nonempty name>.<nonempty extension>` and returns the
// mapped namespace with the parsed name.
func parseSourceKey(key string) (name, namespace string, err error) {
	base := path.Base(key)

	stem, ext, ok := strings.Cut(base, ".")
	if !ok || stem == "" {
		return "", "", fmt.Errorf(
			"malformed object key %s: missing file extension", key,
		)
	}
	if ext == "" {
		return "", "", fmt.Errorf(
			"malformed object key %s: empty file extension", key,
		)
	}

	prefix, name, ok := strings.Cut(stem, "-")
	if !ok {
		return "", "", fmt.Errorf(
			"malformed object key %s: missing namespace prefix", key,
		)
	}
	ns, exists := namespaceMap()[prefix]
	if !exists {
		return "", "", fmt.Errorf(
			"malformed object key %s: unknown namespace prefix %q", key, prefix,
		)
	}
	if name == "" {
		return "", "", fmt.Errorf("malformed object key %s: missing name", key)
	}

	return name, ns, nil
}

func namespaceMap() map[string]string {
	return map[string]string{
		"dfp":  "frontpage",
		"dsc":  "stockcenter",
		"news": "news",
	}
}

// upsertContent creates, updates, or skips a single content record.
func upsertContent(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	item sourceItem,
) error {
	existing, err := client.GetContentBySlug(
		context.Background(),
		&content.ContentRequest{Slug: item.slug},
	)
	if err == nil {
		return applyExistingContent(client, logger, item, existing)
	}
	if status.Code(err) == codes.NotFound {
		return createContent(client, logger, item)
	}
	return fmt.Errorf(
		"error fetching content %s (slug %s): %w",
		item.key,
		item.slug,
		err,
	)
}

// applyExistingContent compares stored content against the fetched payload and
// updates only when they differ.
func applyExistingContent(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	item sourceItem,
	existing *content.Content,
) error {
	id, stored, err := existingRecord(existing)
	if err != nil {
		return fmt.Errorf(
			"invalid existing content %s (slug %s): %w",
			item.key,
			item.slug,
			err,
		)
	}
	if stored == item.payload {
		logger.Infof("unchanged/skipped content %s (slug %s)", item.key, item.slug)
		return nil
	}
	return updateContent(client, logger, item, id)
}

func createContent(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	item sourceItem,
) error {
	resp, err := client.StoreContent(
		context.Background(),
		&content.StoreContentRequest{
			Data: &content.StoreContentRequest_Data{
				Attributes: &content.NewContentAttributes{
					Name:      item.name,
					Namespace: item.namespace,
					CreatedBy: updatedBy,
					Content:   item.payload,
					Slug:      item.slug,
				},
			},
		},
	)
	if status.Code(err) == codes.AlreadyExists {
		return refetchAfterRace(client, logger, item)
	}
	if err != nil {
		return fmt.Errorf(
			"error creating content %s (slug %s): %w",
			item.key,
			item.slug,
			err,
		)
	}
	if resp == nil || resp.Data == nil || resp.Data.Attributes == nil {
		return fmt.Errorf(
			"invalid create response for content %s (slug %s): missing data or attributes",
			item.key,
			item.slug,
		)
	}
	logger.Infof("created content %s (slug %s)", item.key, item.slug)
	return nil
}

func updateContent(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	item sourceItem,
	id int64,
) error {
	resp, err := client.UpdateContent(
		context.Background(),
		&content.UpdateContentRequest{
			Id: id,
			Data: &content.UpdateContentRequest_Data{
				Attributes: &content.ExistingContentAttributes{
					UpdatedBy: updatedBy,
					Content:   item.payload,
				},
			},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"content"}},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"error updating content %s (slug %s): %w",
			item.key,
			item.slug,
			err,
		)
	}
	if resp == nil || resp.Data == nil || resp.Data.Attributes == nil {
		return fmt.Errorf(
			"invalid update response for content %s (slug %s): missing data or attributes",
			item.key,
			item.slug,
		)
	}
	logger.Infof("updated content %s (slug %s)", item.key, item.slug)
	return nil
}

// refetchAfterRace re-fetches the slug after a concurrent create and applies
// the existing-record comparison/update logic.
func refetchAfterRace(
	client content.ContentServiceClient,
	logger *logrus.Entry,
	item sourceItem,
) error {
	existing, err := client.GetContentBySlug(
		context.Background(),
		&content.ContentRequest{Slug: item.slug},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf(
				"content %s (slug %s) disappeared after concurrent create: %w",
				item.key,
				item.slug,
				err,
			)
		}
		return fmt.Errorf(
			"error refetching content %s (slug %s) after concurrent create: %w",
			item.key,
			item.slug,
			err,
		)
	}
	return applyExistingContent(client, logger, item, existing)
}

// existingRecord validates a successful Get response and returns the record id
// and stored content string.
func existingRecord(resp *content.Content) (int64, string, error) {
	if resp == nil || resp.Data == nil {
		return 0, "", fmt.Errorf("missing data in response")
	}
	if resp.Data.Attributes == nil {
		return 0, "", fmt.Errorf("missing attributes in response")
	}
	if resp.Data.Id <= 0 {
		return 0, "", fmt.Errorf("invalid record id %d in response", resp.Data.Id)
	}
	return resp.Data.Id, resp.Data.Attributes.Content, nil
}
