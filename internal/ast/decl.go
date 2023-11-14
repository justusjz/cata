// Copyright (c) 2023 Justus Zorn

package ast

type Param struct {
	Name Ident
	Type TypeNode
}

type FnDecl struct {
	Name       Ident
	Params     []Param
	ReturnType TypeNode
	Body       []StmtNode
}
