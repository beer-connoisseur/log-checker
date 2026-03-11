package logcall

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func IsLoggerCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	var obj types.Object
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		obj = pass.TypesInfo.ObjectOf(fun)
	case *ast.SelectorExpr:
		obj = pass.TypesInfo.ObjectOf(fun.Sel)
	default:
		return false
	}
	if obj == nil {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	path := pkg.Path()
	if path == "log/slog" || path == "go.uber.org/zap" {
		name := obj.Name()
		switch name {
		case "Info", "Error", "Warn", "Debug":
			return true
		}
	}
	return false
}

func ExtractStrings(expr ast.Expr) []*ast.BasicLit {
	var lits []*ast.BasicLit
	ast.Inspect(expr, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				lits = append(lits, node)
			}
			return false
		case *ast.BinaryExpr:
			return node.Op == token.ADD
		case *ast.ParenExpr:
			return true
		default:
			return false
		}
	})
	return lits
}

func FindMessageLiterals(pass *analysis.Pass, node ast.Node) []*ast.BasicLit {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	if !IsLoggerCall(pass, call) {
		return nil
	}
	var all []*ast.BasicLit
	for _, arg := range call.Args {
		all = append(all, ExtractStrings(arg)...)
	}
	return all
}
