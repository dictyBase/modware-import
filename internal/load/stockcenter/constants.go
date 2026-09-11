package stockcenter

type Status int

const (
	Created Status = iota
	Updated
	Deleted
	Read
	Nop
)

const (
	sysnameTag   = "systematic name"
	mutmethodTag = "mutagenesis method"
	muttypeTag   = "mutant type"
	genoTag      = "genotype"
	synTag       = "synonym"
	val          = "novalue"
)

// Logging field keys shared across the stockcenter loaders.
const (
	logTypeKey   = "type"
	logStockKey  = "stock"
	logEventKey  = "event"
	logCountKey  = "count"
	logReadKey   = "read"
	logFolderKey = "folder"
	logBucketKey = "bucket"
)

// Event values recorded under the "event" log field.
const (
	evLoad   = "load"
	evCreate = "create"
	evDelete = "delete"
	evRead   = "read"
)

// Stock type values used in content attributes and log fields.
const (
	plasmidType   = "plasmid"
	strainType    = "strain"
	inventoryType = "inventory"
	phenotypeType = "phenotype"
)
