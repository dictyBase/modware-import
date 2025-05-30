	"github.com/dictyBase/arangomanager"
func queryArango(params *queryArangoParams) error {
	params.logger.Debugf("Executing ArangoDB query: %s", params.aqlQuery)
	cursor, err := params.dbh.Search(params.aqlQuery)
// queryArangoParams holds the parameters for the queryArango function.
type queryArangoParams struct {
	ctx            context.Context
	config         AppConfig // Contains AQLQuery and Logger
	arangoDocsChan chan<- ArangoResultDoc
	mainCancel     context.CancelFunc
}
	if err != nil {
		return fmt.Errorf("failed to execute ArangoDB query: %w", err)
	}
	defer cursor.Close()

	if cursor.IsEmpty() {
		params.logger.Error("No feature props found")
		return nil
	}

// queryArango is responsible for querying ArangoDB and sending documents to a channel.
// It calls params.mainCancel if critical errors occur or the context is done.
func queryArango(params *queryArangoParams) {
	dbh := registry.GetArangodbConnection()
	cursor, err := dbh.Search(params.config.AQLQuery)
	if err != nil {
		params.config.Logger.Errorf("error in arangodb query %v", err)
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	defer cursor.Close()
	if cursor.IsEmpty() {
		params.config.Logger.Warn("ArangoDB querier finished (no data).")
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	docCount := 0
	for cursor.Scan() {
		select {
		case <-params.ctx.Done():
			return
		default:
			var doc ArangoResultDoc
			if errRead := cursor.Read(&doc); errRead != nil {
				params.config.Logger.Errorf(
					"Failed to read document from ArangoDB cursor: %v",
					errRead,
				)
				params.mainCancel() // Signal other goroutines to stop
			}
			// Send the document to the channel, handling context
			// cancellation
			select {
			case params.arangoDocsChan <- doc:
				docCount++
			case <-params.ctx.Done():
				return
			}
		}
	}
	params.config.Logger.Infof(
		"Successfully fetched %d documents from ArangoDB.",
		docCount,
	)
}

func grpcUpdateWorkerFunc(
	config AppConfig,
	grpcClient feature_annotation.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[ProcessedGeneData, GrpcUpdateResult] {
	return func(ctx context.Context, job concurrent.Job[ProcessedGeneData]) (GrpcUpdateResult, error) {
		logger := config.Logger
		processedData := job.Payload
		logger.Debugf(
			"gRPC Worker (Job %s): updating gene ID: %s",
			job.ID,
			processedData.GeneID,
		)
		result := GrpcUpdateResult{GeneID: processedData.GeneID, Success: false}
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx, &feature_annotation.FeatureAnnotationId{
				Id: processedData.GeneID,
			})
		if err != nil {
			result.Message = fmt.Sprintf(
				"failed to GetFeatureAnnotation: %v",
				err,
			)
			result.Error = err
			return result, err
		}
		for _, prop := range processedData.StrippedPropsText {
			_, err := grpcClient.AddTag(ctx,
				&feature_annotation.AddTagRequest{
					Id: featAnno.Id,
					Tag: &feature_annotation.TagPropertyCreate{
						Tag:       prop.OriginalName,
						Value:     prop.StrippedText,
						CreatedBy: config.ArangoUser,
					},
				})
			if err != nil {
				errMsg := fmt.Sprintf(
					"failed to AddTag for property %s: %v",
					prop.OriginalName,
					err,
				)
				result.Message = errMsg
				result.Error = err
				logger.Errorf(
					"gRPC Worker (Job %s): %s for gene ID %s",
					job.ID, errMsg, processedData.GeneID,
				)
				return result, err
			}
			logger.Debugf(
				"gRPC Worker (Job %s): successfully added tag %s for gene ID %s",
				job.ID,
				prop.OriginalName,
				processedData.GeneID,
			)
		}
		result.Success = true
		result.Message = fmt.Sprintf(
			"Successfully added %d tags",
			len(processedData.StrippedPropsText),
		)
		return result, nil
	}
}
