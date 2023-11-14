// Copyright (c) 2023 Justus Zorn

package gen

import (
	"github.com/justusjz/cata/internal/ast"
)

var i32 = &ast.NamedType{Name: ast.Ident{Ident: "i32"}}

type exprResult struct {
	out string
	ty  ast.TypeNode
	mut bool
}

func (g *generator) genExpr(expr ast.ExprNode) exprResult {
	switch e := expr.(type) {
	case *ast.IntExpr:
		return exprResult{out: e.Val, ty: i32, mut: false}
	case *ast.VarExpr:
		entry := g.scope.find(e.Name.Ident)
		if entry == nil {
			g.diagnose(e.At(), "undefined identifier '%s'", e.Name.Ident)
		}
		if v, ok := entry.(*scopeVar); ok {
			return exprResult{out: e.Name.Ident, ty: v.ty, mut: v.mut}
		} else {
			g.diagnose(e.At(), "expected variable name")
		}
	}
	g.diagnose(expr.At(), "expression kind is not implemented yet")
	return exprResult{}
}

func (g *generator) genCoerce(expr ast.ExprNode, ty ast.TypeNode) string {
	result := g.genExpr(expr)
	if !typeEqual(result.ty, ty) {
		g.diagnose(expr.At(), "cannot convert from '%s' to '%s'", result.ty, ty)
	}
	return result.out
}
