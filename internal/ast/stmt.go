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

type BlockStmt struct {
	Pos   scanner.Pos
	Block []StmtNode
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

type AssignStmt struct {
	Left  ExprNode
	Right ExprNode
}

type IfStmt struct {
	Pos    scanner.Pos
	Cond   ExprNode
	Body   []StmtNode
	ElseIf StmtNode
	Else   StmtNode
}

type WhileStmt struct {
	Pos  scanner.Pos
	Cond ExprNode
	Body []StmtNode
}

func (e *ExprStmt) stmtNode()   {}
func (b *BlockStmt) stmtNode()  {}
func (r *ReturnStmt) stmtNode() {}
func (v *VarStmt) stmtNode()    {}
func (a *AssignStmt) stmtNode() {}
func (i *IfStmt) stmtNode()     {}
func (w *WhileStmt) stmtNode()  {}

func (e *ExprStmt) At() scanner.Pos   { return e.Expr.At() }
func (b *BlockStmt) At() scanner.Pos  { return b.Pos }
func (r *ReturnStmt) At() scanner.Pos { return r.Pos }
func (v *VarStmt) At() scanner.Pos    { return v.Pos }
func (a *AssignStmt) At() scanner.Pos { return a.Left.At() }
func (i *IfStmt) At() scanner.Pos     { return i.Pos }
func (w *WhileStmt) At() scanner.Pos  { return w.Pos }
