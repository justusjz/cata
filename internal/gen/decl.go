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
		if g.scope.find(param.Name.Ident) != nil {
			g.diagnose(param.Name.Pos, "duplicate identifier '%s'", param.Name.Ident)
		}
		g.scope.add(param.Name.Ident, &scopeVar{ty: param.Type, mut: false})
	}
	strParams := strings.Join(params, ", ")
	signature := fmt.Sprintf("%s %s(%s)", returnType, fnDecl.Name.Ident, strParams)
	fmt.Fprintf(g.header, "%s;\n", signature)
	fmt.Fprintf(g.body, "%s ", signature)
	returns := g.genBlock(fnDecl.Body, fnDecl.ReturnType)
	if !returns && fnDecl.ReturnType != nil {
		g.diagnose(fnDecl.ReturnType.At(), "not all paths return a value")
	}
	// reset scope
	g.scope = g.scope.parent
}
