package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// HandleS3ObjectParams contains parameters for handleS3Object function
type HandleS3ObjectParams struct {
	Context   *cli.Context
	Object    minio.ObjectInfo
	S3Client  *minio.Client
	OutputDir string
	Log       *logrus.Entry
}

// BuildArangoImportParams contains parameters for buildArangoImportCmd function
type BuildArangoImportParams struct {
	Context    *cli.Context
	Collection string
	FilePath   string
}

// ProcessS3ObjectParams contains parameters for processS3Object function
type ProcessS3ObjectParams struct {
	S3Client  *minio.Client
	Bucket    string
	ObjectKey string
}

// GenericResponse represents a generic JSON structure
type GenericResponse struct {
	Results []struct {
		Items []interface{} `json:"items"`
	} `json:"results"`
}

func buildArangoImportCmd(params BuildArangoImportParams) *exec.Cmd {
	// #nosec G204 -- Using CLI context values that are validated by the CLI framework
	return exec.Command("arangoimport",
		"--server.endpoint", fmt.Sprintf("tcp://%s:%s",
			params.Context.String("arangodb-host"),
			params.Context.String("arangodb-port")),
		"--server.database", params.Context.String("arangodb-database"),
		"--server.username", params.Context.String("arangodb-user"),
		"--server.password", params.Context.String("arangodb-pass"),
		"--collection", params.Collection,
		"--create-collection", "true",
		"--overwrite", "true",
		"--threads", "3",
		"--type", "json",
		"--file", params.FilePath,
	)
}

func processJSONFile(reader io.Reader) ([]interface{}, error) {
	response := &GenericResponse{}
	if err := json.NewDecoder(reader).Decode(response); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}
	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found in JSON")
	}
	return response.Results[0].Items, nil
}

func runArangoImport(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running arangoimport: %w", err)
	}
	return nil
}

func createOutputDirectory(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}
	return nil
}

func writeJSONToFile(items []interface{}, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", outputFile, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(items); err != nil {
		return fmt.Errorf(
			"error writing to output file %s: %w",
			outputFile,
			err,
		)
	}
	return nil
}

func processS3Object(params ProcessS3ObjectParams) ([]interface{}, error) {
	reader, err := params.S3Client.GetObject(
		params.Bucket,
		params.ObjectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting object %s: %w",
			params.ObjectKey,
			err,
		)
	}

	items, err := processJSONFile(reader)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func handleS3Object(params HandleS3ObjectParams) error {
	if filepath.Ext(params.Object.Key) != ".json" {
		return nil
	}

	collection := strings.ToLower(
		filepath.Base(params.Object.Key[:len(params.Object.Key)-5]),
	)
	outputFile := filepath.Join(
		params.OutputDir,
		fmt.Sprintf("%s.json", collection),
	)

	params.Log.WithFields(logrus.Fields{
		"file":        params.Object.Key,
		"collection":  collection,
		"output_file": outputFile,
	}).Info("processing file for import")

	items, err := processS3Object(ProcessS3ObjectParams{
		S3Client:  params.S3Client,
		Bucket:    params.Context.String("s3-bucket"),
		ObjectKey: params.Object.Key,
	})
	if err != nil {
		return err
	}
	params.Log.WithFields(logrus.Fields{
		"file":        params.Object.Key,
		"items_count": len(items),
	}).Info("successfully parsed JSON file")

	if err := writeJSONToFile(items, outputFile); err != nil {
		return err
	}
	params.Log.WithFields(logrus.Fields{
		"output_file": outputFile,
	}).Info("wrote JSON to file")

	cmd := buildArangoImportCmd(BuildArangoImportParams{
		Context:    params.Context,
		Collection: collection,
		FilePath:   outputFile,
	})
	params.Log.WithFields(logrus.Fields{
		"collection": collection,
		"input_file": outputFile,
	}).Info("starting import to ArangoDB")
	if err := runArangoImport(cmd); err != nil {
		return err
	}
	params.Log.WithFields(logrus.Fields{
		"collection": collection,
	}).Info("successfully imported data to ArangoDB")

	return nil
}

func LoadArangodb(cltx *cli.Context) error {
	log := registry.GetLogger()
	s3Client := registry.GetS3Client()
	bucket := cltx.String("s3-bucket")
	prefix := cltx.String("s3-bucket-path")
	outputDir := cltx.String("output-dir")

	if err := createOutputDirectory(outputDir); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	log.WithFields(logrus.Fields{
		"bucket":     bucket,
		"prefix":     prefix,
		"output_dir": outputDir,
	}).Info("starting ArangoDB data import process")

	doneCh := make(chan struct{})
	defer close(doneCh)

	for object := range s3Client.ListObjects(bucket, prefix, true, doneCh) {
		if object.Err != nil {
			return cli.Exit(
				fmt.Sprintf("error listing objects: %s", object.Err),
				2,
			)
		}

		if err := handleS3Object(HandleS3ObjectParams{
			Context:   cltx,
			Object:    object,
			S3Client:  s3Client,
			OutputDir: outputDir,
			Log:       log,
		}); err != nil {
			return cli.Exit(err.Error(), 2)
		}
	}

	log.Info("completed ArangoDB import process")
	return nil
}
