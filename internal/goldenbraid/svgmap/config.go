// Package svgmap renders circular poly records as accessible SVG bytes.
package svgmap

import (
	"errors"

	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/bebop/poly/io/genbank"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

const (
	// DefaultSizePx is the default square viewport size.
	DefaultSizePx = 640
	// MinLabelGapPx is the minimum vertical label separation.
	MinLabelGapPx = 16
	// LabelMarginPx keeps labels away from the viewport edges.
	LabelMarginPx = 16
	// ArrowAngleDeg is the maximum angular arrowhead length.
	ArrowAngleDeg = 2.5
	// DefaultMaxLabels caps the default number of visible feature labels.
	DefaultMaxLabels   = 32
	radiusFraction     = 0.28
	ringStepPx         = 12.0
	featureHalfWidthPx = 4.0
	labelLeaderGapPx   = 8.0
	labelFontSizePx    = 10.0
	glyphWidthPx       = 6.0
	half               = 2.0
	fullTurn           = 2.0
	degreesPerHalfTurn = 180.0
	defaultFill        = "#88aabb"
	backboneStroke     = "#495568"
	labelStroke        = "#8893a3"
	leftSide           = "left"
	rightSide          = "right"
)

var (
	// ErrInvalidConfig identifies invalid size, label limit, or ring geometry.
	ErrInvalidConfig = errors.New("invalid SVG map configuration or geometry")
	// ErrNotCircular identifies a linear record passed to the circular renderer.
	ErrNotCircular = errors.New("SVG map requires a circular record")
	configL        = MakeConfigLenses()
	inputL         = MakeRenderInputLenses()
	stateL         = MakerenderStateLenses()
	featureL       = MakefeatureLayoutLenses()
	defaultConfig  = F.Pipe2(
		Config{},
		configL.Size.Set(DefaultSizePx),
		configL.MaxLabels.Set(DefaultMaxLabels),
	)
)

// Config controls viewport size, primer visibility, and visible label count.
// fp-go:Lens
type Config struct {
	Size        int
	ShowPrimers bool
	MaxLabels   int
}

// Option updates a value copy of Config.
type Option = En.Endomorphism[Config]

// RenderInput supplies the single shared poly record and rendering options.
// fp-go:Lens
type RenderInput struct {
	Record  genbank.Genbank
	Options []Option
}

// fp-go:Lens
type renderState struct {
	Record   genbank.Genbank
	Config   Config
	Length   int
	Features []featureLayout
	Labels   []featureLayout
}

// fp-go:Lens
type featureLayout struct {
	ID      int
	Feature genbank.Feature
	Leaves  []genbank.Location
	Ring    int
	Radius  float64
	Angle   float64
	X       float64
	Y       float64
	Side    string
}

// WithSize sets the square viewport size in pixels.
func WithSize(px int) Option { return configL.Size.Set(px) }

// WithShowPrimers includes primer sites in paths and visible labels.
func WithShowPrimers(show bool) Option { return configL.ShowPrimers.Set(show) }

// WithMaxLabels limits visible labels; zero hides all labels, not feature paths.
func WithMaxLabels(limit int) Option { return configL.MaxLabels.Set(limit) }
