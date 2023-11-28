// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type Param struct {
	Name Ident
	Type *NamedType
}

type FnDecl struct {
	Name    Ident
	Params  []Param
	Return  *NamedType
	Body    []StmtNode
	Scanner *scanner.Scanner
}

type StructDecl struct {
	Name    Ident
	Params  []Ident
	Fields  []Param
	Scanner *scanner.Scanner
	Started bool
	Done    bool
}
