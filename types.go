package sundrylint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func IsVar(pass *analysis.Pass, e ast.Expr) bool {
	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() && tv.Assignable() {
		return true
	}
	return false
}

func IsConst(pass *analysis.Pass, e ast.Expr) bool {
	switch e.(type) {
	case *ast.SelectorExpr, *ast.BasicLit, *ast.Ident:
	case *ast.TypeAssertExpr, *ast.CallExpr, *ast.IndexExpr, *ast.StarExpr:
		return false
	default:
		return false
	}

	tv := pass.TypesInfo.Types[e]
	if tv.Addressable() || tv.Assignable() {
		return false
	}
	return true
}

func IsPkg(pass *analysis.Pass, e ast.Expr, pkgPath string) bool {
	if selectExpr, ok := e.(*ast.SelectorExpr); ok {
		obj := pass.TypesInfo.ObjectOf(selectExpr.Sel)
		if obj == nil || obj.Pkg().Path() != pkgPath {
			return false
		}
		return true
	}
	return false
}

type FuncType struct {
	ArgsNum    int
	Signature  string
	ResultsNum int
}

func IsFunc(pass *analysis.Pass, node *ast.CallExpr, fnType FuncType) bool {
	if len(node.Args) != fnType.ArgsNum {
		return false
	}

	sign, ok := pass.TypesInfo.TypeOf(node.Fun).(*types.Signature)
	if !ok {
		return false
	}
	if sign.String() != fnType.Signature {
		return false
	}
	if sign.Params().Len() != fnType.ArgsNum {
		return false
	}
	if sign.Results().Len() != fnType.ResultsNum {
		return false
	}
	return true
}

// restOfBlock, given a traversal stack, finds the innermost containing
// block and returns the suffix of its statements starting with the current
// node, along with the number of call expressions encountered.
func restOfBlock(stack []ast.Node) ([]ast.Stmt, int) {
	var ncalls int
	for i := len(stack) - 1; i >= 0; i-- {
		if b, ok := stack[i].(*ast.BlockStmt); ok {
			for j, v := range b.List {
				if v == stack[i+1] {
					return b.List[j:], ncalls
				}
			}
			break
		}

		if _, ok := stack[i].(*ast.CallExpr); ok {
			ncalls++
		}
	}
	return nil, 0
}
