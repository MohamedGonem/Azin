package lexer

import (
	token2 "github.com/azin-lang/Azin/pkg/token"
)

func (l *Lexer) lexIdentifier(start token2.Position) token2.Token {
	l.consumeWhile(isIdentifierContinue)

	if kind, ok := token2.Keywords[string(l.file.Slice(start.Offset, l.cursor))]; ok {
		return l.emit(kind, start)
	}

	return l.emit(token2.Identifier, start)
}
