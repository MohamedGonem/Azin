//nolint:unused
package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitStruct(s *ast.StructStmt) {
	for _, field := range s.Fields {
		if field.SynType != nil {
			a.MarkTypeUsed(field.SynType.Value)
		}
	}
}
