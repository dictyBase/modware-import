package gbexport_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExportAvoidsEndomorphismBuildersAndNonShortCircuitTraverse keeps
// struct construction direct (no endomorphism/BuildFrom tables) and prevents
// non-short-circuiting IOE.TraverseArraySeq from re-entering fallible I/O.
func TestExportAvoidsEndomorphismBuildersAndNonShortCircuitTraverse(t *testing.T) {
	fset := token.NewFileSet()
	for _, name := range []string{"state.go", "source.go", "effects.go", "pipeline.go"} {
		file, err := parser.ParseFile(fset, name, nil, 0)
		require.NoError(t, err)
		for _, imp := range file.Imports {
			path, unquoteErr := strconv.Unquote(imp.Path.Value)
			require.NoError(t, unquoteErr)
			if path == "github.com/IBM/fp-go/v2/endomorphism" ||
				path == "github.com/dictyBase/fp-go-loom/endomorphism" {
				t.Errorf("%s imports %q; use direct lens setters instead", name, path)
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if ok && sel.Sel.Name == "TraverseArraySeq" {
				if id, isIdent := sel.X.(*ast.Ident); isIdent && id.Name == "IOE" {
					t.Errorf("%s:%d uses IOE.TraverseArraySeq (does not short-circuit effects)",
						name, fset.Position(sel.Pos()).Line)
				}
			}
			return true
		})
	}
}
