// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genBlock(block []ast.StmtNode, returnType ast.TypeNode) bool {
	returns := false
	fmt.Fprint(g.body, "{\n")
	g.indent += 4
	for _, stmt := range block {
		if g.genStmt(stmt, returnType) {
			returns = true
		}
	}
	g.indent -= 4
	g.writeIndent()
	fmt.Fprint(g.body, "}\n")
	return returns
}

func (g *generator) genStmt(stmt ast.StmtNode, returnType ast.TypeNode) bool {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		g.diagnose(s.Expr.At(), "expression cannot be used as statements")
	case *ast.ReturnStmt:
		if returnType != nil && s.Expr == nil {
			g.diagnose(s.Pos, "expected return value of type '%s'", returnType)
		}
		if returnType == nil && s.Expr != nil {
			g.diagnose(s.Pos, "function cannot return a value")
		}
		g.writeIndent()
		if returnType == nil {
			fmt.Fprint(g.body, "return;\n")
		} else {
			expr := g.genCoerce(s.Expr, returnType)
			fmt.Fprintf(g.body, "return %s;\n", expr)
		}
		return true
	}
	return false
}
