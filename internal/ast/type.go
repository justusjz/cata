// Copyright (c) 2023 Justus Zorn

package ast

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/scanner"
)

type TypeNode interface {
	typeNode()
	At() scanner.Pos
	String() string
}

type NamedType struct {
	Name Ident
}

type FnType struct {
	Pos        scanner.Pos
	Params     []TypeNode
	ReturnType TypeNode
}

func (n *NamedType) typeNode() {}
func (f *FnType) typeNode()    {}

func (n *NamedType) At() scanner.Pos { return n.Name.Pos }
func (f *FnType) At() scanner.Pos    { return f.Pos }

func (n *NamedType) String() string { return n.Name.Ident }

func (f *FnType) String() string {
	params := []string{}
	for _, param := range f.Params {
		params = append(params, param.String())
	}
	strParams := strings.Join(params, ", ")
	if f.ReturnType != nil {
		return fmt.Sprintf("fn (%s) %s", strParams, f.ReturnType)
	} else {
		return fmt.Sprintf("fn (%s)", strParams)
	}
}
