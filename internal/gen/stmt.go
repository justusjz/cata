// Copyright (c) 2023 Justus Zorn

package gen

import (
	"fmt"

	"github.com/justusjz/cata/internal/ast"
)

func (g *generator) genBlock(block []ast.StmtNode) {
	fmt.Fprint(g.body, "{\n")
	g.indent++
	for _, stmt := range block {
		g.genStmt(stmt, true)
	}
	g.indent--
	g.writeIndent()
	fmt.Fprint(g.body, "}")
}

func (g *generator) genStmt(stmt ast.StmtNode, indent bool) {
	if indent {
		g.writeIndent()
	}
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		expr := g.genExpr(s.Expr)
		fmt.Fprintf(g.body, "%s;\n", expr)
		return
	case *ast.BlockStmt:
		g.genBlock(s.Block)
		return
	case *ast.ReturnStmt:
		if s.Expr == nil {
			fmt.Fprint(g.body, "return;\n")
		} else {
			expr := g.genExpr(s.Expr)
			fmt.Fprintf(g.body, "return %s;\n", expr)
		}
		return
	case *ast.VarStmt:
		ty := g.genType(s.Type)
		if s.Expr != nil {
			expr := g.genExpr(s.Expr)
			fmt.Fprintf(g.body, "%s %s = %s;\n", ty, s.Name.Ident, expr)
		} else {
			fmt.Fprintf(g.body, "%s %s;\n", ty, s.Name.Ident)
		}
		return
	case *ast.AssignStmt:
		left := g.genExpr(s.Left)
		right := g.genExpr(s.Right)
		fmt.Fprintf(g.body, "%s = %s;\n", left, right)
		return
	case *ast.IfStmt:
		cond := g.genExpr(s.Cond)
		fmt.Fprintf(g.body, "if (%s) ", cond)
		g.genBlock(s.Body)
		if s.Else != nil {
			fmt.Fprint(g.body, " else ")
			g.genStmt(s.Else, false)
		}
		if indent {
			// only print newline if not in if stmt
			fmt.Fprintln(g.body)
		}
		return
	case *ast.WhileStmt:
		cond := g.genExpr(s.Cond)
		fmt.Fprintf(g.body, "while (%s) ", cond)
		g.genBlock(s.Body)
		fmt.Fprintln(g.body)
		return
	}
	fmt.Printf("%T\n", stmt)
	panic("statement kind not implemented yet")
}
