// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genBlock(block []ast.StmtNode, returnType ast.TypeNode) bool {
	returns := false
	fmt.Fprint(g.body, "{\n")
	g.indent++
	for _, stmt := range block {
		if g.genStmt(stmt, returnType, true) {
			returns = true
		}
	}
	g.indent--
	g.writeIndent()
	fmt.Fprint(g.body, "}")
	return returns
}

func (g *generator) genStmt(stmt ast.StmtNode, returnType ast.TypeNode, indent bool) bool {
	if indent {
		// only indent if not in if stmt
		g.writeIndent()
	}
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		if e, ok := s.Expr.(*ast.CallExpr); ok {
			expr := g.genExpr(e)
			fmt.Fprintf(g.body, "%s;\n", expr.out)
		} else {
			g.diagnose(s.Expr.At(), "expression cannot be used as statement")
		}
		return false
	case *ast.BlockStmt:
		return g.genBlock(s.Block, returnType)
	case *ast.ReturnStmt:
		if returnType != nil && s.Expr == nil {
			g.diagnose(s.Pos, "expected return value of type '%s'", returnType)
		}
		if returnType == nil && s.Expr != nil {
			g.diagnose(s.Pos, "function cannot return a value")
		}
		if returnType == nil {
			fmt.Fprint(g.body, "return;\n")
		} else {
			expr := g.genCoerce(s.Expr, returnType)
			fmt.Fprintf(g.body, "return %s;\n", expr)
		}
		return true
	case *ast.VarStmt:
		if g.scope.findVar(s.Name.Ident) != nil {
			g.diagnose(s.Name.Pos, "duplicate identifier '%s'", s.Name.Ident)
		}
		g.scope.addVar(s.Name.Ident, &scopeVar{ty: s.Type, mut: true})
		ty := g.genType(s.Type)
		if s.Expr != nil {
			expr := g.genCoerce(s.Expr, s.Type)
			fmt.Fprintf(g.body, "%s %s = %s;\n", ty, s.Name.Ident, expr)
		} else {
			fmt.Fprintf(g.body, "%s %s;\n", ty, s.Name.Ident)
		}
		return false
	case *ast.AssignStmt:
		left := g.genExpr(s.Left)
		if !left.mut {
			g.diagnose(s.Left.At(), "cannot assign to constant value")
		}
		right := g.genCoerce(s.Right, left.ty)
		fmt.Fprintf(g.body, "%s = %s;\n", left.out, right)
		return false
	case *ast.IfStmt:
		cond := g.genCoerce(s.Cond, tyBool)
		fmt.Fprintf(g.body, "if (%s) ", cond)
		g.newScope()
		returns := g.genBlock(s.Body, returnType)
		g.popScope()
		if s.Else != nil {
			fmt.Fprint(g.body, " else ")
			g.newScope()
			if !g.genStmt(s.Else, returnType, false) {
				returns = false
			}
			g.popScope()
		} else {
			returns = false
		}
		if indent {
			// only print newline if not in if stmt
			fmt.Fprintln(g.body)
		}
		return returns
	case *ast.WhileStmt:
		cond := g.genCoerce(s.Cond, tyBool)
		fmt.Fprintf(g.body, "while (%s) ", cond)
		g.newScope()
		g.genBlock(s.Body, returnType)
		g.popScope()
		fmt.Fprintln(g.body)
		return false
	}
	g.diagnose(stmt.At(), "statement kind '%T' is not implemented yet", stmt)
	return false
}
