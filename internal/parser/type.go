// Copyright (c) 2023 Justus Zorn

package parser

import (
	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

func (p *parser) parseType(expected string) *ast.NamedType {
	name := p.parseIdent(expected)
	args := []*ast.NamedType{}
	if p.s.Skip(scanner.LBRACKET) {
		// generic arguments
		for {
			arg := p.parseType("type")
			args = append(args, arg)
			if !p.s.Skip(scanner.COMMA) {
				break
			}
		}
		p.s.Expect(scanner.RBRACKET, "',' or ']'")
	}
	return &ast.NamedType{Name: name, Args: args}
}
