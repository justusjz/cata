// Copyright (c) 2023 Justus Zorn

package ast

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/scanner"
)

type Type interface {
	ty()
	String() string
}

type NamedType struct {
	Name Ident
	Args []*NamedType
}

type FnType struct {
	Params []*NamedType
	Return *NamedType
}

var Uint8 = &NamedType{Name: Ident{Ident: "u8"}}
var Int32 = &NamedType{Name: Ident{Ident: "i32"}}
var Bool = &NamedType{Name: Ident{Ident: "bool"}}

func SliceType(eltype *NamedType) *NamedType {
	return &NamedType{Name: Ident{Ident: "slice"}, Args: []*NamedType{eltype}}
}

func (n *NamedType) At() scanner.Pos { return n.Name.Pos }
func (n *NamedType) ty()             {}

func (n *NamedType) String() string {
	if len(n.Args) == 0 {
		return n.Name.Ident
	}
	args := []string{}
	for _, arg := range n.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s[%s]", n.Name.Ident, strings.Join(args, ", "))
}

func (f *FnType) ty() {}

func (f *FnType) String() string {
	params := []string{}
	for _, param := range f.Params {
		params = append(params, param.String())
	}
	result := "fn (" + strings.Join(params, ", ") + ")"
	if f.Return != nil {
		return result + " " + f.Return.String()
	} else {
		return result
	}
}
