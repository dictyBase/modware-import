	"github.com/dictyBase/arangomanager"
func queryArango(params *queryArangoParams) error {
	params.logger.Debugf("Executing ArangoDB query: %s", params.aqlQuery)
	cursor, err := params.dbh.Search(params.aqlQuery)
	if err != nil {
		return fmt.Errorf("failed to execute ArangoDB query: %w", err)
	}
	defer cursor.Close()

	if cursor.IsEmpty() {
		params.logger.Error("No feature props found")
		return nil
	}

	docCount := 0
	for cursor.Scan() {
		select {
		case <-params.ctx.Done():
			params.logger.Warn("Main context cancelled during ArangoDB query.")
			return params.ctx.Err()
		default:
			var doc ArangoResultDoc
			if err := cursor.Read(&doc); err != nil {
				params.logger.Errorf(
					"Failed to read document from ArangoDB cursor: %v",
					err,
				)
				continue
			}
			select {
			case params.arangoDocsChan <- doc:
				docCount++
				if docCount%100 == 0 {
					params.logger.Infof(
						"Fetched %d documents from ArangoDB...",
						docCount,
					)
				}
			case <-params.ctx.Done():
				params.logger.Warn(
					"Main context cancelled while sending Arango doc to channel.",
				)
				return params.ctx.Err()
			}
		}
	}
	params.logger.Infof(
		"Successfully fetched %d documents from ArangoDB.",
		docCount,
	)
	return nil
}
