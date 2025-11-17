package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/config"
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
		"--server.endpoint", fmt.Sprintf("http+tcp://%s:%s",
			params.Context.String("arangodb-host"),
			params.Context.String("arangodb-port")),
		"--server.database", params.Context.String("arangodb-database"),
		"--server.username", params.Context.String("arangodb-user"),
		"--server.password", params.Context.String("arangodb-pass"),
		"--collection", params.Collection,
		"--create-collection", "true",
		"--overwrite", "true",
		"--type", "jsonl",
		"--file", params.FilePath,
	)
}

func runArangoImport(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running arangoimport: %w", err)
	}
	return nil
}

func createOutputDirectory(outputDir string) error {
	if err := os.MkdirAll(outputDir, config.SharedDirectoryPermission); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}
	return nil
}

func getCollectionAndOutputFile(objectKey, outputDir string) (string, string) {
	collection := strings.ToLower(
		filepath.Base(objectKey[:len(objectKey)-5]),
	)
	outputFile := filepath.Join(
		outputDir,
		fmt.Sprintf("%s.json", collection),
	)
	return collection, outputFile
}

func processAndWriteData(
	params HandleS3ObjectParams,
	outputFile string,
) error {
	// Get the object from S3
	reader, err := params.S3Client.GetObject(
		params.Context.String("s3-bucket"),
		params.Object.Key,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf(
			"error getting object %s: %w",
			params.Object.Key,
			err,
		)
	}
	defer reader.Close()
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	// Process the JSON
	if err := processJSON(reader, outFile); err != nil {
		return fmt.Errorf("error processing JSON: %w", err)
	}

	params.Log.WithFields(logrus.Fields{
		"output_file": outputFile,
	}).Info("wrote JSON to file")
	return nil
}

func importToArangoDB(
	params HandleS3ObjectParams,
	collection, outputFile string,
) error {
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

func handleS3Object(params HandleS3ObjectParams) error {
	if filepath.Ext(params.Object.Key) != ".json" {
		return nil
	}

	collection, outputFile := getCollectionAndOutputFile(
		params.Object.Key,
		params.OutputDir,
	)

	params.Log.WithFields(logrus.Fields{
		"file":        params.Object.Key,
		"collection":  collection,
		"output_file": outputFile,
	}).Info("processing file for import")

	if err := processAndWriteData(params, outputFile); err != nil {
		return err
	}

	if params.Context.Bool("skip-import") {
		params.Log.WithFields(logrus.Fields{
			"collection": collection,
		}).Info("skipping ArangoDB import due to skip-import flag")
		return nil
	}

	return importToArangoDB(params, collection, outputFile)
}

func LoadArangodb(cltx *cli.Context) error {
	log := registry.GetLogger()
	s3Client := registry.GetS3Client()
	bucket := cltx.String("s3-bucket")
	prefix := cltx.String("s3-bucket-path")
	outputDir := cltx.String("output-dir")

	if err := createOutputDirectory(outputDir); err != nil {
		return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
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
				config.DefaultRetryBackoffFactor,
			)
		}

		if err := handleS3Object(HandleS3ObjectParams{
			Context:   cltx,
			Object:    object,
			S3Client:  s3Client,
			OutputDir: outputDir,
			Log:       log,
		}); err != nil {
			return cli.Exit(err.Error(), config.DefaultRetryBackoffFactor)
		}
	}

	log.Info("completed ArangoDB import process")
	return nil
}

// processJSONToken handles a single JSON token and returns whether to continue processing
func processJSONToken(
	decoder *json.Decoder,
	encoder *json.Encoder,
	token json.Token,
	depth *int,
) (bool, error) {
	switch t := token.(type) {
	case json.Delim:
		return handleDelimiter(t, depth)
	case string:
		if t == "items" && *depth == 3 {
			return true, processItems(decoder, encoder)
		}
	}
	return true, nil
}

// handleDelimiter processes JSON delimiters and tracks depth
func handleDelimiter(delim json.Delim, depth *int) (bool, error) {
	switch delim {
	case '{', '[':
		*depth++
	case '}', ']':
		*depth--
		if *depth == 0 {
			return false, nil
		}
	}
	return true, nil
}

// processItems handles the contents of an "items" array
func processItems(decoder *json.Decoder, encoder *json.Encoder) error {
	// Skip the opening bracket of items array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error in start decoding items: %w", err)
	}

	// Process each item in the array
	for decoder.More() {
		var item interface{}
		if err := decoder.Decode(&item); err != nil {
			return fmt.Errorf("error decoding item: %w", err)
		}
		if err := encoder.Encode(item); err != nil {
			return fmt.Errorf("error encoding item: %w", err)
		}
	}
	return nil
}

// processJSON reads a JSON input file and writes specific items to an output file.
// It expects a JSON structure with a specific nesting pattern:
// {                  // depth 1
//
//	  "results": [     // depth 2
//	    {              // depth 3
//	      "items": [   // depth 4
//	        {...},     // individual items
//	        {...}
//	      ]
//	    }
//	  ]
//	}
func processJSON(inputFile io.Reader, outputFile io.Writer) error {
	decoder := json.NewDecoder(inputFile)
	encoder := json.NewEncoder(outputFile)
	depth := 0

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding token: %w", err)
		}

		continueProcessing, err := processJSONToken(
			decoder,
			encoder,
			token,
			&depth,
		)
		if err != nil {
			return err
		}
		if !continueProcessing {
			return nil
		}
	}
	return nil
}
