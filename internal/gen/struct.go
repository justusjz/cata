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

func (g *generator) genSlice(name string, eltype *ast.NamedType) {
	// TODO: since slice only contains a pointer to the type,
	// this doesn't require generating the struct beforehand.
	// currently, however, genType always does that
	fmt.Fprintf(g.header, "struct %s {\n\t%s *data;\n\tsize_t length;\n};\n", name, g.genType(eltype))
}
