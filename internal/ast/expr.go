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

func (i *IntExpr) exprNode()  {}
func (v *VarExpr) exprNode()  {}
func (c *CallExpr) exprNode() {}

func (i *IntExpr) At() scanner.Pos  { return i.Pos }
func (i *VarExpr) At() scanner.Pos  { return i.Name.Pos }
func (c *CallExpr) At() scanner.Pos { return c.Fn.At() }
