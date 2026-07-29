//nolint:unused,unparam
package lexer

import (
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) eofToken() token2.Token {
	return token2.Token{
		Kind:     token2.EOF,
		Position: l.pos(),
	}
}

func (l *Lexer) eof() bool {
	return l.file.EOF(l.cursor)
}

func (l *Lexer) peek() rune {
	if l.eof() {
		return 0
	}
	r, _ := l.file.Rune(l.cursor)
	return r
}

func (l *Lexer) peekNext() rune {
	if l.eof() {
		return 0
	}

	_, size := l.file.Rune(l.cursor)
	nextCursor := l.cursor + size

	if l.file.EOF(nextCursor) {
		return 0
	}

	nextRune, _ := l.file.Rune(nextCursor)
	return nextRune
}

func (l *Lexer) advance() (r rune, size uint32) {
	if l.eof() {
		return 0, 0
	}
	r, size = l.file.Rune(l.cursor)
	l.cursor += size
	return r, size
}

func (l *Lexer) match(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	_, _ = l.advance()
	return true
}

func (l *Lexer) matchAny(chars string) bool {
	r := l.peek()
	for _, ch := range chars {
		if r == ch {
			_, _ = l.advance()
			return true
		}
	}
	return false
}

func (l *Lexer) consumeWhile(pred func(rune) bool) {
	for pred(l.peek()) {
		_, _ = l.advance()
	}
}

func (l *Lexer) emit(kind token2.Kind, start token2.Position) token2.Token {
	return token2.Token{
		Kind:     kind,
		Position: start,
		Length:   l.cursor - start.Offset,
	}
}

func (l *Lexer) either(ch rune, ifMatch, otherwise token2.Kind, start token2.Position) token2.Token {
	if l.match(ch) {
		return l.emit(ifMatch, start)
	}
	return l.emit(otherwise, start)
}

func (l *Lexer) pos() token2.Position {
	return token2.Position{Offset: l.cursor}
}
