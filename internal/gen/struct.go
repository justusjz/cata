// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/check"
)

func (g *generator) genStruct(name string, decl *ast.StructDecl, args []*ast.NamedType) {
	out := "struct " + name + " {\n"
	env := check.GetGenericEnv(decl, args)
	for _, field := range decl.Fields {
		ty := check.ResolveGenericType(field.Type, env)
		out += fmt.Sprintf("\t%s %s;\n", g.genType(ty), field.Name.Ident)
	}
	out += "};\n\n"
	fmt.Fprint(g.header, out)
}
