package english

import (
	"go/ast"
	"strconv"
	"unicode"

	"github.com/beer-connoisseur/log-checker/analysis/logcall"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "english",
	Doc:      "checks that log messages contain only english",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		for _, lit := range logcall.FindMessageLiterals(pass, n) {
			s, err := strconv.Unquote(lit.Value)
			if err != nil {
				continue
			}
			if !isEnglish(s) {
				pass.Reportf(lit.Pos(), "log message should be in english")
			}
		}
	})

	return nil, nil
}

func isEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && (r < 'A' || r > 'z' || (r > 'Z' && r < 'a')) {
			return false
		}
	}
	return true
}
