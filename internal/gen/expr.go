// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

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
	case *ast.CallExpr:
		fn := g.genExpr(e.Fn)
		if fnTy, ok := fn.ty.(*ast.FnType); ok {
			if len(e.Args) != len(fnTy.Params) {
				g.diagnose(e.At(), "expected %d arguments, but got %d", len(fnTy.Params), len(e.Args))
			}
			args := []string{}
			for i := 0; i < len(e.Args); i++ {
				arg := g.genCoerce(e.Args[i], fnTy.Params[i])
				args = append(args, arg)
			}
			strArgs := strings.Join(args, ", ")
			out := fmt.Sprintf("%s(%s)", fn.out, strArgs)
			return exprResult{out: out, ty: fnTy.ReturnType, mut: false}
		} else if fn.ty == nil {
			g.diagnose(e.Fn.At(), "cannot call value of type 'void'")
		} else {
			g.diagnose(e.Fn.At(), "cannot call value of type '%s'", fn.ty)
		}
	}
	g.diagnose(expr.At(), "expression kind is not implemented yet")
	return exprResult{}
}

func (g *generator) genCoerce(expr ast.ExprNode, ty ast.TypeNode) string {
	result := g.genExpr(expr)
	if result.ty == nil {
		g.diagnose(expr.At(), "expression does not have a value")
	}
	if !typeEqual(result.ty, ty) {
		g.diagnose(expr.At(), "cannot convert from '%s' to '%s'", result.ty, ty)
	}
	return result.out
}
