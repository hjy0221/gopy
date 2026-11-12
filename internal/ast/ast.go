// Package ast defines the syntax tree consumed by the bytecode compiler.
package ast

type Module struct{ Body []Stmt }

type Stmt interface{ stmt() }
type Expr interface{ expr() }

type Assign struct {
	Name  string
	Value Expr
}
type ExprStmt struct{ Value Expr }
type Function struct {
	Name   string
	Params []string
	Body   []Stmt
}
type If struct {
	Test       Expr
	Body, Else []Stmt
}
type Return struct{ Value Expr }

func (Assign) stmt()   {}
func (ExprStmt) stmt() {}
func (Function) stmt() {}
func (If) stmt()       {}
func (Return) stmt()   {}

type Name struct{ Identifier string }
type Literal struct{ Kind, Source string }
type Unary struct {
	Op      string
	Operand Expr
}
type Binary struct {
	Op          string
	Left, Right Expr
}
type Compare struct {
	Left   Expr
	Ops    []string
	Rights []Expr
}
type Call struct {
	Callee Expr
	Args   []Expr
}

func (Name) expr()    {}
func (Literal) expr() {}
func (Unary) expr()   {}
func (Binary) expr()  {}
func (Compare) expr() {}
func (Call) expr()    {}
