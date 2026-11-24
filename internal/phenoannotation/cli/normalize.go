package cli

import (
	"context"
	"errors"
	"io"

	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/predicate"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	cli "github.com/urfave/cli/v3"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

var (
	// ErrReview marks a complete batch that contains rejected rows.
	ErrReview = errors.New("batch contains rejected rows: review rejects.csv")
	errFlags  = errors.New("use one input flag and no positional arguments")
	errNoSub  = errors.New("expected normalize subcommand")

	argsL    = MakeargumentsLenses()
	stateL   = norm.MakeStateLenses()
	optionsL = norm.MakeOptionsLenses()
	batchL   = norm.MakeBatchLenses()

	stateRejected = F.Pipe1(stateL.Batch, L.Compose[norm.State](batchL.Rejected))

	cmdName     = L.MakeLensRef(getCmdName, setCmdName)
	cmdUsage    = L.MakeLensRef(getCmdUsage, setCmdUsage)
	cmdAction   = L.MakeLensRef(getCmdAction, setCmdAction)
	cmdFlags    = L.MakeLensRef(getCmdFlags, setCmdFlags)
	cmdWriter   = L.MakeLensRef(getCmdWriter, setCmdWriter)
	cmdCommands = L.MakeLensRef(getCmdCommands, setCmdCommands)

	flagName     = L.MakeLensRef(getFlagName, setFlagName)
	flagUsage    = L.MakeLensRef(getFlagUsage, setFlagUsage)
	flagRequired = L.MakeLensRef(getFlagRequired, setFlagRequired)

	// Flag shape: no positional argument and not both input flags. An
	// explicit empty value still counts as a set flag.
	flagSet      = (*cli.Command).IsSet
	noPositional = F.Flow2((*cli.Command).NArg, P.IsZero[int]())
	bothSources  = F.Pipe1(F.Bind2nd(flagSet, "folder"), P.And(F.Bind2nd(flagSet, "input")))
	plainCall    = F.Pipe1(noPositional, P.And(P.Not(bothSources)))

	requireFlagShape = E.FromPredicate(
		F.Pipe1(plainCall, P.ContraMap(argsL.Cmd.Get)),
		F.Constant1[arguments](errFlags),
	)

	// requireNoRejects turns a batch with rejected rows into ErrReview.
	// The batch is already published when this step runs.
	requireNoRejects = E.FromPredicate(
		F.Flow2(stateRejected.Get, P.IsZero[int]()),
		F.Constant1[norm.State](ErrReview),
	)

	flagText = (*cli.Command).String

	// toOptions copies the flag values into the options.
	toOptions = endoutils.BuildFrom([]func(*cli.Command) En.Endomorphism[norm.Options]{
		F.Flow2(F.Bind2nd(flagText, "input"), optionsL.Input.Set),
		F.Flow2(F.Bind2nd(flagText, "folder"), optionsL.Folder.Set),
		F.Flow2(F.Bind2nd(flagText, "sheet"), optionsL.Sheet.Set),
		F.Flow2(F.Bind2nd(flagText, "out"), optionsL.Out.Set),
	})

	// toState builds a new normalize state from one invocation.
	toState = endoutils.BuildFrom([]func(arguments) En.Endomorphism[norm.State]{
		F.Flow2(argsL.Ctx.Get, stateL.Context.Set),
		F.Flow3(argsL.Cmd.Get, toOptions, stateL.Options.Set),
	})
)

func getCmdName(c *cli.Command) string                 { return c.Name }
func setCmdName(c *cli.Command, v string) *cli.Command { c.Name = v; return c }

func getCmdUsage(c *cli.Command) string                 { return c.Usage }
func setCmdUsage(c *cli.Command, v string) *cli.Command { c.Usage = v; return c }

func getCmdAction(c *cli.Command) cli.ActionFunc { return c.Action }

func setCmdAction(c *cli.Command, v cli.ActionFunc) *cli.Command { c.Action = v; return c }

func getCmdFlags(c *cli.Command) []cli.Flag                 { return c.Flags }
func setCmdFlags(c *cli.Command, v []cli.Flag) *cli.Command { c.Flags = v; return c }

func getCmdWriter(c *cli.Command) io.Writer                 { return c.Writer }
func setCmdWriter(c *cli.Command, v io.Writer) *cli.Command { c.Writer = v; return c }

func getCmdCommands(c *cli.Command) []*cli.Command { return c.Commands }

func setCmdCommands(c *cli.Command, v []*cli.Command) *cli.Command {
	c.Commands = v
	return c
}

func getFlagName(f *cli.StringFlag) string                     { return f.Name }
func setFlagName(f *cli.StringFlag, v string) *cli.StringFlag  { f.Name = v; return f }
func getFlagUsage(f *cli.StringFlag) string                    { return f.Usage }
func setFlagUsage(f *cli.StringFlag, v string) *cli.StringFlag { f.Usage = v; return f }
func getFlagRequired(f *cli.StringFlag) bool                   { return f.Required }

func setFlagRequired(f *cli.StringFlag, v bool) *cli.StringFlag {
	f.Required = v
	return f
}

// runNormalize validates the flags, selects files, reads and classifies
// every row, publishes it, prints the summary, and returns the
// review status. It forces the effect once.
func runNormalize(ctx context.Context, c *cli.Command) error {
	return F.Pipe13(
		arguments{},
		argsL.Ctx.Set(ctx),
		argsL.Cmd.Set(c),
		IOE.Of[error, arguments],
		IOE.ChainEitherK(requireFlagShape),
		IOE.Map[error](F.Flow2(toState, stateL.Writer.Set(norm.WriteFile))),
		IOE.ChainEitherK(norm.Valid),
		IOE.Bind(stateL.Files.Set, norm.SelectFiles),
		IOE.Bind(stateL.Batch.Set, norm.ReadFiles),
		IOE.Chain(norm.Publish),
		IOE.ChainFirstIOK[error](norm.PrintSummary),
		IOE.ChainEitherK(requireNoRejects),
		ioeutils.ToEither[error, norm.State],
		E.ToError[norm.State],
	)
}

func needCommand(context.Context, *cli.Command) error { return errNoSub }

func textFlag(name, usage string, required bool) cli.Flag {
	return F.Pipe3(
		&cli.StringFlag{},
		flagName.Set(name),
		flagUsage.Set(usage),
		flagRequired.Set(required),
	)
}

func normalizeCommand() *cli.Command {
	return F.Pipe4(
		&cli.Command{},
		cmdName.Set("normalize"),
		cmdUsage.Set("convert create workbooks into reviewable CSV files"),
		cmdAction.Set(runNormalize),
		cmdFlags.Set([]cli.Flag{
			textFlag("input", "one create workbook", false),
			textFlag("folder", "folder of create workbooks", false),
			textFlag("out", "new output directory", true),
			textFlag("sheet", "exact worksheet name", false),
		}),
	)
}

// Command returns the root command. It writes the summary to w.
func Command(w io.Writer) *cli.Command {
	return F.Pipe4(
		&cli.Command{},
		cmdName.Set("phenoannotation"),
		cmdAction.Set(needCommand),
		cmdWriter.Set(w),
		cmdCommands.Set([]*cli.Command{normalizeCommand()}),
	)
}
