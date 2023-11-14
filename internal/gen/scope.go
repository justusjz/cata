// Copyright (c) 2023 Justus Zorn

package gen

import "github.com/justusjz/cata/internal/ast"

type scopeEntry interface {
	scopeEntry()
}

type scopeVar struct {
	ty  ast.TypeNode
	mut bool
}

func (s *scopeVar) scopeEntry() {}

type scope struct {
	parent  *scope
	entries map[string]scopeEntry
}

func newGlobalScope() *scope {
	return &scope{parent: nil, entries: map[string]scopeEntry{}}
}

func newScope(parent *scope) *scope {
	return &scope{parent: parent, entries: map[string]scopeEntry{}}
}

func (s *scope) add(name string, entry scopeEntry) {
	if s.find(name) != nil {
		panic("duplicate variable")
	}
	s.entries[name] = entry
}

func (s *scope) find(name string) scopeEntry {
	if e, ok := s.entries[name]; ok {
		return e
	} else if s.parent != nil {
		return s.parent.find(name)
	} else {
		return nil
	}
}
