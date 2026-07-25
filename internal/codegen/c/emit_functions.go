//nolint:goconst
package c

import (
	"github.com/azin-lang/Azin/internal/ast"
)

func (t *Transpiler) emitFunction(
	fn *ast.FuncStmt,
) {
	t.defers = nil

	t.emitFunctionSignature(fn)

	t.write(" {\n")

	t.pushIndent()

	for _, stmt := range fn.Body {
		t.emitStatement(stmt)
	}

	t.flushDefers()

	if t.functionName(fn) == "main" &&
		fn.SynReturnType != nil &&
		emitType(fn.SynReturnType.Value) == "void" {

		t.indentLine()
		t.write("return 0;\n")
	}

	t.popIndent()

	t.write("}\n")
}

func (t *Transpiler) emitFunctionSignature(
	fn *ast.FuncStmt,
) {
	name := t.functionName(fn)

	if fn.SynReturnType == nil {
		t.printf("void %s(", name)
		t.write("void")
		t.write(") /* missing return type */")
		return
	}

	ret := emitType(
		fn.SynReturnType.Value,
	)

	if name == "main" &&
		ret == "void" {

		ret = "int"
	}

	t.printf(
		"%s %s(",
		ret,
		name,
	)

	if len(fn.Params) == 0 {
		t.write("void")
	} else {
		for i, param := range fn.Params {

			if i > 0 {
				t.write(", ")
			}

			if param.SynType == nil {
				t.printf("void %s /* missing type */", param.Name.Value)
				continue
			}

			t.printf(
				"%s %s",
				emitType(param.SynType.Value),
				param.Name.Value,
			)
		}
	}

	t.write(")")
}

func (t *Transpiler) functionName(
	fn *ast.FuncStmt,
) string {
	if fn.CName != "" {
		return fn.CName
	}

	return fn.Name.Value
}
