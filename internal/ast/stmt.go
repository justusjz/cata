// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type StmtNode interface {
	stmtNode()
	At() scanner.Pos
}

type ExprStmt struct {
	Expr ExprNode
}

type ReturnStmt struct {
	Pos  scanner.Pos
	Expr ExprNode
}

type VarStmt struct {
	Pos  scanner.Pos
	Name Ident
	Type TypeNode
	Expr ExprNode
}

func (e *ExprStmt) stmtNode()   {}
func (r *ReturnStmt) stmtNode() {}
func (v *VarStmt) stmtNode()    {}

func (e *ExprStmt) At() scanner.Pos   { return e.Expr.At() }
func (r *ReturnStmt) At() scanner.Pos { return r.Pos }
func (v *VarStmt) At() scanner.Pos    { return v.Pos }
