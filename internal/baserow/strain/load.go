package strain

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	R "github.com/IBM/fp-go/context/readerioeither"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/datasource/xls/strain"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const ConcurrentStrainLoader = 10

type StrainPayload struct {
	Descriptor              string `json:"strain_descriptor,omitempty"`
	Species                 string `json:"species,omitempty"`
	Reference               string `json:"reference,omitempty"`
	Summary                 string `json:"strain_summary,omitempty"`
	GeneticModificationId   []int  `json:"genetic_modification_id,omitempty"`
	StrainCharacteristicsId []int  `json:"strain_characteristics_id,omitempty"`
	MutagenesisMethodId     []int  `json:"mutagenesis_method_id,omitempty"`
	AssignedBy              []int  `json:"assigned_by,omitempty"`
	Names                   string `json:"strain_names,omitempty"`
	SystematicName          string `json:"systematic_name,omitempty"`
	Plasmid                 string `json:"plasmid,omitempty"`
	ParentId                string `json:"parent_strain_id,omitempty"`
	Genes                   string `json:"associated_genes,omitempty"`
	Genotype                string `json:"genotype,omitempty"`
	Depositor               string `json:"depositor,omitempty"`
}

type fnRunnerProperties struct {
	fn    func(*strain.StrainAnnotation) (string, error)
	props *strain.StrainAnnotation
}

type StrainLoader struct {
	Host             string
	Token            string
	TableId          int
	Logger           *logrus.Entry
	OntologyTableMap map[string]int
	TableManager     *database.TableManager
	Payload          *StrainPayload
	Annotation       *strain.StrainAnnotation
}

func NewStrainLoader(
	host, token string,
	tableId int,
	logger *logrus.Entry,
	tblMap map[string]int,
	manager *database.TableManager,
) *StrainLoader {
	return &StrainLoader{
		Host:             host,
		Token:            token,
		TableId:          tableId,
		Logger:           logger,
		OntologyTableMap: tblMap,
		TableManager:     manager,
	}
}

func (loader *StrainLoader) Load(reader *strain.StrainAnnotationReader) error {
	loaderSlice := make([]*fnRunnerProperties, 0, ConcurrentStrainLoader)
	for reader.Next() {
		strain, err := reader.Value()
		if strain.IsEmpty() {
			continue
		}
		if err != nil {
			return err
		}
		loader.Logger.Infof("got strain descriptor %s", strain.Descriptor())
		loaderSlice = append(loaderSlice, &fnRunnerProperties{
			fn:    loader.addStrainRow,
			props: strain,
		})
		if len(loaderSlice) == ConcurrentStrainLoader {
			loader.Logger.Debug("going to load strain")
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
		loader.Logger.Debug("going to load remaining strains")
		if err := processFnRunnerProperties(loaderSlice, loader.Logger); err != nil {
			return err
		}
	}

	return nil
}

func (loader *StrainLoader) addStrain(
	strn *strain.StrainAnnotation,
) E.Either[error, *StrainLoader] {
	loader.Annotation = strn
	return E.Right[error](loader)
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
	content := F.Pipe7(
		E.Do[error](strn),
		E.Bind(initialPayload, loader.addStrain),
		E.Bind(charIdsHandler, characteristicIds),
		E.Bind(genModIdHandler, genmodId),
		E.Bind(mutagenesisIdHandler, mutagenesisId),
		E.Map[error, *StrainLoader](loaderToPayload),
		E.Chain[error, *StrainPayload](marshalPayload),
		E.Fold(httpapi.OnJSONPayloadError, httpapi.OnJSONPayloadSuccess),
	)
	if content.Error != nil {
		return empty, content.Error
	}
	resp := F.Pipe3(
		loader.createStrainURL(),
		httpapi.MakeHTTPRequest("POST", bytes.NewBuffer(content.Payload)),
		R.Map(httpapi.SetHeaderWithJWT(loader.Token)),
		strainCreateHTTP,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold(onStrainCreateFeedbackError, onStrainCreateFeedbackSuccess),
	)
	return output.Msg, output.Err
}

func executeLoaderSlice(
	loaderSlice []*fnRunnerProperties,
) (chan string, chan error) {
	// channel to communicate error and result
	resultCh := make(chan string)
	errCh := make(chan error)
	var wg sync.WaitGroup

	// Run each function in a goroutine
	for _, loader := range loaderSlice {
		wg.Add(1)
		go func(ldr *fnRunnerProperties) {
			defer wg.Done()
			result, err := ldr.fn(ldr.props)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- result
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
	logger.Debugf("going process %d records", len(loaderSlice))
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
			logger.Infof(result)
		}
	}
}
