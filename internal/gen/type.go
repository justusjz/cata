// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/check"
)

func (g *generator) genType(ty *ast.NamedType) string {
	if ty.Name.Ident == "u8" {
		return "uint8_t"
	} else if ty.Name.Ident == "i32" {
		return "int32_t"
	} else if ty.Name.Ident == "bool" {
		return "_Bool"
	} else {
		for _, instance := range g.instances {
			if check.TypeEqual(instance.ty, ty) {
				// struct was already generated
				return "struct " + instance.name
			}
		}
		// generate the struct
		name := "ty" + fmt.Sprint(g.counter)
		g.counter++
		g.instances = append(g.instances, genericInstance{ty: ty, name: name})
		if ty.Name.Ident == "slice" {
			g.genSlice(name, ty.Args[0])
		} else {
			decl := g.structs[ty.Name.Ident]
			g.genStruct(name, decl, ty.Args)
		}
		return fmt.Sprintf("struct %s", name)
	}
}
