// Copyright (c) 2023 Justus Zorn

package ast

import (
	"github.com/justusjz/cata/internal/scanner"
)

type Ident struct {
	Pos   scanner.Pos
	Ident string
}

type Module struct {
	Fns     []*FnDecl
	Structs []*StructDecl
}
