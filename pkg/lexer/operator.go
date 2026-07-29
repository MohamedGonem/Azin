package lexer

import (
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) lexOperator(ch rune, start token2.Position) token2.Token {
	switch ch {
	case '+':
		return l.lexPlus(start)
	case '-':
		return l.lexMinus(start)
	case '*':
		return l.either('=', token2.StarEqual, token2.Star, start)
	case '/':
		return l.either('=', token2.SlashEqual, token2.Slash, start)
	case '%':
		return l.either('=', token2.ModuloEqual, token2.Modulo, start)
	case '=':
		return l.either('=', token2.EqualEqual, token2.Equal, start)
	case '!':
		return l.either('=', token2.BangEqual, token2.Bang, start)
	case '<':
		if l.match('=') {
			return l.emit(token2.LessEqual, start)
		}
		if l.match('<') {
			return l.emit(token2.LessLess, start)
		}
		return l.emit(token2.Less, start)
	case '>':
		if l.match('=') {
			return l.emit(token2.GreaterEqual, start)
		}
		if l.match('>') {
			return l.emit(token2.GreaterGreater, start)
		}
		return l.emit(token2.Greater, start)
	case '&':
		if l.match('&') {
			return l.emit(token2.LogicalAnd, start)
		}
		if l.match('=') {
			return l.emit(token2.AmpersandEqual, start)
		}
		return l.emit(token2.Ampersand, start)
	case '|':
		if l.match('|') {
			return l.emit(token2.LogicalOr, start)
		}
		if l.match('=') {
			return l.emit(token2.PipeEqual, start)
		}
		return l.emit(token2.Pipe, start)
	case '"':
		return l.lexString(start)
	default:
		return l.lexUnknown(start)
	}
}

func (l *Lexer) lexPlus(start token2.Position) token2.Token {
	if l.match('=') {
		return l.emit(token2.PlusEqual, start)
	}
	if l.match('+') {
		return l.emit(token2.PlusPlus, start)
	}
	return l.emit(token2.Plus, start)
}

func (l *Lexer) lexMinus(start token2.Position) token2.Token {
	if l.match('=') {
		return l.emit(token2.MinusEqual, start)
	}
	if l.match('-') {
		return l.emit(token2.MinusMinus, start)
	}
	if l.match('>') {
		return l.emit(token2.Arrow, start)
	}
	return l.emit(token2.Minus, start)
}

func (l *Lexer) lexUnknown(start token2.Position) token2.Token {
	l.consumeWhile(func(r rune) bool {
		// Stop on EOF
		if r == 0 {
			return false
		}

		// Stop on whitespace
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return false
		}

		// Stop on any character that could start a valid token
		if isIdentifierStart(r) || isDigit(r) || isPunctuation(r) {
			return false
		}

		// Stop on valid operator characters and quotes
		switch r {
		case '+', '-', '*', '/', '%', '=', '!', '<', '>', '&', '|', '"':
			return false
		}

		// Otherwise, it's more garbage. Keep eating it!
		return true
	})

	length := l.cursor - start.Offset
	text := string(l.file.Slice(start.Offset, l.cursor))

	l.diag.ReportError(start, int(length), "unexpected characters: %q", text)
	return l.emit(token2.Unknown, start)
}
