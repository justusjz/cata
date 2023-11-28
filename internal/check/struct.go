// Copyright (c) 2023 Justus Zorn

package check

import (
	"github.com/justusjz/cata/internal/ast"
)

func (c *checker) checkStruct(decl *ast.StructDecl) {
	if decl.Done {
		// struct was already generated
		return
	}
	if decl.Started {
		// recursive struct inclusion
		c.diagnose(decl.Name.Pos, "recursive struct inclusion is not allowed")
	}
	c.newScope()
	// add generic parameters to scope
	for _, param := range decl.Params {
		if c.scope.findType(param.Ident) != nil {
			c.diagnose(param.Pos, "duplicate type '%s'", param.Ident)
		}
		c.scope.addType(param.Ident, &GenericType{})
	}
	// mark struct as started
	decl.Started = true
	// check for duplicate fields
	fields := map[string]bool{}
	if len(decl.Fields) == 0 {
		c.diagnose(decl.Name.Pos, "struct must have at least one field")
	}
	for _, field := range decl.Fields {
		if _, ok := fields[field.Name.Ident]; ok {
			c.diagnose(field.Name.Pos, "duplicate field name '%s'", field.Name.Ident)
		}
		fields[field.Name.Ident] = true
		// check the field types
		c.checkType(field.Type)
	}
	c.popScope()
	// mark struct as done
	decl.Done = true
}
