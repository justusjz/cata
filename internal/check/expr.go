// Copyright (c) 2023 Justus Zorn

package check

import (
	"github.com/justusjz/cata/internal/ast"
)

type result struct {
	ty  ast.Type
	mut bool
}

func (c *checker) checkExpr(expr ast.ExprNode) result {
	switch e := expr.(type) {
	case *ast.IntExpr:
		return result{ty: ast.Int32, mut: false}
	case *ast.StringExpr:
		return result{ty: ast.SliceType(ast.Uint8), mut: false}
	case *ast.VarExpr:
		v := c.scope.findVar(e.Name.Ident)
		if v == nil {
			// variable not found
			c.diagnose(e.At(), "undefined identifier '%s'", e.Name.Ident)
		}
		return result{ty: v.ty, mut: v.mut}
	case *ast.CallExpr:
		fn := c.checkExpr(e.Fn)
		if ty, ok := fn.ty.(*ast.FnType); ok {
			if len(e.Args) != len(ty.Params) {
				// incorrect argument count
				c.diagnose(e.At(), "expected %d arguments, but got %d", len(ty.Params), len(e.Args))
			}
			for i := 0; i < len(e.Args); i++ {
				// coerce arguments to parameter types
				c.checkCoerce(e.Args[i], ty.Params[i])
			}
			return result{ty: ty.Return, mut: false}
		} else if fn.ty == nil {
			c.diagnose(e.Fn.At(), "cannot call expression that does not have a value")
		} else {
			c.diagnose(e.Fn.At(), "cannot call value of type '%s'", fn.ty)
		}
	case *ast.FieldExpr:
		expr := c.checkExpr(e.Expr)
		field := c.getField(expr.ty, e.Field.Ident)
		if field == nil {
			c.diagnose(e.Expr.At(), "value does not have field '%s'", e.Field.Ident)
		}
		return result{ty: field, mut: expr.mut}
	}
	panic("expression kind not implemented")
}

func (c *checker) checkCoerce(expr ast.ExprNode, ty ast.Type) {
	result := c.checkExpr(expr)
	if result.ty == nil {
		c.diagnose(expr.At(), "expression does not have a value")
	}
	if !TypeEqual(result.ty, ty) {
		c.diagnose(expr.At(), "cannot convert from '%s' to '%s'", result.ty, ty)
	}
}
