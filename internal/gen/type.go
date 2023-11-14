// Copyright (c) 2023 Justus Zorn

package gen

import "github.com/justusjz/cata/internal/ast"

func typeEqual(left ast.TypeNode, right ast.TypeNode) bool {
	if nl, ok := left.(*ast.NamedType); ok {
		if nr, ok := right.(*ast.NamedType); ok {
			if nl.Name.Ident == nr.Name.Ident {
				return true
			}
		}
	}
	return false
}

func (g *generator) genType(ty ast.TypeNode) string {
	if n, ok := ty.(*ast.NamedType); ok {
		if n.Name.Ident == "i32" {
			return "int32_t"
		} else if n.Name.Ident == "bool" {
			return "_Bool"
		} else {
			g.diagnose(ty.At(), "undefined type '%s'", n.Name.Ident)
		}
	}
	g.diagnose(ty.At(), "type kind not implemented")
	return ""
}
