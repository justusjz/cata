// Copyright (c) 2023 Justus Zorn

package gen

import (
	"github.com/justusjz/cata/internal/ast"
)

var i32 = &ast.NamedType{Name: ast.Ident{Ident: "i32"}}

func (g *generator) genExpr(expr ast.ExprNode) (string, ast.TypeNode) {
	switch e := expr.(type) {
	case *ast.IntExpr:
		return e.Val, i32
	}
	g.diagnose(expr.At(), "expression kind is not implemented yet")
	return "", nil
}

func (g *generator) genCoerce(expr ast.ExprNode, ty ast.TypeNode) string {
	e, realTy := g.genExpr(expr)
	if !typeEqual(realTy, ty) {
		g.diagnose(expr.At(), "cannot convert from '%s' to '%s'", realTy, ty)
	}
	return e
}
