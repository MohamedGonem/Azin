package token_test

import (
	"testing"
)

func TestKeywordsContainAllRegistered(t *testing.T) {
	expected := map[string]Kind{
		"fn":      KwFn,
		"do":      KwDo,
		"var":     KwVar,
		"mut":     KwMut,
		"return":  KwReturn,
		"end":     KwEnd,
		"char":    KwChar,
		"int":     KwInt,
		"bool":    KwBool,
		"unit":    KwUnit,
		"string":  KwString,
		"float":   KwFloat,
		"if":      KwIf,
		"then":    KwThen,
		"else":    KwElse,
		"struct":  KwStruct,
		"is":      KwIs,
		"importc": KwImportC,
		"loop":    KwLoop,
		"stop":    KwStop,
		"null":    KwNull,
		"enum":    KwEnum,
		"defer":   KwDefer,
	}

	for word, kind := range expected {
		got, ok := Keywords[word]
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

	for word := range Keywords {
		if !known[word] {
			t.Errorf("Unexpected keyword entry: %q", word)
		}
	}

	if len(Keywords) != len(known) {
		t.Errorf("Keywords map has %d entries, want %d", len(Keywords), len(known))
	}
}
