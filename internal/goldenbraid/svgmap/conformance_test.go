package svgmap_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// A local projection followed by a returned pipeline escapes the existing gates.
func TestRendererFunctionsRemainSingleExpression(t *testing.T) {
	for _, name := range []string{"config.go", "layout.go", "render.go"} {
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		require.NoError(t, err)
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && splitPipeline(fn) {
				t.Errorf("%s %s: local derivations before terminal pipeline", name, fn.Name.Name)
			}
		}
	}
}

func splitPipeline(fn *ast.FuncDecl) bool {
	if fn.Body == nil || len(fn.Body.List) < 2 {
		return false
	}
	ret, ok := fn.Body.List[len(fn.Body.List)-1].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	call, ok := ret.Results[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && strings.HasPrefix(sel.Sel.Name, "Pipe")
}
