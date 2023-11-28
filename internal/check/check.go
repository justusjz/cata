// Copyright (c) 2023 Justus Zorn

package check

import (
	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

type checker struct {
	scanner *scanner.Scanner
	scope   *scope
	structs map[string]*ast.StructDecl
}

func (c *checker) diagnose(pos scanner.Pos, format string, a ...any) {
	c.scanner.Diagnose(pos, format, a...)
}

func (c *checker) newScope() {
	c.scope = newScope(c.scope)
}

func (c *checker) popScope() {
	c.scope = c.scope.parent
}

type GenericType struct {
	Params []string
}

func CheckModule(module *ast.Module) {
	c := checker{scope: newGlobalScope(), structs: map[string]*ast.StructDecl{}}
	// add primitives to scope
	c.scope.addType("u8", &GenericType{})
	c.scope.addType("i32", &GenericType{})
	c.scope.addType("bool", &GenericType{})
	c.scope.addType("slice", &GenericType{Params: []string{"eltype"}})
	// add structs to scope
	for _, st := range module.Structs {
		c.scanner = st.Scanner
		if c.scope.findType(st.Name.Ident) != nil {
			c.diagnose(st.Name.Pos, "duplicate type '%s'", st.Name.Ident)
		}
		params := []string{}
		for _, param := range st.Params {
			params = append(params, param.Ident)
		}
		c.scope.addType(st.Name.Ident, &GenericType{Params: params})
		c.structs[st.Name.Ident] = st
	}
	// check struct definitions
	for _, st := range module.Structs {
		c.scanner = st.Scanner
		c.checkStruct(st)
	}
	// add functions to scope
	for _, fn := range module.Fns {
		c.scanner = fn.Scanner
		if c.scope.findVar(fn.Name.Ident) != nil {
			c.diagnose(fn.Name.Pos, "duplicate identifier '%s'", fn.Name.Ident)
		}
		params := []*ast.NamedType{}
		for _, param := range fn.Params {
			params = append(params, param.Type)
		}
		c.scope.addVar(fn.Name.Ident, &scopeVar{ty: &ast.FnType{Params: params, Return: fn.Return}, mut: false})
	}
	// check function definitions
	for _, fn := range module.Fns {
		c.scanner = fn.Scanner
		c.checkFn(fn)
	}
}
