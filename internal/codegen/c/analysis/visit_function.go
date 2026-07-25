package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitFunction(fn *ast.FuncStmt) {
	name := FunctionName(fn)
	if fn.SynReturnType != nil {
		a.MarkTypeUsed(fn.SynReturnType.Value)
	}
	for _, param := range fn.Params {
		if param.SynType != nil {
			a.MarkTypeUsed(param.SynType.Value)
			a.registerVariable(name, param.Name.Value)
			a.useVariable(name, param.Name.Value)
		}
	}
	for _, stmt := range fn.Body {
		a.visitStmt(name, stmt)
	}
}
