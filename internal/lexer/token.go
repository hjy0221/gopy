package lexer

import "fmt"

type TokenType string

const (
	Illegal      TokenType = "ILLEGAL"
	EOF          TokenType = "EOF"
	Newline      TokenType = "NEWLINE"
	Indent       TokenType = "INDENT"
	Dedent       TokenType = "DEDENT"
	Identifier   TokenType = "IDENTIFIER"
	Number       TokenType = "NUMBER"
	String       TokenType = "STRING"
	Plus         TokenType = "+"
	Minus        TokenType = "-"
	Star         TokenType = "*"
	Slash        TokenType = "/"
	FloorDiv     TokenType = "//"
	Percent      TokenType = "%"
	Equal        TokenType = "="
	EqualEqual   TokenType = "=="
	NotEqual     TokenType = "!="
	Less         TokenType = "<"
	LessEqual    TokenType = "<="
	Greater      TokenType = ">"
	GreaterEqual TokenType = ">="
	LParen       TokenType = "("
	RParen       TokenType = ")"
	LBracket     TokenType = "["
	RBracket     TokenType = "]"
	LBrace       TokenType = "{"
	RBrace       TokenType = "}"
	Colon        TokenType = ":"
	Comma        TokenType = ","
	Dot          TokenType = "."
	Semicolon    TokenType = ";"
	Arrow        TokenType = "->"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Type, t.Literal, t.Line, t.Column)
}
