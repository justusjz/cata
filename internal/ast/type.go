// Copyright (c) 2023 Justus Zorn

package ast

import "github.com/justusjz/cata/internal/scanner"

type TypeNode interface {
	typeNode()
	At() scanner.Pos
}

type NamedType struct {
	Name Ident
}

func (n *NamedType) typeNode() {}

func (n *NamedType) At() scanner.Pos { return n.Name.Pos }
