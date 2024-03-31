package sundrylint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const appendFuncName = `append`

var appendFunc = FuncType{
	ArgsNum:     2,
	Signature:   ``,
	ResultsNum:  1,
	Variadic:    true,
	UseVariadic: true,
}

func LintRangeAppendAll(pass *analysis.Pass, node *ast.CallExpr, stack []ast.Node) (ds []analysis.Diagnostic) {
	if !(IsFuncType(pass, node, appendFunc) && IsBuiltinFunc(pass, node, appendFuncName)) {
		return nil
	}

	appendIdent, ok := node.Args[1].(*ast.Ident)
	if !ok {
		return nil
	}
	blockStmt, rangeStmt := fetchBlockAndRangeStmt(stack)
	if blockStmt == nil || rangeStmt == nil {
		return nil
	}
	rangeIdent, ok := rangeStmt.X.(*ast.Ident)
	if !ok {
		return nil
	}
	appendObj := pass.TypesInfo.ObjectOf(appendIdent)
	rangeObj := pass.TypesInfo.ObjectOf(rangeIdent)
	if appendObj == nil || appendObj != rangeObj {
		return nil
	}
	return []analysis.Diagnostic{
		{
			Pos:      node.Pos(),
			End:      node.End(),
			Category: LinterName,
			Message:  SubLinterRangeappendallMessage,
		},
	}
}

func fetchBlockAndRangeStmt(stack []ast.Node) (*ast.BlockStmt, *ast.RangeStmt) {
	if len(stack) < 2 {
		return nil, nil
	}
	var block *ast.BlockStmt
	for i := len(stack) - 2; i >= 0; i-- {
		switch node := stack[i].(type) {
		case *ast.BlockStmt:
			if block == nil {
				block = node
			}
		case *ast.RangeStmt:
			return block, node
		case *ast.CallExpr:
			return block, nil
		default:
		}
	}
	return block, nil
}
