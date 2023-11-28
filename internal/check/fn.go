// Copyright (c) 2023 Justus Zorn

package check

import "github.com/justusjz/cata/internal/ast"

func (c *checker) checkFn(fn *ast.FnDecl) {
	c.newScope()
	if fn.Return != nil {
		// check return type
		c.checkType(fn.Return)
	}
	for _, param := range fn.Params {
		// check parameter and add to scope
		c.checkType(param.Type)
		if c.scope.findVar(param.Name.Ident) != nil {
			c.diagnose(param.Name.Pos, "duplicate identifier '%s'", param.Name.Ident)
		}
		c.scope.addVar(param.Name.Ident, &scopeVar{ty: param.Type, mut: false})
	}
	// check body
	returns := c.checkBlock(fn.Body, fn.Return)
	if !returns && fn.Return != nil {
		c.diagnose(fn.Return.At(), "not all paths return a value")
	}
	c.popScope()
}
