// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genFn(fnDecl *ast.FnDecl) {
	returnType := "void"
	if fnDecl.Return != nil {
		returnType = g.genType(fnDecl.Return)
	}
	params := []string{}
	for _, param := range fnDecl.Params {
		paramType := g.genType(param.Type)
		params = append(params, paramType+" "+param.Name.Ident)
	}
	strParams := strings.Join(params, ", ")
	signature := fmt.Sprintf("%s %s(%s)", returnType, fnDecl.Name.Ident, strParams)
	fmt.Fprintf(g.header, "%s;\n", signature)
	if fnDecl.Body != nil {
		fmt.Fprintf(g.body, "%s ", signature)
		g.genBlock(fnDecl.Body)
		fmt.Fprint(g.body, "\n\n")
	}
}
