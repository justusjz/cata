// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genFnDecl(fnDecl *ast.FnDecl) {
	// create new scope
	g.scope = newScope(g.scope)
	returnType := "void"
	if fnDecl.ReturnType != nil {
		returnType = g.genType(fnDecl.ReturnType)
	}
	params := []string{}
	for _, param := range fnDecl.Params {
		paramType := g.genType(param.Type)
		params = append(params, paramType+" "+param.Name.Ident)
		// add parameters as locals
		if g.scope.findVar(param.Name.Ident) != nil {
			g.diagnose(param.Name.Pos, "duplicate identifier '%s'", param.Name.Ident)
		}
		g.scope.addVar(param.Name.Ident, &scopeVar{ty: param.Type, mut: false})
	}
	strParams := strings.Join(params, ", ")
	signature := fmt.Sprintf("%s %s(%s)", returnType, fnDecl.Name.Ident, strParams)
	fmt.Fprintf(g.header, "%s;\n", signature)
	fmt.Fprintf(g.body, "%s ", signature)
	returns := g.genBlock(fnDecl.Body, fnDecl.ReturnType)
	if !returns && fnDecl.ReturnType != nil {
		g.diagnose(fnDecl.ReturnType.At(), "not all paths return a value")
	}
	fmt.Fprint(g.body, "\n\n")
	// reset scope
	g.scope = g.scope.parent
}

func (g *generator) genStructDecl(structDecl *ast.StructDecl) {
	if structDecl.Done {
		// struct was already generated
		return
	}
	if structDecl.Started {
		// recursive struct
		g.diagnose(structDecl.Name.Pos, "recursive struct inclusion is not allowed")
	}
	// mark struct as started
	structDecl.Started = true
	fields := map[string]bool{}
	out := "struct " + structDecl.Name.Ident + " {\n"
	if len(structDecl.Fields) == 0 {
		g.diagnose(structDecl.Name.Pos, "struct must have at least one field")
	}
	for _, field := range structDecl.Fields {
		if _, ok := fields[field.Name.Ident]; ok {
			// check for duplicate fields
			g.diagnose(field.Name.Pos, "duplicate field name '%s'", field.Name.Ident)
		}
		fields[field.Name.Ident] = true
		out += fmt.Sprintf("\t%s %s;\n", g.genType(field.Type), field.Name.Ident)
	}
	out += "};\n\n"
	fmt.Fprint(g.header, out)
	// mark struct as done
	structDecl.Done = true
}
