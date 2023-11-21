// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type ExprNode interface {
	exprNode()
	At() scanner.Pos
}

type IntExpr struct {
	Pos scanner.Pos
	Val string
}

type VarExpr struct {
	Name Ident
}

type CallExpr struct {
	Fn   ExprNode
	Args []ExprNode
}

type StructExpr struct {
	Struct Ident
	Fields []ExprNode
}

type FieldExpr struct {
	Expr  ExprNode
	Field Ident
}

func (i *IntExpr) exprNode()    {}
func (v *VarExpr) exprNode()    {}
func (c *CallExpr) exprNode()   {}
func (s *StructExpr) exprNode() {}
func (f *FieldExpr) exprNode()  {}

func (i *IntExpr) At() scanner.Pos    { return i.Pos }
func (i *VarExpr) At() scanner.Pos    { return i.Name.Pos }
func (c *CallExpr) At() scanner.Pos   { return c.Fn.At() }
func (s *StructExpr) At() scanner.Pos { return s.Struct.Pos }
func (f *FieldExpr) At() scanner.Pos  { return f.Expr.At() }
