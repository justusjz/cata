// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genFnDecl(fnDecl *ast.FnDecl) {
	returnType := "void"
	if fnDecl.ReturnType != nil {
		returnType = g.genType(fnDecl.ReturnType)
	}
	params := []string{}
	for _, param := range fnDecl.Params {
		paramType := g.genType(param.Type)
		params = append(params, paramType+" "+param.Name.Ident)
	}
	strParams := strings.Join(params, ", ")
	signature := fmt.Sprintf("%s %s(%s)", returnType, fnDecl.Name.Ident, strParams)
	fmt.Fprintf(g.header, "%s;\n", signature)
	fmt.Fprintf(g.body, "%s ", signature)
	returns := g.genBlock(fnDecl.Body, fnDecl.ReturnType)
	if !returns && fnDecl.ReturnType != nil {
		g.diagnose(fnDecl.ReturnType.At(), "not all paths return a value")
	}
}
