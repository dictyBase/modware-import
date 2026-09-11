package cli

const (
	// DefaultUserName is the default creator/updater for annotations
	DefaultUserName = "dcr@dictycr.org"
)

// AnnMap maps legacy creator usernames to their email addresses.
var AnnMap = map[string]string{
	"CGM_DDB_PASC": "pgaudet@northwestern.edu",
	"CGM_DDB_PFEY": "pfey@northwestern.edu",
	"CGM_DDB_BOBD": "robert-dodson@northwestern.edu",
	"CGM_DDB_KPIL": "kpilchar@northwestern.edu",
	"CGM_DDB":      "dictybase@northwestern.edu",
}

// Logging and metric field keys shared across the annotation loaders.
const (
	jobIDKey          = "job_id"
	stageKey          = "stage"
	errorKey          = "error"
	geneIDKey         = "gene_id"
	readFromDBKey     = "read_from_db"
	totalProcessedKey = "total_processed"
	successCountKey   = "success_count"
	errorCountKey     = "error_count"
	processingRateKey = "processing_rate"
	elapsedTimeKey    = "elapsed_time"
	grpcSubmittedKey  = "grpc_submitted"
	grpcCompletedKey  = "grpc_completed"
	skippedCountKey   = "skipped_count"
)

// Pipeline stage value recorded under the "stage" field.
const stageSubmittedToGRPCPool = "submitted_to_grpc_pool"

// CLI flag names and help text.
const (
	inputFlagName       = "input"
	workersFlagName     = "workers"
	batchSizeFlagName   = "batch-size"
	userFlagName        = "user"
	grpcWorkersFlagName = "grpc-workers"
	grpcWorkersEnvVar   = "GRPC_WORKERS"
	userEmailUsage      = "email address of the user running the load"
	grpcWorkersUsage    = "Number of gRPC update workers"
)

// Annotation tag value.
const productTag = "product"
