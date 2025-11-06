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
	RE "github.com/IBM/fp-go/readereither"

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

// ParseContext carries data through the parsing pipeline
type ParseContext struct {
	Record  []string
	Plasmid *Plasmid
}

// Dependencies holds lookup services for enrichment
type Dependencies struct {
	Alookup StockAnnotatorLookup
	Plookup StockPubLookup
	Glookup StockGeneLookup
}

// PlasmidParser represents a parsing step in the pipeline
type PlasmidParser = RE.ReaderEither[Dependencies, error, ParseContext]

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
func parseId(ctx ParseContext) PlasmidParser {
	return RE.FromEither[Dependencies](
		F.Pipe1(
			getRecordField(ctx.Record, 0),
			O.Fold(
				func() E.Either[error, ParseContext] {
					return E.Left[ParseContext](
						fmt.Errorf("missing id at index 0"),
					)
				},
				func(id string) E.Either[error, ParseContext] {
					ctx.Plasmid.Id = id
					return E.Right[error](ctx)
				},
			),
		),
	)
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
func parseName(ctx ParseContext) PlasmidParser {
	return RE.FromEither[Dependencies](
		F.Pipe1(
			getRecordField(ctx.Record, 1),
			O.Fold(
				func() E.Either[error, ParseContext] {
					return E.Left[ParseContext](
						fmt.Errorf("missing name at index 1"),
					)
				},
				func(name string) E.Either[error, ParseContext] {
					ctx.Plasmid.Name = F.Pipe1(name, ensurePlasmidPrefix)
					return E.Right[error](ctx)
				},
			),
		),
	)
}

// parseSummary extracts plasmid summary from record
func parseSummary(ctx ParseContext) PlasmidParser {
	return RE.FromEither[Dependencies](
		F.Pipe1(
			getRecordField(ctx.Record, 2),
			O.Fold(
				func() E.Either[error, ParseContext] {
					return E.Left[ParseContext](
						fmt.Errorf("missing summary at index 2"),
					)
				},
				func(summary string) E.Either[error, ParseContext] {
					ctx.Plasmid.Summary = summary
					return E.Right[error](ctx)
				},
			),
		),
	)
}

// applyAnnotator applies annotator lookup to plasmid context
func applyAnnotator(ctx ParseContext, deps Dependencies) ParseContext {
	user, createdOn, updatedOn, ok := deps.Alookup.StockAnnotator(
		ctx.Plasmid.Id,
	)
	if ok {
		ctx.Plasmid.User = user
		ctx.Plasmid.CreatedOn = createdOn
		ctx.Plasmid.UpdatedOn = updatedOn
	}
	return ctx
}

// applyPublications applies publication lookup to plasmid context
func applyPublications(ctx ParseContext, deps Dependencies) ParseContext {
	return F.Pipe2(
		deps.Plookup.StockPub(ctx.Plasmid.Id),
		O.FromPredicate(isNonEmptySlice),
		O.Fold(
			func() ParseContext { return ctx },
			func(pubs []string) ParseContext {
				filteredPubs := F.Pipe1(
					pubs,
					A.Filter(func(pub string) bool {
						return pub != ""
					}),
				)
				ctx.Plasmid.Publications = append(
					ctx.Plasmid.Publications,
					filteredPubs...)
				return ctx
			},
		),
	)
}

// applyGenes applies gene lookup to plasmid context
func applyGenes(ctx ParseContext, deps Dependencies) ParseContext {
	return F.Pipe2(
		deps.Glookup.StockGene(ctx.Plasmid.Id),
		O.FromPredicate(isNonEmptySlice),
		O.Fold(
			func() ParseContext { return ctx },
			func(genes []string) ParseContext {
				filteredGenes := F.Pipe1(
					genes,
					A.Filter(func(gene string) bool {
						return gene != ""
					}),
				)
				ctx.Plasmid.Genes = append(ctx.Plasmid.Genes, filteredGenes...)
				return ctx
			},
		),
	)
}

// enrichWithAnnotator enriches plasmid with annotator data using Reader
func enrichWithAnnotator(ctx ParseContext) PlasmidParser {
	return F.Pipe1(
		RE.Ask[Dependencies, error](),
		RE.Map[Dependencies, error](
			func(deps Dependencies) ParseContext {
				return applyAnnotator(ctx, deps)
			},
		),
	)
}

// isNonEmptySlice checks if a slice is non-empty
func isNonEmptySlice(slice []string) bool {
	return len(slice) > 0
}

// enrichWithPublications enriches plasmid with publications using Reader
func enrichWithPublications(ctx ParseContext) PlasmidParser {
	return F.Pipe1(
		RE.Ask[Dependencies, error](),
		RE.Map[Dependencies, error](func(deps Dependencies) ParseContext {
			return applyPublications(ctx, deps)
		}),
	)
}

// enrichWithGenes enriches plasmid with genes using Reader
func enrichWithGenes(ctx ParseContext) PlasmidParser {
	return F.Pipe1(
		RE.Ask[Dependencies, error](),
		RE.Map[Dependencies, error](
			func(deps Dependencies) ParseContext {
				return applyGenes(ctx, deps)
			},
		),
	)
}

// parseAllFields composes all field parsing functions
func parseAllFields(ctx ParseContext) PlasmidParser {
	return F.Pipe2(
		parseId(ctx),
		RE.Chain(parseName),
		RE.Chain(parseSummary),
	)
}

// Value gets a new Plasmid instance using fp-go Reader pattern
func (plr *csvPlasmidReader) Value() E.Either[error, *Plasmid] {
	// Check for CSV reader errors first
	if plr.Err != nil {
		return E.Left[*Plasmid](plr.Err)
	}

	// Create initial context
	ctx := ParseContext{
		Record:  plr.Record,
		Plasmid: &Plasmid{},
	}

	// Run the point-free pipeline
	result := F.Pipe4(
		parseAllFields(ctx),
		RE.Chain(enrichWithAnnotator),
		RE.Chain(enrichWithPublications),
		RE.Chain(enrichWithGenes),
		RE.Map[Dependencies, error](
			func(ctx ParseContext) *Plasmid {
				return ctx.Plasmid
			},
		),
	)

	return result(Dependencies{
		Alookup: plr.alookup,
		Plookup: plr.plookup,
		Glookup: plr.glookup,
	})
}
