package strain

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"slices"

	R "github.com/IBM/fp-go/context/readerioeither"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	J "github.com/IBM/fp-go/json"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/datasource/xls/strain"
	"github.com/sirupsen/logrus"
)

const ConcurrentStrainLoader = 10

type fnRunnerProperties struct {
	fn    func(*strain.StrainAnnotation) (string, error)
	props *strain.StrainAnnotation
}

type StrainLoader struct {
	Host         string
	Token        string
	TableId      int
	Logger       *logrus.Entry
	StrainReader *strain.StrainAnnotationReader
}

func (loader *StrainLoader) Load() error {
	loaderSlice := make([]*fnRunnerProperties, 0, ConcurrentStrainLoader)
	strainReader := loader.StrainReader

	for strainReader.Next() {
		strain, err := strainReader.Value()
		if err != nil {
			return err
		}
		loaderSlice = append(loaderSlice, &fnRunnerProperties{
			fn:    loader.addStrainRow,
			props: strain,
		})

		if len(loaderSlice) == ConcurrentStrainLoader {
			if err := processFnRunnerProperties(loaderSlice, loader.Logger); err != nil {
				return err
			}
			loaderSlice = slices.Delete(loaderSlice, 0, len(loaderSlice))
			// Another way to do this
			// loaderSlice = loaderSlice[:0] // Reset the slice without allocating new memory
		}
	}

	// Process remaining items in loaderSlice
	if len(loaderSlice) > 0 {
		if err := processFnRunnerProperties(loaderSlice, loader.Logger); err != nil {
			return err
		}
	}

	return nil
}

func (loader *StrainLoader) addRowParameters(
	strn *strain.StrainAnnotation,
) map[string]interface{} {
	params := map[string]interface{}{
		"strain_descriptor":         strn.Descriptor(),
		"species":                   strn.Species(),
		"assigned_by":               []string{strn.AssignedBy()},
		"reference":                 strn.Reference(),
		"strain_summary":            strn.Summary(),
		"strain_characteristics_id": []string{strn.Characteristic()},
		"genetic_modification_id":   []string{strn.GeneticModification()},
		"mutagenesis_method_id":     []string{strn.MutagenesisMethod()},
	}
	if strn.HasName() {
		params["strain_names"] = strn.Name()
	}
	if strn.HasSystematicName() {
		params["systematic_name"] = strn.SystematicName()
	}
	if strn.HasPlasmid() {
		params["plasmid"] = strn.Plasmid()
	}
	if strn.HasParentId() {
		params["parent_strain_id"] = strn.ParentId()
	}
	if strn.HasGenes() {
		params["associated_genes"] = strn.Genes()
	}
	if strn.HasGenotype() {
		params["genotype"] = strn.Genotype()
	}
	if strn.HasDepositor() {
		params["depositor"] = strn.Depositor()
	}
	return params
}

func (loader *StrainLoader) createStrainURL() string {
	return fmt.Sprintf(
		"https://%s/api/database/rows/table/%d/?user_field_names=true",
		loader.Host,
		loader.TableId,
	)
}

func (loader *StrainLoader) addStrainRow(
	strn *strain.StrainAnnotation,
) (string, error) {
	var empty string
	createPayload := F.Pipe3(
		strn,
		loader.addRowParameters,
		J.Marshal,
		E.Fold(httpapi.OnJSONPayloadError, httpapi.OnJSONPayloadSuccess),
	)
	if createPayload.Error != nil {
		return empty, createPayload.Error
	}
	resp := F.Pipe3(
		loader.createStrainURL(),
		httpapi.MakeHTTPRequest("POST", bytes.NewBuffer(createPayload.Payload)),
		R.Map(httpapi.SetHeaderWithJWT(loader.Token)),
		strainCreateHTTP,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold[error, strainCreateResp, httpapi.ResponseFeedback](
			onStrainCreateFeedbackError,
			onStrainCreateFeedbackSuccess,
		),
	)
	return output.Msg, output.Err
}

func executeLoaderSlice(
	loaderSlice []*fnRunnerProperties,
) (chan string, chan error) {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure all paths cancel the context to avoid context leak

	// channel to communicate error and result
	resultCh := make(chan string)
	errCh := make(chan error)
	var wg sync.WaitGroup

	// Run each function in a goroutine
	for _, loader := range loaderSlice {
		wg.Add(1)
		go func(ldr *fnRunnerProperties) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				result, err := ldr.fn(ldr.props)
				if err != nil {
					cancel()
					errCh <- err
					return
				}
				resultCh <- result
			}
		}(loader)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	return resultCh, errCh
}

func processFnRunnerProperties(
	loaderSlice []*fnRunnerProperties,
	logger *logrus.Entry,
) error {
	resultCh, errCh := executeLoaderSlice(loaderSlice)
	for {
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case result, ok := <-resultCh:
			if !ok {
				return nil
			}
			logger.Debugf(result)
		}
	}
}
