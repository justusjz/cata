// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type Param struct {
	Name Ident
	Type TypeNode
}

type FnDecl struct {
	Name       Ident
	Params     []Param
	ReturnType TypeNode
	Body       []StmtNode
	Scanner    *scanner.Scanner
}
