package nospecials

import (
	"go/ast"
	"strconv"
	"strings"
	"unicode"

	"github.com/beer-connoisseur/log-checker/analysis/logcall"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nospecials",
	Doc:      "checks that log messages do not contain special characters or emojis",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var allowedPunct = map[rune]bool{
	'.':  true,
	',':  true,
	':':  true,
	';':  true,
	'-':  true,
	'\'': true,
	'"':  true,
	'=':  true,
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
			if !hasNoSpecials(s) {
				pass.Report(analysis.Diagnostic{
					Pos:            lit.Pos(),
					Message:        "log message should not contain special characters or emojis",
					SuggestedFixes: makeNoSpecialsFix(lit, s),
				})
			}
		}
	})

	return nil, nil
}

func isEmoji(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F: // Emoticons
		return true
	case r >= 0x1F300 && r <= 0x1F5FF: // Miscellaneous Symbols and Pictographs
		return true
	case r >= 0x1F680 && r <= 0x1F6FF: // Transport and Map Symbols
		return true
	case r >= 0x1F700 && r <= 0x1F77F: // Alchemical Symbols
		return true
	case r >= 0x1F780 && r <= 0x1F7FF: // Geometric Shapes Extended
		return true
	case r >= 0x1F800 && r <= 0x1F8FF: // Supplemental Arrows-C
		return true
	case r >= 0x1F900 && r <= 0x1F9FF: // Supplemental Symbols and Pictographs
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F: // Chess Symbols
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF: // Symbols and Pictographs Extended-A
		return true
	case r >= 0x2702 && r <= 0x27B0: // Dingbats
		return true
	case r >= 0x24C2 && r <= 0x1F251: // Enclosed characters
		return true
	}
	return false
}

func hasNoSpecials(s string) bool {
	for _, r := range s {
		if allowedPunct[r] {
			continue
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || isEmoji(r) {
			return false
		}
	}
	return true
}

func makeNoSpecialsFix(lit *ast.BasicLit, original string) []analysis.SuggestedFix {
	var cleanText strings.Builder
	for _, r := range original {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || isEmoji(r) {
			continue
		} else {
			cleanText.WriteRune(r)
		}
	}

	newLiteral := strconv.Quote(cleanText.String())

	return []analysis.SuggestedFix{{
		Message: "delete special symbols and emojis",
		TextEdits: []analysis.TextEdit{{
			Pos:     lit.Pos(),
			End:     lit.End(),
			NewText: []byte(newLiteral),
		}},
	}}
}
