// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
)

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
			decl := g.scope.findType(n.Name.Ident)
			if decl == nil {
				g.diagnose(ty.At(), "undefined type '%s'", n.Name.Ident)
			}
			// generate the struct
			g.genStructDecl(decl.decl)
			return fmt.Sprintf("struct %s", n.Name.Ident)
		}
	}
	g.diagnose(ty.At(), "type kind not implemented")
	return ""
}

func (g *generator) isLinear(ty ast.TypeNode) bool {
	if _, ok := ty.(*ast.FnType); ok {
		return false
	}
	if n, ok := ty.(*ast.NamedType); ok {
		if n.Name.Ident == "i32" || n.Name.Ident == "bool" {
			return false
		}
		decl := g.scope.findType(n.Name.Ident)
		if decl.decl.Linear {
			// the struct itself is linear
			return true
		}
		for _, field := range decl.decl.Fields {
			if g.isLinear(field.Type) {
				// one of the fields is linear
				return true
			}
		}
		return false
	}
	panic("type kind not implemented")
}
