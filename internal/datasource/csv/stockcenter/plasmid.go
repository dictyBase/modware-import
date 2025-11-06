package stockcenter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

// PlasmidGenbank is the container for genbank link for plasmid
type PlasmidGenbank struct {
	Id      string
	Genbank string
}

// PlasmidGenbankReader is the defined interface for reading the data
type PlasmidGenbankReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidGenbank, error)
}

type csvPlasmidGenbankReader struct {
	*csource.CsvReader
}

// NewPlasmidGenbankReader is to get an instance of PlasmidGenbankReader
func NewPlasmidGenbankReader(r io.Reader) PlasmidGenbankReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidGenbankReader{&csource.CsvReader{Reader: cr}}
}

// Value gets a new PlasmidGenbank instance
func (pgr *csvPlasmidGenbankReader) Value() (*PlasmidGenbank, error) {
	g := new(PlasmidGenbank)
	if pgr.Err != nil {
		return g, pgr.Err
	}
	g.Id = pgr.Record[0]
	g.Genbank = pgr.Record[1]
	return g, nil
}

// Plasmid is the container for plasmid data
type Plasmid struct {
	Id           string
	Summary      string
	User         string
	CreatedOn    time.Time
	UpdatedOn    time.Time
	Name         string
	Publications []string
	Genes        []string
}

// PlasmidReader is the defined interface for reading the plasmid data
type PlasmidReader interface {
	datasource.IteratorWithoutValue
	Value() E.Either[error, *Plasmid]
}

type csvPlasmidReader struct {
	*csource.CsvReader
	alookup StockAnnotatorLookup
	plookup StockPubLookup
	glookup StockGeneLookup
}

// NewCsvPlasmidReader is to get an instance of PlasmidReader instance
func NewCsvPlasmidReader(
	r io.Reader,
	al StockAnnotatorLookup,
	pl StockPubLookup,
	gl StockGeneLookup,
) PlasmidReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidReader{
		CsvReader: &csource.CsvReader{Reader: cr},
		alookup:   al,
		plookup:   pl,
		glookup:   gl,
	}
}

// getRecordField safely accesses CSV record field by index using Option
func getRecordField(record []string, index int) O.Option[string] {
	if index < 0 || index >= len(record) {
		return O.None[string]()
	}
	return O.Some(record[index])
}

// parseId extracts and sets the plasmid ID from record
func parseId(record []string) func(*Plasmid) E.Either[error, *Plasmid] {
	return func(plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe1(
			getRecordField(record, 0),
			O.Fold(
				func() E.Either[error, *Plasmid] {
					return E.Left[*Plasmid](fmt.Errorf("missing id at index 0"))
				},
				func(id string) E.Either[error, *Plasmid] {
					plasmid.Id = id
					return E.Right[error](plasmid)
				},
			),
		)
	}
}

// ensurePlasmidPrefix ensures name starts with 'p' prefix using Option pattern
func ensurePlasmidPrefix(name string) string {
	return F.Pipe2(
		O.Some(name),
		O.Filter(func(n string) bool {
			return strings.HasPrefix(n, "p")
		}),
		O.GetOrElse(F.Constant(fmt.Sprintf("p%s", name))),
	)
}

// parseName extracts and normalizes plasmid name with prefix
func parseName(record []string) func(*Plasmid) E.Either[error, *Plasmid] {
	return func(plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe1(
			getRecordField(record, 1),
			O.Fold(
				func() E.Either[error, *Plasmid] {
					return E.Left[*Plasmid](fmt.Errorf("missing name at index 1"))
				},
				func(name string) E.Either[error, *Plasmid] {
					plasmid.Name = F.Pipe1(name, ensurePlasmidPrefix)
					return E.Right[error](plasmid)
				},
			),
		)
	}
}

// parseSummary extracts plasmid summary from record
func parseSummary(record []string) func(*Plasmid) E.Either[error, *Plasmid] {
	return func(plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe1(
			getRecordField(record, 2),
			O.Fold(
				func() E.Either[error, *Plasmid] {
					return E.Left[*Plasmid](fmt.Errorf("missing summary at index 2"))
				},
				func(summary string) E.Either[error, *Plasmid] {
					plasmid.Summary = summary
					return E.Right[error](plasmid)
				},
			),
		)
	}
}

// enrichWithAnnotator is a curried function that enriches plasmid with annotator data
var enrichWithAnnotator = F.Curry2(
	func(alookup StockAnnotatorLookup, plasmid *Plasmid) E.Either[error, *Plasmid] {
		user, createdOn, updatedOn, ok := alookup.StockAnnotator(plasmid.Id)
		if ok {
			plasmid.User = user
			plasmid.CreatedOn = createdOn
			plasmid.UpdatedOn = updatedOn
		}
		return E.Right[error](plasmid)
	},
)

// isNonEmptySlice checks if a slice is non-empty
func isNonEmptySlice(slice []string) bool {
	return len(slice) > 0
}

// enrichWithPublications is a curried function that enriches plasmid with publications
var enrichWithPublications = F.Curry2(
	func(plookup StockPubLookup, plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe2(
			plookup.StockPub(plasmid.Id),
			O.FromPredicate(isNonEmptySlice),
			O.Fold(
				func() E.Either[error, *Plasmid] {
					// No publications found, keep empty slice
					return E.Right[error](plasmid)
				},
				func(pubs []string) E.Either[error, *Plasmid] {
					// Filter out empty strings from publications
					filteredPubs := F.Pipe1(
						pubs,
						A.Filter(func(pub string) bool {
							return pub != ""
						}),
					)
					plasmid.Publications = append(plasmid.Publications, filteredPubs...)
					return E.Right[error](plasmid)
				},
			),
		)
	},
)

// enrichWithGenes is a curried function that enriches plasmid with genes
var enrichWithGenes = F.Curry2(
	func(glookup StockGeneLookup, plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe2(
			glookup.StockGene(plasmid.Id),
			O.FromPredicate(isNonEmptySlice),
			O.Fold(
				func() E.Either[error, *Plasmid] {
					// No genes found, keep empty slice
					return E.Right[error](plasmid)
				},
				func(genes []string) E.Either[error, *Plasmid] {
					// Filter out empty strings from genes
					filteredGenes := F.Pipe1(
						genes,
						A.Filter(func(gene string) bool {
							return gene != ""
						}),
					)
					plasmid.Genes = append(plasmid.Genes, filteredGenes...)
					return E.Right[error](plasmid)
				},
			),
		)
	},
)

// parsePlasmidFields is a curried function that composes all field parsers
var parsePlasmidFields = F.Curry2(
	func(record []string, plasmid *Plasmid) E.Either[error, *Plasmid] {
		return F.Pipe3(
			E.Of[error](plasmid),
			E.Chain(parseId(record)),
			E.Chain(parseName(record)),
			E.Chain(parseSummary(record)),
		)
	},
)

// Value gets a new Plasmid instance using fp-go Either pattern
func (plr *csvPlasmidReader) Value() E.Either[error, *Plasmid] {
	plasmid := new(Plasmid)

	// Check for CSV reader errors first
	if plr.Err != nil {
		return E.Left[*Plasmid](plr.Err)
	}

	// Parse all fields and enrich with lookups using functional composition
	return F.Pipe3(
		parsePlasmidFields(plr.Record)(plasmid),
		E.Chain(enrichWithAnnotator(plr.alookup)),
		E.Chain(enrichWithPublications(plr.plookup)),
		E.Chain(enrichWithGenes(plr.glookup)),
	)
}
