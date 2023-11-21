// Copyright (c) 2023 Justus Zorn

package gen

import "github.com/justusjz/cata/internal/ast"

type scopeVar struct {
	ty  ast.TypeNode
	mut bool
}

type scopeType struct {
	decl *ast.StructDecl
}

type scope struct {
	parent *scope
	vars   map[string]*scopeVar
	types  map[string]*scopeType
}

func newGlobalScope() *scope {
	return &scope{parent: nil, vars: map[string]*scopeVar{}, types: map[string]*scopeType{}}
}

func newScope(parent *scope) *scope {
	return &scope{parent: parent, vars: map[string]*scopeVar{}, types: map[string]*scopeType{}}
}

func (s *scope) addVar(name string, v *scopeVar) {
	if s.findVar(name) != nil || name == "true" || name == "false" {
		panic("duplicate variable name")
	}
	s.vars[name] = v
}

func (s *scope) addType(name string, t *scopeType) {
	if s.findType(name) != nil {
		panic("duplicate type name")
	}
	s.types[name] = t
}

func (s *scope) findVar(name string) *scopeVar {
	if name == "true" || name == "false" {
		return &scopeVar{ty: tyBool, mut: false}
	}
	if e, ok := s.vars[name]; ok {
		return e
	} else if s.parent != nil {
		return s.parent.findVar(name)
	} else {
		return nil
	}
}

func (s *scope) findType(name string) *scopeType {
	if t, ok := s.types[name]; ok {
		return t
	} else if s.parent != nil {
		return s.parent.findType(name)
	} else {
		return nil
	}
}
