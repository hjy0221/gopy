// Package lexer turns Python-like source text into a stream of tokens.
package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	source       string
	pos          int
	line         int
	column       int
	lineStart    bool
	indents      []int
	brackets     int
	pending      []Token
	ended        bool
	finalNewline bool
}

func New(source string) *Lexer {
	return &Lexer{source: strings.ReplaceAll(source, "\r\n", "\n"), line: 1, column: 1, lineStart: true, indents: []int{0}}
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}
func (l *Lexer) peekNext() byte {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.pos+1]
}
func (l *Lexer) advance() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	c := l.source[l.pos]
	l.pos++
	if c == '\n' {
		l.line++
		l.column = 1
		l.lineStart = true
	} else {
		l.column++
	}
	return c
}
func (l *Lexer) token(kind TokenType, lit string, line, col int) Token {
	return Token{Type: kind, Literal: lit, Line: line, Column: col}
}

// NextToken returns the next token. Lexical errors are reported as Illegal tokens.
func (l *Lexer) NextToken() Token {
	if len(l.pending) > 0 {
		t := l.pending[0]
		l.pending = l.pending[1:]
		return t
	}
	if l.ended {
		return l.token(EOF, "", l.line, l.column)
	}
	for {
		if l.pos >= len(l.source) {
			if l.pos > 0 && l.source[l.pos-1] != '\n' && l.brackets == 0 && !l.finalNewline {
				l.finalNewline = true
				return l.token(Newline, "\n", l.line, l.column)
			}
			if len(l.indents) > 1 {
				l.indents = l.indents[:len(l.indents)-1]
				return l.token(Dedent, "", l.line, l.column)
			}
			l.ended = true
			return l.token(EOF, "", l.line, l.column)
		}
		if l.lineStart && l.brackets == 0 {
			col := l.column
			width := 0
			for l.peek() == ' ' || l.peek() == '\t' {
				if l.advance() == '\t' {
					width = (width/8 + 1) * 8
				} else {
					width++
				}
			}
			if l.peek() == '\n' || l.peek() == '#' { // blank/comment-only lines do not affect indentation
				l.lineStart = false
				if l.peek() == '#' {
					for l.peek() != 0 && l.peek() != '\n' {
						l.advance()
					}
				}
				if l.peek() == '\n' {
					line, c := l.line, l.column
					l.advance()
					if l.brackets == 0 {
						return l.token(Newline, "\n", line, c)
					}
					continue
				}
			}
			l.lineStart = false
			current := l.indents[len(l.indents)-1]
			if width > current {
				l.indents = append(l.indents, width)
				return l.token(Indent, strings.Repeat(" ", width-current), l.line, col)
			}
			if width < current {
				for len(l.indents) > 1 && width < l.indents[len(l.indents)-1] {
					l.indents = l.indents[:len(l.indents)-1]
					l.pending = append(l.pending, l.token(Dedent, "", l.line, col))
				}
				if width != l.indents[len(l.indents)-1] {
					return l.token(Illegal, fmt.Sprintf("inconsistent indentation: column %d does not match any outer level", width), l.line, col)
				}
				if len(l.pending) > 0 {
					t := l.pending[0]
					l.pending = l.pending[1:]
					return t
				}
			}
		}
		c := l.peek()
		if c == ' ' || c == '\t' || c == '\f' {
			l.advance()
			continue
		}
		if c == '#' {
			for l.peek() != 0 && l.peek() != '\n' {
				l.advance()
			}
			continue
		}
		if c == '\n' {
			line, col := l.line, l.column
			l.advance()
			if l.brackets == 0 {
				return l.token(Newline, "\n", line, col)
			}
			continue
		}
		break
	}
	line, col := l.line, l.column
	l.lineStart = false
	c := l.advance()
	if isIdentStart(rune(c)) {
		start := l.pos - 1
		for isIdentContinue(rune(l.peek())) {
			l.advance()
		}
		return l.token(Identifier, l.source[start:l.pos], line, col)
	}
	if c >= '0' && c <= '9' {
		return l.number(line, col, c)
	}
	if c == '\'' || c == '"' {
		return l.readString(line, col, c)
	}
	makeToken := func(kind TokenType) Token { return l.token(kind, string(c), line, col) }
	switch c {
	case '+':
		return makeToken(Plus)
	case '-':
		if l.peek() == '>' {
			l.advance()
			return l.token(Arrow, "->", line, col)
		}
		return makeToken(Minus)
	case '*':
		return makeToken(Star)
	case '/':
		if l.peek() == '/' {
			l.advance()
			return l.token(FloorDiv, "//", line, col)
		}
		return makeToken(Slash)
	case '%':
		return makeToken(Percent)
	case '=':
		if l.peek() == '=' {
			l.advance()
			return l.token(EqualEqual, "==", line, col)
		}
		return makeToken(Equal)
	case '!':
		if l.peek() == '=' {
			l.advance()
			return l.token(NotEqual, "!=", line, col)
		}
		return l.token(Illegal, "!", line, col)
	case '<':
		if l.peek() == '=' {
			l.advance()
			return l.token(LessEqual, "<=", line, col)
		}
		return makeToken(Less)
	case '>':
		if l.peek() == '=' {
			l.advance()
			return l.token(GreaterEqual, ">=", line, col)
		}
		return makeToken(Greater)
	case '(':
		l.brackets++
		return makeToken(LParen)
	case '[':
		l.brackets++
		return makeToken(LBracket)
	case '{':
		l.brackets++
		return makeToken(LBrace)
	case ')':
		if l.brackets > 0 {
			l.brackets--
		}
		return makeToken(RParen)
	case ']':
		if l.brackets > 0 {
			l.brackets--
		}
		return makeToken(RBracket)
	case '}':
		if l.brackets > 0 {
			l.brackets--
		}
		return makeToken(RBrace)
	case ':':
		return makeToken(Colon)
	case ',':
		return makeToken(Comma)
	case '.':
		return makeToken(Dot)
	case ';':
		return makeToken(Semicolon)
	default:
		return l.token(Illegal, string(c), line, col)
	}
}

func (l *Lexer) number(line, col int, first byte) Token {
	start := l.pos - 1
	for (l.peek() >= '0' && l.peek() <= '9') || l.peek() == '_' {
		l.advance()
	}
	if l.peek() == '.' && l.peekNext() >= '0' && l.peekNext() <= '9' {
		l.advance()
		for (l.peek() >= '0' && l.peek() <= '9') || l.peek() == '_' {
			l.advance()
		}
	}
	if l.peek() == 'e' || l.peek() == 'E' {
		save := l.pos
		l.advance()
		if l.peek() == '+' || l.peek() == '-' {
			l.advance()
		}
		if l.peek() < '0' || l.peek() > '9' {
			l.pos = save
			l.column -= l.pos - save
		} else {
			for (l.peek() >= '0' && l.peek() <= '9') || l.peek() == '_' {
				l.advance()
			}
		}
	}
	return l.token(Number, l.source[start:l.pos], line, col)
}

func (l *Lexer) readString(line, col int, quote byte) Token {
	start := l.pos - 1
	triple := l.peek() == quote && l.peekNext() == quote
	if triple {
		l.advance()
		l.advance()
	}
	for {
		if l.peek() == 0 {
			return l.token(Illegal, "unterminated string", line, col)
		}
		if l.peek() == '\\' {
			l.advance()
			if l.peek() != 0 {
				l.advance()
			}
			continue
		}
		if l.peek() == quote {
			l.advance()
			if !triple || (l.peek() == quote && l.peekNext() == quote) {
				if triple {
					l.advance()
					l.advance()
				}
				l.lineStart = false
				return l.token(String, l.source[start:l.pos], line, col)
			}
			continue
		}
		if l.peek() == '\n' && !triple {
			return l.token(Illegal, "unterminated string", line, col)
		}
		l.advance()
	}
}

func isIdentStart(r rune) bool    { return r == '_' || unicode.IsLetter(r) }
func isIdentContinue(r rune) bool { return isIdentStart(r) || unicode.IsDigit(r) }

// All returns tokens through EOF.
func (l *Lexer) All() []Token {
	var result []Token
	for {
		t := l.NextToken()
		result = append(result, t)
		if t.Type == EOF {
			return result
		}
	}
}
