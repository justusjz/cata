// Copyright (c) 2023 Justus Zorn

package check

import "github.com/justusjz/cata/internal/ast"

func (c *checker) checkBlock(block []ast.StmtNode, ret *ast.NamedType) bool {
	returns := false
	for _, stmt := range block {
		if c.checkStmt(stmt, ret) {
			returns = true
		}
	}
	return returns
}

func (c *checker) checkStmt(stmt ast.StmtNode, ret *ast.NamedType) bool {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		if e, ok := s.Expr.(*ast.CallExpr); ok {
			c.checkExpr(e)
		} else {
			c.diagnose(s.Expr.At(), "expression cannot be used as statement")
		}
		return false
	case *ast.BlockStmt:
		return c.checkBlock(s.Block, ret)
	case *ast.ReturnStmt:
		if ret != nil && s.Expr == nil {
			c.diagnose(s.Pos, "expected return value of type '%s'", ret.Name.Ident)
		}
		if ret == nil && s.Expr != nil {
			c.diagnose(s.Expr.At(), "function cannot return a value")
		}
		if s.Expr != nil {
			c.checkCoerce(s.Expr, ret)
		}
		return true
	case *ast.VarStmt:
		if c.scope.findVar(s.Name.Ident) != nil {
			c.diagnose(s.Name.Pos, "duplicate identifier '%s'", s.Name.Ident)
		}
		c.scope.addVar(s.Name.Ident, &scopeVar{ty: s.Type, mut: true})
		c.checkType(s.Type)
		if s.Expr != nil {
			// check initializer expression
			c.checkCoerce(s.Expr, s.Type)
		}
		return false
	case *ast.AssignStmt:
		left := c.checkExpr(s.Left)
		if !left.mut {
			c.diagnose(s.Left.At(), "cannot assign to constant value")
		}
		c.checkCoerce(s.Right, left.ty)
		return false
	case *ast.IfStmt:
		c.checkCoerce(s.Cond, ast.Bool)
		// check if branch
		c.newScope()
		returns := c.checkBlock(s.Body, ret)
		c.popScope()
		if s.Else != nil {
			// check else branch
			c.newScope()
			if !c.checkStmt(s.Else, ret) {
				// else does not return
				returns = false
			}
			c.popScope()
		} else {
			// if stmt does not have an else, does not return
			returns = false
		}
		return returns
	case *ast.WhileStmt:
		c.checkCoerce(s.Cond, ast.Bool)
		c.newScope()
		c.checkBlock(s.Body, ret)
		c.popScope()
		return false
	}
	panic("statement kind not implemented")
}
