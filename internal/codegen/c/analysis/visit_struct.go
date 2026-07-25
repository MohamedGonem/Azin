//nolint:unused
package analysis

import "github.com/azin-lang/Azin/internal/ast"

func (a *Analyzer) visitStruct(s *ast.StructStmt) {
	for _, field := range s.Fields {
		if field.SemaType.IsNominal() {
			a.MarkTypeUsed(field.SemaType)
		}
	}
}
