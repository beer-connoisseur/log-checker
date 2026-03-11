package lowercase

import (
	"go/ast"
	"log/slog"
	"strconv"
	"strings"
	"unicode"

	"github.com/beer-connoisseur/log-checker/analysis/logcall"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "lowercase",
	Doc:      "checks that log messages start with a lowercase letter",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	slog.Warn("    \nTest")
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		for _, lit := range logcall.FindMessageLiterals(pass, n) {
			s, err := strconv.Unquote(lit.Value)
			if err != nil {
				continue
			}
			if !isLowercaseStart(s) {
				pass.Report(analysis.Diagnostic{
					Pos:            lit.Pos(),
					Message:        "log message should start with a lowercase letter",
					SuggestedFixes: makeLowercaseFix(lit, s),
				})
			}
		}
	})

	return nil, nil
}

func isLowercaseStart(s string) bool {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if s == "" {
		return true
	}
	r := []rune(s)[0]
	if unicode.IsLetter(r) {
		return unicode.IsLower(r)
	}
	return true
}

func makeLowercaseFix(lit *ast.BasicLit, original string) []analysis.SuggestedFix {
	trimmed := strings.TrimLeftFunc(original, unicode.IsSpace)
	if trimmed == "" {
		return nil
	}

	runes := []rune(trimmed)
	if !unicode.IsLetter(runes[0]) {
		return nil
	}
	runes[0] = unicode.ToLower(runes[0])

	var whitespace strings.Builder
	for _, r := range original {
		if unicode.IsSpace(r) {
			whitespace.WriteRune(r)
		} else {
			break
		}
	}
	newText := whitespace.String() + string(runes)

	newLiteral := strconv.Quote(newText)

	return []analysis.SuggestedFix{{
		Message: "change first letter to lowercase",
		TextEdits: []analysis.TextEdit{{
			Pos:     lit.Pos(),
			End:     lit.End(),
			NewText: []byte(newLiteral),
		}},
	}}
}
