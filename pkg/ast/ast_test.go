package ast_test

import (
	"strings"
	"testing"

	token2 "github.com/azin-lang/Azin/pkg/token"
)

func tok(kind token2.Kind, offset, length uint32) token2.Token {
	return token2.Token{Kind: kind, Position: token2.Position{Offset: offset}, Length: length}
}

func ident(value string) *Identifier {
	return &Identifier{
		Token: tok(token2.Identifier, 0, uint32(len(value))), //nolint:gosec
		Value: value,
	}
}

func TestProgram(t *testing.T) {
	p := &Program{Statements: nil}
	if p.TokenLiteral() != "" {
		t.Errorf("empty program TokenLiteral = %q", p.TokenLiteral())
	}
	if p.Pos() != (token2.Position{}) {
		t.Errorf("empty program Pos = %v", p.Pos())
	}
	if p.Label() != "Program" {
		t.Errorf("Label = %q", p.Label())
	}
}

func TestBadNodes(t *testing.T) {
	be := &BadExpr{Token: tok(token2.Error, 0, 1)}
	if be.Label() != "BadExpr" {
		t.Errorf("BadExpr Label = %q", be.Label())
	}

	bs := &BadStmt{Token: tok(token2.Error, 1, 2)}
	if bs.Label() != "BadStmt" {
		t.Errorf("BadStmt Label = %q", bs.Label())
	}
	if bs.Pos() != (token2.Position{Offset: 1}) {
		t.Errorf("BadStmt Pos = %v", bs.Pos())
	}
}

func TestVarStmt(t *testing.T) {
	v := &VarStmt{
		Token:   tok(token2.KwVar, 0, 3),
		Name:    ident("x"),
		SynType: ident("int"),
		Mutable: true,
	}
	if !strings.Contains(v.Label(), "var") {
		t.Errorf("Label missing 'var': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "mut") {
		t.Errorf("Label missing 'mut': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "x") {
		t.Errorf("Label missing 'x': %q", v.Label())
	}
	if !strings.Contains(v.Label(), "int") {
		t.Errorf("Label missing 'int': %q", v.Label())
	}
}

func TestFuncStmt(t *testing.T) {
	f := &FuncStmt{
		Token:         tok(token2.KwFn, 0, 2),
		Name:          ident("add"),
		Params:        []*FieldDecl{{Name: ident("a"), SynType: ident("int")}},
		SynReturnType: ident("int"),
	}
	label := f.Label()
	if !strings.Contains(label, "add") {
		t.Errorf("Label missing 'add': %q", label)
	}
	if !strings.Contains(label, "int") {
		t.Errorf("Label missing return type: %q", label)
	}
	if !strings.Contains(label, "a") {
		t.Errorf("Label missing param: %q", label)
	}
}

func TestIfStmt(t *testing.T) {
	s := &IfStmt{
		Token:     tok(token2.KwIf, 0, 2),
		Condition: ident("true"),
	}
	if s.Label() != "if" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestLoopStmt(t *testing.T) {
	s := &LoopStmt{
		Token: tok(token2.KwLoop, 0, 4),
	}
	if s.Label() != "loop" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestReturnStmt(t *testing.T) {
	s := &ReturnStmt{
		Token: tok(token2.KwReturn, 0, 6),
	}
	if s.Label() != "return" {
		t.Errorf("Label = %q", s.Label())
	}
}

func TestStructStmt(t *testing.T) {
	s := &StructStmt{
		Token: tok(token2.KwStruct, 0, 6),
		Name:  ident("Point"),
	}
	if s.Label() != "struct Point" {
		t.Errorf("Label = %q, want 'struct Point'", s.Label())
	}
}

func TestIdentExpr(t *testing.T) {
	id := ident("foobar")
	if id.Label() != "foobar" {
		t.Errorf("Label = %q", id.Label())
	}
	if id.TokenLiteral() != "foobar" {
		t.Errorf("TokenLiteral = %q", id.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	lit := &IntegerLiteral{Token: tok(token2.IntegerLiteral, 0, 2), Value: 42}
	if lit.Label() != "42" {
		t.Errorf("Label = %q", lit.Label())
	}
}

func TestFloatLiteral(t *testing.T) {
	lit := &FloatLiteral{Token: tok(token2.FloatLiteral, 0, 4), Value: 3.14}
	if lit.Label() != "3.14" {
		t.Errorf("Label = %q", lit.Label())
	}
}

func TestStringLiteral(t *testing.T) {
	lit := &StringLiteral{Token: tok(token2.StringLiteral, 0, 5), Value: "hello"}
	if !strings.Contains(lit.Label(), "hello") {
		t.Errorf("Label missing 'hello': %q", lit.Label())
	}
}

func TestCharacterLiteral(t *testing.T) {
	lit := &CharacterLiteral{Token: tok(token2.CharacterLiteral, 0, 3), Value: 'x'}
	if !strings.Contains(lit.Label(), "x") {
		t.Errorf("Label missing 'x': %q", lit.Label())
	}
}

func TestBooleanLiteral(t *testing.T) {
	tLit := &BooleanLiteral{Token: tok(token2.Identifier, 0, 4), Value: true}
	if tLit.Label() != "true" {
		t.Errorf("true literal Label = %q", tLit.Label())
	}
	fLit := &BooleanLiteral{Token: tok(token2.Identifier, 0, 5), Value: false}
	if fLit.Label() != "false" {
		t.Errorf("false literal Label = %q", fLit.Label())
	}
}

func TestCallExpr(t *testing.T) {
	call := &CallExpr{
		Callee: ident("foo"),
		Args:   []Expr{},
	}
	if call.Label() != "call foo" {
		t.Errorf("Label = %q", call.Label())
	}

	callMember := &CallExpr{
		Callee: &MemberExpr{
			Object:   ident("obj"),
			Property: ident("method"),
		},
	}
	if !strings.Contains(callMember.Label(), "call") {
		t.Errorf("Label missing 'call': %q", callMember.Label())
	}
}

func TestMemberExpr(t *testing.T) {
	m := &MemberExpr{
		Object:   ident("point"),
		Property: ident("x"),
	}
	if m.Label() != "point.x" {
		t.Errorf("Label = %q, want 'point.x'", m.Label())
	}
}

func TestBinaryExpr(t *testing.T) {
	b := &BinaryExpr{
		Left:     &IntegerLiteral{Value: 1},
		Operator: tok(token2.Plus, 0, 1),
		Right:    &IntegerLiteral{Value: 2},
	}
	// Label returns the Kind.String(), e.g. "plus" not "+"
	if b.Label() == "" {
		t.Errorf("Label = empty")
	}
}

func TestFieldDecl(t *testing.T) {
	f := &FieldDecl{
		Name:    ident("name"),
		SynType: ident("string"),
	}
	if !strings.Contains(f.Label(), "name") {
		t.Errorf("Label missing 'name': %q", f.Label())
	}
	if !strings.Contains(f.Label(), "string") {
		t.Errorf("Label missing 'string': %q", f.Label())
	}
}

func TestAssignmentStmt(t *testing.T) {
	a := &AssignmentStmt{
		Token: tok(token2.Equal, 0, 1),
		Left:  ident("x"),
		Value: &IntegerLiteral{Value: 42},
	}
	if a.Label() != "assign" {
		t.Errorf("Label = %q, want 'assign'", a.Label())
	}
}

func TestImportCStmt(t *testing.T) {
	i := &ImportCStmt{
		Token: tok(token2.KwImportC, 0, 7),
		Path:  &StringLiteral{Value: "stdio.h"},
	}
	if !strings.Contains(i.Label(), "stdio.h") {
		t.Errorf("Label missing stdio.h: %q", i.Label())
	}
}

func TestExpressionStmt(t *testing.T) {
	e := &ExpressionStmt{
		Token:      tok(token2.Identifier, 0, 4),
		Expression: ident("test"),
	}
	if e.Label() != "test" {
		t.Errorf("Label = %q, want 'test'", e.Label())
	}

	nilExpr := &ExpressionStmt{}
	if nilExpr.Label() != "expr" {
		t.Errorf("nil expr Label = %q, want 'expr'", nilExpr.Label())
	}
}
