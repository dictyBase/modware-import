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
