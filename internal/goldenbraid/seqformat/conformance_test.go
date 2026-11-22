package seqformat_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// A local projection followed by a return used to escape the existing gates.
func TestConvertFunctionsRemainSingleExpression(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "convert.go", nil, 0)
	require.NoError(t, err)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && splitPipeline(fn) {
			t.Errorf("convert.go %s: local derivations before terminal pipeline", fn.Name.Name)
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
