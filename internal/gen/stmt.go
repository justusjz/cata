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
		if e, ok := s.Expr.(*ast.CallExpr); ok {
			g.writeIndent()
			expr := g.genExpr(e)
			fmt.Fprintf(g.body, "%s;\n", expr.out)
		} else {
			g.diagnose(s.Expr.At(), "expression cannot be used as statement")
		}
		return false
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
	case *ast.VarStmt:
		if g.scope.findVar(s.Name.Ident) != nil {
			g.diagnose(s.Name.Pos, "duplicate identifier '%s'", s.Name.Ident)
		}
		g.scope.addVar(s.Name.Ident, &scopeVar{ty: s.Type, mut: true})
		ty := g.genType(s.Type)
		expr := g.genCoerce(s.Expr, s.Type)
		g.writeIndent()
		fmt.Fprintf(g.body, "%s %s = %s;\n", ty, s.Name.Ident, expr)
		return false
	case *ast.AssignStmt:
		left := g.genExpr(s.Left)
		if !left.mut {
			g.diagnose(s.Left.At(), "cannot assign to constant value")
		}
		right := g.genCoerce(s.Right, left.ty)
		g.writeIndent()
		fmt.Fprintf(g.body, "%s = %s;\n", left.out, right)
		return false
	}
	g.diagnose(stmt.At(), "statement kind '%T' is not implemented yet", stmt)
	return false
}
