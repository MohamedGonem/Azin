package lexer

import (
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) lexPunctuation(ch rune, start token2.Position) token2.Token {
	switch ch {
	case '(':
		return l.emit(token2.LeftParen, start)
	case ')':
		return l.emit(token2.RightParen, start)
	case '{':
		return l.emit(token2.LeftBrace, start)
	case '}':
		return l.emit(token2.RightBrace, start)
	case '[':
		return l.emit(token2.LeftBracket, start)
	case ']':
		return l.emit(token2.RightBracket, start)
	case ',':
		return l.emit(token2.Comma, start)
	case ';':
		return l.emit(token2.Semicolon, start)
	case ':':
		return l.emit(token2.Colon, start)
	case '.':
		return l.emit(token2.Dot, start)
	}

	// Should be unreachable as long as isPunctuation guards this method.
	l.diag.ReportError(start, 1, "internal error: unexpected character '%c'", ch)
	return l.emit(token2.Error, start)
}
