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
