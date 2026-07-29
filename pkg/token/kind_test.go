package token_test

import (
	"testing"
)

func TestKindDisplayName(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{Unknown, "unknown"},
		{Identifier, "identifier"},
		{IntegerLiteral, "integer literal"},
		{FloatLiteral, "float literal"},
		{StringLiteral, "string literal"},
		{CharacterLiteral, "character literal"},
		{KwFn, "'fn'"},
		{KwDo, "'do'"},
		{KwVar, "'var'"},
		{KwMut, "'mut'"},
		{KwReturn, "'return'"},
		{KwEnd, "'end'"},
		{KwIf, "'if'"},
		{KwThen, "'then'"},
		{KwElse, "'else'"},
		{KwStruct, "'struct'"},
		{KwIs, "'is'"},
		{KwImportC, "'importC'"},
		{KwChar, "'char'"},
		{KwInt, "'int'"},
		{KwBool, "'bool'"},
		{KwNull, "'null'"},
		{KwUnit, "'unit'"},
		{KwString, "'string'"},
		{KwFloat, "'float'"},
		{Plus, "'+'"},
		{Minus, "'-'"},
		{Star, "'*'"},
		{Slash, "'/'"},
		{Equal, "'='"},
		{EqualEqual, "'=='"},
		{Bang, "'!'"},
		{BangEqual, "'!='"},
		{Less, "'<'"},
		{LessEqual, "'<='"},
		{Greater, "'>'"},
		{GreaterEqual, "'>='"},
		{LeftParen, "'('"},
		{RightParen, "')'"},
		{Comma, "','"},
		{Colon, "':'"},
		{Semicolon, "';'"},
		{Dot, "'.'"},
		{Newline, "newline"},
		{EOF, "end of file"},
	}

	for _, tt := range tests {
		got := tt.kind.DisplayName()
		if got != tt.want {
			t.Errorf("DisplayName(%d) = %q, want %q", tt.kind, got, tt.want)
		}
	}
}

func TestKindDisplayNameNonEmpty(t *testing.T) {
	for k := Unknown; k <= Error; k++ {
		name := k.DisplayName()
		if name == "" {
			t.Errorf("DisplayName(%d) is empty", k)
		}
	}
}

func TestKindHasText(t *testing.T) {
	tests := []struct {
		kind Kind
		want bool
	}{
		{Identifier, true},
		{IntegerLiteral, true},
		{FloatLiteral, true},
		{StringLiteral, true},
		{CharacterLiteral, true},
		{Plus, false},
		{KwFn, false},
		{EOF, false},
	}

	for _, tt := range tests {
		got := tt.kind.HasText()
		if got != tt.want {
			t.Errorf("HasText(%d) = %v, want %v", tt.kind, got, tt.want)
		}
	}
}
