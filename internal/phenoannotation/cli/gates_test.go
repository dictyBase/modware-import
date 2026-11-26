package cli

import (
	"strings"
	"testing"

	"github.com/dictyBase/fp-go-loom/loomcheck"
	pc "github.com/dictyBase/fp-go-loom/pipelinecheck"
)

// roots lists the fp-go packages of the phenoannotation command. The
// cmd/phenoannotation main package is a plain urfave/cli definition and
// is not checked.
var roots = []string{".", "../normalize"}

func isEntrypoint(name string) bool {
	return strings.HasPrefix(name, "run") || strings.HasPrefix(name, "Run")
}

// TestFunctionalGates runs every pipelinecheck rule and loomcheck on the
// two packages. The checks are syntax only. A pass is not a proof.
func TestFunctionalGates(t *testing.T) {
	var cfg pc.Config
	cfg.Roots = roots
	cfg.IsEntrypoint = isEntrypoint
	cfg.RequirePointFreeSeed = true
	cfg.RequirePointFreeBranching = true
	cfg.RequireTrivialDelegate = true
	cfg.TerminalFunctions = []string{"Publish", "PrintSummary", "WriteFile"}
	pc.Require(t, cfg)
	pc.RequireBareIf(t, cfg)
	pc.RequireNoIfErrInTryCatch(t, cfg)
	pc.RequireNoNonRawTryCatchCallback(t, cfg)
	pc.RequireSafePrintf(t, cfg)
	pc.RequireNoRedundantFold(t, cfg)
	pc.RequireNoHandRolledFileJoinInArrayMap(t, cfg)
	pc.RequireNoImperativeLoop(t, cfg)
	pc.RequireNoImperativeBranch(t, cfg)
	pc.RequireNoDirectApply(t, cfg)
	pc.RequireNoNewPipe(t, cfg)
	pc.RequireNoBoolExpr(t, cfg)
	pc.RequireNoPopulatedStruct(t, cfg)
	pc.RequireNoLocalPipeClosure(t, cfg)
	pc.RequireNoFlap(t, cfg)
	pc.RequireTerminalEffect(t, cfg)
	var lc loomcheck.Config
	lc.Roots = roots
	loomcheck.Require(t, lc)
}
