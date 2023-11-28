// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"
	"strings"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genExpr(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.IntExpr:
		return e.Val
	case *ast.VarExpr:
		if e.Name.Ident == "true" {
			return "1"
		} else if e.Name.Ident == "false" {
			return "0"
		} else {
			return e.Name.Ident
		}
	case *ast.CallExpr:
		fn := g.genExpr(e.Fn)
		args := []string{}
		for i := 0; i < len(e.Args); i++ {
			arg := g.genExpr(e.Args[i])
			args = append(args, arg)
		}
		strArgs := strings.Join(args, ", ")
		out := fmt.Sprintf("%s(%s)", fn, strArgs)
		return out
	case *ast.FieldExpr:
		expr := g.genExpr(e.Expr)
		return fmt.Sprintf("%s.%s", expr, e.Field.Ident)
	}
	panic("expression kind not implemented")
}
