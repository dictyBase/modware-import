package gbcli_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRunExportRemainsSinglePipeline guards against local derivations
// splitting the CLI entrypoint before its terminal F.Pipe* expression.
func TestRunExportRemainsSinglePipeline(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "export.go", nil, 0)
	require.NoError(t, err)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "RunExport" {
			continue
		}
		require.NotNil(t, fn.Body)
		require.Len(t, fn.Body.List, 1, "RunExport must stay a single pipeline return")
		ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
		require.True(t, ok)
		require.Len(t, ret.Results, 1)
		call, ok := ret.Results[0].(*ast.CallExpr)
		require.True(t, ok)
		sel, ok := call.Fun.(*ast.SelectorExpr)
		require.True(t, ok && strings.HasPrefix(sel.Sel.Name, "Pipe"))
	}
}
