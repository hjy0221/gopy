package lexer

import "testing"

func TestTokensAndIndentation(t *testing.T) {
	source := "def fib(n):\n    if n < 2:\n        return n\n    return fib(n-1) + fib(n-2)\n"
	want := []TokenType{Identifier, Identifier, LParen, Identifier, RParen, Colon, Newline, Indent, Identifier, Identifier, Less, Number, Colon, Newline, Indent, Identifier, Identifier, Newline, Dedent, Identifier, Identifier, LParen, Identifier, Minus, Number, RParen, Plus, Identifier, LParen, Identifier, Minus, Number, RParen, Newline, Dedent, EOF}
	got := New(source).All()
	if len(got) != len(want) {
		t.Fatalf("got %d tokens, want %d: %#v", len(got), len(want), got)
	}
	for i, typ := range want {
		if got[i].Type != typ {
			t.Errorf("token %d = %s (%q), want %s", i, got[i].Type, got[i].Literal, typ)
		}
	}
}

func TestLiteralsAndOperators(t *testing.T) {
	got := New("x = 12.5e-2 + 'hi\\n' // 2 != 3\n").All()
	want := []TokenType{Identifier, Equal, Number, Plus, String, FloorDiv, Number, NotEqual, Number, Newline, EOF}
	if len(got) != len(want) {
		t.Fatalf("got %d tokens: %#v", len(got), got)
	}
	for i, typ := range want {
		if got[i].Type != typ {
			t.Errorf("token %d = %s, want %s", i, got[i].Type, typ)
		}
	}
}

func TestEOFEmitsFinalNewlineAndDedents(t *testing.T) {
	got := New("if x:\n  y = 1").All()
	want := []TokenType{Identifier, Identifier, Colon, Newline, Indent, Identifier, Equal, Number, Newline, Dedent, EOF}
	if len(got) != len(want) {
		t.Fatalf("got %#v", got)
	}
	for i, typ := range want {
		if got[i].Type != typ {
			t.Errorf("token %d = %s, want %s", i, got[i].Type, typ)
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	if got := New("'unfinished").NextToken(); got.Type != Illegal {
		t.Fatalf("got %v, want ILLEGAL", got)
	}
}
