// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type StmtNode interface {
	stmtNode()
}

type ExprStmt struct {
	Expr ExprNode
}

type ReturnStmt struct {
	Pos  scanner.Pos
	Expr ExprNode
}

func (e *ExprStmt) stmtNode()   {}
func (r *ReturnStmt) stmtNode() {}
