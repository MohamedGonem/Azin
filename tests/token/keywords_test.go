package token_test

import (
	"testing"

	token2 "github.com/azin-lang/Azin/pkg/token"
)

func TestKeywordsContainAllRegistered(t *testing.T) {
	expected := map[string]token2.Kind{
		"fn":      token2.KwFn,
		"do":      token2.KwDo,
		"var":     token2.KwVar,
		"mut":     token2.KwMut,
		"return":  token2.KwReturn,
		"end":     token2.KwEnd,
		"char":    token2.KwChar,
		"int":     token2.KwInt,
		"bool":    token2.KwBool,
		"unit":    token2.KwUnit,
		"string":  token2.KwString,
		"float":   token2.KwFloat,
		"if":      token2.KwIf,
		"then":    token2.KwThen,
		"else":    token2.KwElse,
		"struct":  token2.KwStruct,
		"is":      token2.KwIs,
		"importc": token2.KwImportC,
		"loop":    token2.KwLoop,
		"stop":    token2.KwStop,
		"null":    token2.KwNull,
		"enum":    token2.KwEnum,
		"defer":   token2.KwDefer,
	}

	for word, kind := range expected {
		got, ok := token2.Keywords[word]
		if !ok {
			t.Errorf("Keywords map missing entry for %q", word)
			continue
		}
		if got != kind {
			t.Errorf("Keywords[%q] = %d, want %d", word, got, kind)
		}
	}
}

func TestKeywordsNoExtraEntries(t *testing.T) {
	known := map[string]bool{
		"fn": true, "do": true, "var": true, "mut": true,
		"return": true, "end": true, "char": true, "int": true,
		"bool": true, "unit": true, "string": true, "float": true,
		"if": true, "then": true, "else": true, "struct": true,
		"is": true, "importc": true, "loop": true, "stop": true,
		"null": true, "enum": true, "defer": true,
	}

	for word := range token2.Keywords {
		if !known[word] {
			t.Errorf("Unexpected keyword entry: %q", word)
		}
	}

	if len(token2.Keywords) != len(known) {
		t.Errorf("Keywords map has %d entries, want %d", len(token2.Keywords), len(known))
	}
}
