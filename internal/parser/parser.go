// Package parser builds a syntax tree from the lexer token stream.
package parser

import (
	"fmt"
	"gopy/internal/ast"
	"gopy/internal/lexer"
)

type parser struct {
	lex       *lexer.Lexer
	current   lexer.Token
	err       error
	functions int
}

func Parse(source string) (ast.Module, error) {
	p := parser{lex: lexer.New(source)}
	p.next()
	body := p.body(lexer.EOF)
	if p.err != nil {
		return ast.Module{}, p.err
	}
	return ast.Module{Body: body}, nil
}

func (p *parser) fail(message string) {
	if p.err == nil {
		p.err = fmt.Errorf("%d:%d: %s (found %s %q)", p.current.Line, p.current.Column, message, p.current.Type, p.current.Literal)
	}
}

func (p *parser) next() lexer.Token {
	previous := p.current
	if p.err == nil {
		p.current = p.lex.NextToken()
		if p.current.Type == lexer.Illegal {
			p.fail("lexical error")
		}
	}
	return previous
}

func (p *parser) take(kind lexer.TokenType) bool {
	if p.err != nil || p.current.Type != kind {
		return false
	}
	p.next()
	return true
}

func (p *parser) expect(kind lexer.TokenType) lexer.Token {
	token := p.current
	if !p.take(kind) {
		p.fail("expected " + string(kind))
	}
	return token
}

func (p *parser) keyword(word string) bool {
	return p.current.Type == lexer.Identifier && p.current.Literal == word
}

func (p *parser) name() string {
	switch p.current.Literal {
	case "def", "if", "else", "return", "True", "False", "None":
		p.fail("expected a name")
	}
	return p.expect(lexer.Identifier).Literal
}

func (p *parser) body(end lexer.TokenType) []ast.Stmt {
	statements := []ast.Stmt{}
	for p.err == nil && p.current.Type != end {
		if p.take(lexer.Newline) {
			continue
		}
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *parser) suite() []ast.Stmt {
	p.expect(lexer.Colon)
	p.expect(lexer.Newline)
	for p.take(lexer.Newline) {
	}
	p.expect(lexer.Indent)
	body := p.body(lexer.Dedent)
	p.expect(lexer.Dedent)
	return body
}

func (p *parser) statement() ast.Stmt {
	switch {
	case p.keyword("def"):
		p.next()
		name := p.name()
		p.expect(lexer.LParen)
		params := []string{}
		seen := map[string]bool{}
		for p.err == nil && p.current.Type != lexer.RParen {
			if seen[p.current.Literal] {
				p.fail("duplicate parameter")
				break
			}
			param := p.name()
			seen[param] = true
			params = append(params, param)
			if !p.take(lexer.Comma) {
				break
			}
		}
		p.expect(lexer.RParen)
		p.functions++
		body := p.suite()
		p.functions--
		return ast.Function{Name: name, Params: params, Body: body}
	case p.keyword("if"):
		p.next()
		test := p.expression(1)
		body := p.suite()
		for p.take(lexer.Newline) {
		}
		var otherwise []ast.Stmt
		if p.keyword("else") {
			p.next()
			otherwise = p.suite()
		}
		return ast.If{Test: test, Body: body, Else: otherwise}
	case p.keyword("return"):
		if p.functions == 0 {
			p.fail("return outside function")
		}
		p.next()
		var value ast.Expr
		if p.current.Type != lexer.Newline {
			value = p.expression(1)
		}
		p.expect(lexer.Newline)
		return ast.Return{Value: value}
	default:
		value := p.expression(1)
		if p.current.Type == lexer.Equal {
			name, ok := value.(ast.Name)
			if !ok {
				p.fail("expected a name as assignment target")
			}
			p.next()
			rhs := p.expression(1)
			p.expect(lexer.Newline)
			return ast.Assign{Name: name.Identifier, Value: rhs}
		}
		p.expect(lexer.Newline)
		return ast.ExprStmt{Value: value}
	}
}

func precedence(kind lexer.TokenType) int {
	switch kind {
	case lexer.EqualEqual, lexer.NotEqual, lexer.Less, lexer.LessEqual, lexer.Greater, lexer.GreaterEqual:
		return 1
	case lexer.Plus, lexer.Minus:
		return 2
	case lexer.Star, lexer.Slash, lexer.FloorDiv, lexer.Percent:
		return 3
	default:
		return 0
	}
}

func (p *parser) expression(min int) ast.Expr {
	left := p.prefix()
	for p.err == nil {
		level := precedence(p.current.Type)
		if level < min {
			break
		}
		if level == 1 {
			comparison := ast.Compare{Left: left}
			for p.err == nil && precedence(p.current.Type) == 1 {
				comparison.Ops = append(comparison.Ops, p.next().Literal)
				comparison.Rights = append(comparison.Rights, p.expression(2))
			}
			left = comparison
			continue
		}
		op := p.next().Literal
		left = ast.Binary{Op: op, Left: left, Right: p.expression(level + 1)}
	}
	return left
}

func (p *parser) prefix() ast.Expr {
	token := p.current
	var value ast.Expr
	switch token.Type {
	case lexer.Plus, lexer.Minus:
		p.next()
		return ast.Unary{Op: token.Literal, Operand: p.prefix()}
	case lexer.Identifier:
		switch token.Literal {
		case "True", "False":
			p.next()
			value = ast.Literal{Kind: "BOOL", Source: token.Literal}
		case "None":
			p.next()
			value = ast.Literal{Kind: "NONE", Source: token.Literal}
		default:
			value = ast.Name{Identifier: p.name()}
		}
	case lexer.Number, lexer.String:
		p.next()
		value = ast.Literal{Kind: string(token.Type), Source: token.Literal}
	case lexer.LParen:
		p.next()
		value = p.expression(1)
		p.expect(lexer.RParen)
	default:
		p.fail("expected expression")
		return nil
	}
	for p.take(lexer.LParen) {
		args := []ast.Expr{}
		for p.err == nil && p.current.Type != lexer.RParen {
			args = append(args, p.expression(1))
			if !p.take(lexer.Comma) {
				break
			}
		}
		p.expect(lexer.RParen)
		value = ast.Call{Callee: value, Args: args}
	}
	return value
}
