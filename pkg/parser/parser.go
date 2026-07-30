package parser

import (
	"slices"

	"github.com/azin-lang/Azin/pkg/ast"
	token2 "github.com/azin-lang/Azin/pkg/token"
)

type ErrorReporter interface {
	ReportError(pos token2.Position, length int, format string, args ...any)
	Err() error
}

const (
	_ int = iota
	PrecLowest
	PrecBitwiseAnd // &
	PrecEquality   // ==, !=
	PrecComparison // <, >, <=, >=
	PrecShift      // <<, >>
	PrecTerm       // +, -
	PrecFactor     // *, /
	PrecCall       // (
	PrecMember     // .
)

type Parser struct {
	source  string
	tokens  []token2.Token
	current int
	diag    ErrorReporter
}

func Parse(source string, tokens []token2.Token, diag ErrorReporter) (*ast.Program, error) {
	p := New(source, tokens, diag)
	return p.ParseProgram(), p.diag.Err()
}

func New(source string, tokens []token2.Token, diag ErrorReporter) *Parser {
	return &Parser{
		source: source,
		tokens: tokens,
		diag:   diag,
	}
}

func (p *Parser) Err() error {
	return p.diag.Err()
}

func (p *Parser) synchronize() {
	if p.isAtEnd() {
		return
	}
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == token2.Newline || p.previous().Kind == token2.Semicolon {
			return
		}

		if isSyncPoint(p.peek().Kind) {
			return
		}
		p.advance()
	}
}

func (p *Parser) lexeme(tok token2.Token) string {
	start := int(tok.Position.Offset)
	end := start + int(tok.Length)
	if start < 0 {
		return ""
	}
	if end > len(p.source) {
		end = len(p.source)
	}
	if start > end {
		return ""
	}
	return p.source[start:end]
}

func (p *Parser) reportError(tok token2.Token, format string, args ...any) {
	p.diag.ReportError(tok.Position, int(tok.Length), format, args...)
}

func badStmt(tok token2.Token) *ast.BadStmt {
	return &ast.BadStmt{Token: tok}
}

func badExpr(tok token2.Token) *ast.BadExpr {
	return &ast.BadExpr{Token: tok}
}

func isBadStmt(stmt ast.Stmt) bool {
	_, ok := stmt.(*ast.BadStmt)
	return ok
}

func isBadExpr(expr ast.Expr) bool {
	_, ok := expr.(*ast.BadExpr)
	return ok
}

func (p *Parser) expect(kind token2.Kind, context string) (token2.Token, bool) {
	if p.check(kind) {
		return p.advance(), true
	}

	got := p.peek()
	p.reportError(
		got,
		"expected %s %s, found %s",
		kind.DisplayName(),
		context,
		got.Kind.DisplayName(),
	)

	return got, false
}

func (p *Parser) parseBlock(until ...token2.Kind) []ast.Stmt {
	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.isAtEnd() || p.checkAny(until...) {
			break
		}
		if stmt := p.parseStatement(); stmt != nil {
			body = append(body, stmt)
		}
	}
	return body
}

func (p *Parser) skipNewlines() {
	for p.match(token2.Newline) {
	}
}

func (p *Parser) peek() token2.Token {
	if p.current >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() token2.Token {
	if p.current == 0 {
		return p.tokens[0]
	}
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.tokens[p.current].Kind == token2.EOF
}

func (p *Parser) advance() token2.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(kind token2.Kind) bool {
	return p.peek().Kind == kind
}

func (p *Parser) checkAny(kinds ...token2.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	return slices.Contains(kinds, p.peek().Kind)
}

func (p *Parser) match(kinds ...token2.Kind) bool {
	if p.isAtEnd() {
		return false
	}
	if slices.Contains(kinds, p.peek().Kind) {
		p.advance()
		return true
	}
	return false
}

//nolint:unparam
func (p *Parser) consumeStatementEnd() bool {
	switch p.peek().Kind {
	case token2.Semicolon:
		p.advance()
		return true
	case token2.Newline:
		p.skipNewlines()
		return true
	case token2.EOF, token2.KwEnd, token2.KwElse:
		return true
	default:
		p.reportError(p.peek(), "expected end of statement (newline or ';')")
		return false
	}
}

func isBuiltinType(kind token2.Kind) bool {
	switch kind {
	case token2.KwUnit, token2.KwInt, token2.KwFloat, token2.KwString, token2.KwChar, token2.KwBool:
		return true
	default:
		return false
	}
}

func isSyncPoint(kind token2.Kind) bool {
	switch kind {
	case token2.KwFn, token2.KwStruct, token2.KwEnum, token2.KwVar, token2.KwIf, token2.KwLoop,
		token2.KwReturn, token2.KwElse, token2.KwImportC, token2.KwImport, token2.KwDefer, token2.KwEnd:
		return true
	default:
		return false
	}
}

func getPrecedence(kind token2.Kind) int {
	switch kind {
	case token2.LeftParen:
		return PrecCall
	case token2.Dot:
		return PrecMember
	case token2.Star, token2.Slash:
		return PrecFactor
	case token2.Plus, token2.Minus:
		return PrecTerm
	case token2.LessLess, token2.GreaterGreater:
		return PrecShift
	case token2.Less, token2.LessEqual, token2.Greater, token2.GreaterEqual:
		return PrecComparison
	case token2.EqualEqual, token2.BangEqual:
		return PrecEquality
	case token2.Ampersand:
		return PrecBitwiseAnd
	default:
		return PrecLowest
	}
}
