package nosensitive

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/beer-connoisseur/log-checker/analysis/logcall"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nosensitive",
	Doc:      "checks that log messages do not contain sensitive data",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var wordsFlag string
var patterns []*regexp.Regexp

func init() {
	Analyzer.Flags.StringVar(
		&wordsFlag,
		"words",
		"password,pass,pwd,api_key,apikey,token,secret,key,auth,credential,credit,card,ssn",
		"comma-separated list of sensitive words",
	)
}

func run(pass *analysis.Pass) (any, error) {
	if patterns == nil {
		words := strings.Split(wordsFlag, ",")
		for _, w := range words {
			w = strings.TrimSpace(w)
			if w == "" {
				continue
			}
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(w) + `[:=]\s*`)
			patterns = append(patterns, re)
		}
	}
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	i.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		if !logcall.IsLoggerCall(pass, call) {
			return
		}
		for _, arg := range call.Args {
			if isConcatWithSensitive(arg) {
				pass.Reportf(arg.Pos(), "log message may contain sensitive data")
			}
		}
	})

	return nil, nil
}

func isConcatWithSensitive(expr ast.Expr) bool {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return false
	}

	return hasSensitiveLiteral(bin.X) && isValueExpr(bin.Y)
}

func isValueExpr(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.Ident, *ast.SelectorExpr, *ast.StarExpr:
		return true
	default:
		return false
	}
}

func hasSensitiveLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return false
	}
	lower := strings.ToLower(s)
	for _, re := range patterns {
		if re.MatchString(lower) {
			return true
		}
	}
	return false
}
