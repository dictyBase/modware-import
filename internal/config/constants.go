// Package config provides shared configuration constants used across the codebase
package config

const (
	// Database and Service Ports
	DefaultArangoDBPort = 8529  // Standard port for ArangoDB connections
	DefaultGRPCPort     = 9000  // Default gRPC service port
	DefaultMetricsPort  = 9560  // Default metrics/monitoring port
	UniProtPort         = 44689 // UniProt service port

	// File Permissions
	DefaultFilePermission      = 0o600 // Owner read/write only
	DefaultDirectoryPermission = 0o700 // Owner read/write/execute only
	SharedDirectoryPermission  = 0o755 // Owner full, others read/execute

	// Worker Pool Defaults
	DefaultWorkerPoolSize     = 4
	DefaultBatchSize          = 100
	DefaultPlasmidBatchSize   = 30
	DefaultCSVWorkerPoolSize  = 20
	DefaultTimeoutSeconds     = 300
	DefaultRetryBackoffFactor = 2

	// CSV/Data Processing
	MinimumFieldCount     = 2  // Minimum fields required in CSV records
	MinimumTokenCount     = 3  // Minimum tokens for valid parsing
	GoldenBraidFieldCount = 7  // Required fields in GoldenBraid CSV
	InventoryFieldCount   = 9  // Required fields in inventory TSV
	PhenotypeColumnCount  = 10 // Standard phenotype column count
	ExtendedColumnCount   = 12 // Extended phenotype column count

	// Pagination and Limits
	DefaultPageSize   = 10
	DefaultMaxResults = 1000
	DefaultLineWidth  = 80

	// Buffer and Formatting
	DefaultBufferSize = 8
	DefaultLineLimit  = 15
)
