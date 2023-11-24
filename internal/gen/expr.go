// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/ast"
)

var tyI32 = &ast.NamedType{Name: ast.Ident{Ident: "i32"}}
var tyBool = &ast.NamedType{Name: ast.Ident{Ident: "bool"}}

type exprResult struct {
	out string
	ty  ast.TypeNode
	mut bool
}

func (g *generator) genExpr(expr ast.ExprNode, consume bool, revive bool) exprResult {
	switch e := expr.(type) {
	case *ast.IntExpr:
		return exprResult{out: e.Val, ty: tyI32, mut: false}
	case *ast.VarExpr:
		v := g.scope.findVar(e.Name.Ident)
		if v == nil {
			g.diagnose(e.At(), "undefined identifier '%s'", e.Name.Ident)
		}
		if e.Name.Ident == "true" {
			return exprResult{out: "1", ty: v.ty, mut: v.mut}
		} else if e.Name.Ident == "false" {
			return exprResult{out: "0", ty: v.ty, mut: v.mut}
		} else {
			if revive {
				g.liveness.makeAlive(e.Name.Ident)
			}
			if !g.liveness.isAlive(e.Name.Ident) {
				g.diagnose(e.At(), "'%s' is not initialized", e.Name.Ident)
			}
			if consume && g.isLinear(v.ty) {
				if !v.mut {
					g.diagnose(e.At(), "cannot consume constant value '%s'", e.Name.Ident)
				}
				if !g.liveness.isOwned(e.Name.Ident) {
					g.diagnose(e.At(), "cannot consume '%s' here", e.Name.Ident)
				}
				g.liveness.makeDead(e.Name.Ident)
			}
			return exprResult{out: e.Name.Ident, ty: v.ty, mut: v.mut}
		}
	case *ast.CallExpr:
		fn := g.genExpr(e.Fn, false, false)
		if fnTy, ok := fn.ty.(*ast.FnType); ok {
			if len(e.Args) != len(fnTy.Params) {
				g.diagnose(e.At(), "expected %d arguments, but got %d", len(fnTy.Params), len(e.Args))
			}
			args := []string{}
			for i := 0; i < len(e.Args); i++ {
				arg := g.genCoerce(e.Args[i], fnTy.Params[i], true)
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
	case *ast.StructExpr:
		ty := g.scope.findType(e.Struct.Ident)
		if ty == nil {
			g.diagnose(e.Struct.Pos, "expected struct type")
		}
		if len(e.Fields) != len(ty.decl.Fields) {
			g.diagnose(e.At(), "expected %d fields, but got %d", len(ty.decl.Fields), len(e.Fields))
		}
		fields := []string{}
		for i := 0; i < len(e.Fields); i++ {
			field := g.genCoerce(e.Fields[i], ty.decl.Fields[i].Type, true)
			fields = append(fields, field)
		}
		strFields := strings.Join(fields, ", ")
		out := fmt.Sprintf("(struct %s){%s}", e.Struct.Ident, strFields)
		return exprResult{out: out, ty: &ast.NamedType{Name: e.Struct}, mut: false}
	case *ast.FieldExpr:
		expr := g.genExpr(e.Expr, false, false)
		if namedType, ok := expr.ty.(*ast.NamedType); ok {
			ty := g.scope.findType(namedType.Name.Ident)
			if ty == nil {
				g.diagnose(e.Expr.At(), "expected expression of struct type")
			}
			decl := ty.decl
			for _, field := range decl.Fields {
				if e.Field.Ident == field.Name.Ident {
					if consume && g.isLinear(field.Type) {
						g.diagnose(e.Field.Pos, "cannot consume struct field '%s'", e.Field.Ident)
					}
					out := fmt.Sprintf("%s.%s", expr.out, field.Name.Ident)
					return exprResult{out: out, ty: field.Type, mut: expr.mut}
				}
			}
			g.diagnose(e.Field.Pos, "struct '%s' does not have field '%s'", decl.Name.Ident, e.Field.Ident)
		}
	}
	g.diagnose(expr.At(), "expression kind is not implemented yet")
	return exprResult{}
}

func (g *generator) genCoerce(expr ast.ExprNode, ty ast.TypeNode, consume bool) string {
	result := g.genExpr(expr, consume, false)
	if result.ty == nil {
		g.diagnose(expr.At(), "expression does not have a value")
	}
	if !typeEqual(result.ty, ty) {
		g.diagnose(expr.At(), "cannot convert from '%s' to '%s'", result.ty, ty)
	}
	return result.out
}
