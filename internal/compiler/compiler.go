package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"gopy/internal/ast"
	bc "gopy/internal/bytecode"
)

func Compile(module ast.Module) (*bc.Code, error) {
	code := &bc.Code{Name: "<module>"}
	c := compiler{code: code}
	c.body(module.Body)
	return code, c.err
}

type compiler struct {
	code *bc.Code
	err  error
}

func (c *compiler) emit(op bc.Op, arg any) int {
	i := len(c.code.Instructions)
	c.code.Instructions = append(c.code.Instructions, bc.Instruction{Op: op, Arg: arg})
	return i
}

func (c *compiler) patch(i int) {
	c.code.Instructions[i].Arg = len(c.code.Instructions)
}

func (c *compiler) body(body []ast.Stmt) {
	for _, stmt := range body {
		if c.err != nil {
			return
		}
		switch s := stmt.(type) {
		case ast.Assign:
			c.expr(s.Value)
			c.emit(bc.Store, s.Name)
		case ast.ExprStmt:
			c.expr(s.Value)
			c.emit(bc.Pop, nil)
		case ast.Return:
			c.expr(s.Value)
			c.emit(bc.Return, nil)
		case ast.If:
			c.expr(s.Test)
			skip := c.emit(bc.JumpFalse, 0)
			c.body(s.Body)
			end := c.emit(bc.Jump, 0)
			c.patch(skip)
			c.body(s.Else)
			c.patch(end)
		case ast.Function:
			if c.code.Name != "<module>" {
				c.err = fmt.Errorf("nested functions are not supported")
				return
			}
			fn := &bc.Code{Name: s.Name, Params: s.Params}
			sub := compiler{code: fn}
			sub.body(s.Body)
			c.err = sub.err
			c.emit(bc.Push, fn)
			c.emit(bc.Store, s.Name)
		}
	}
}

func (c *compiler) expr(expr ast.Expr) {
	if c.err != nil {
		return
	}
	switch e := expr.(type) {
	case nil:
		c.emit(bc.Push, nil)
	case ast.Name:
		c.emit(bc.Load, e.Identifier)
	case ast.Literal:
		switch e.Kind {
		case "NUMBER":
			n, err := strconv.ParseInt(strings.ReplaceAll(e.Source, "_", ""), 10, 64)
			if err != nil {
				c.err = fmt.Errorf("unsupported integer literal %q", e.Source)
				return
			}
			c.emit(bc.Push, n)
		case "BOOL":
			c.emit(bc.Push, e.Source == "True")
		case "NONE":
			c.emit(bc.Push, nil)
		default:
			c.err = fmt.Errorf("unsupported literal %q", e.Source)
		}
	case ast.Unary:
		c.expr(e.Operand)
		c.emit(bc.Unary, e.Op)
	case ast.Binary:
		c.expr(e.Left)
		c.expr(e.Right)
		c.emit(bc.Binary, e.Op)
	case ast.Compare:
		if len(e.Ops) != 1 {
			c.err = fmt.Errorf("chained comparisons are not supported")
			return
		}
		c.expr(e.Left)
		c.expr(e.Rights[0])
		c.emit(bc.Compare, e.Ops[0])
	case ast.Call:
		c.expr(e.Callee)
		for _, arg := range e.Args {
			c.expr(arg)
		}
		c.emit(bc.Call, len(e.Args))
	}
}
