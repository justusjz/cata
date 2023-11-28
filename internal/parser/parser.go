// Copyright (c) 2023 Justus Zorn

package parser

import (
	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

type parser struct {
	s *scanner.Scanner
}

func (p *parser) parseIdent(expected string) ast.Ident {
	pos := p.s.Pos()
	ident := p.s.Expect(scanner.IDENT, expected)
	return ast.Ident{Pos: pos, Ident: ident}
}

func (p *parser) parseOperand(expected string) ast.ExprNode {
	pos := p.s.Pos()
	if p.s.Has(scanner.INT) {
		val := p.s.Expect(scanner.INT, "")
		return &ast.IntExpr{Pos: pos, Val: val}
	} else {
		ident := p.parseIdent(expected)
		return &ast.VarExpr{Name: ident}
	}
}

func (p *parser) parsePrimary(expected string) ast.ExprNode {
	expr := p.parseOperand(expected)
	for p.s.Has(scanner.LPAREN) || p.s.Has(scanner.PERIOD) {
		if p.s.Skip(scanner.LPAREN) {
			// function call
			args := []ast.ExprNode{}
			if !p.s.Skip(scanner.RPAREN) {
				for {
					arg := p.parseExpr("expression")
					args = append(args, arg)
					if !p.s.Skip(scanner.COMMA) {
						break
					}
				}
				p.s.Expect(scanner.RPAREN, "',' or ')'")
			}
			expr = &ast.CallExpr{Fn: expr, Args: args}
		} else if p.s.Skip(scanner.PERIOD) {
			// field access
			field := p.parseIdent("identifier")
			expr = &ast.FieldExpr{Expr: expr, Field: field}
		}
	}
	return expr
}

func (p *parser) parseUnary(expected string) ast.ExprNode {
	return p.parsePrimary(expected)
}

func (p *parser) parseExpr(expected string) ast.ExprNode {
	return p.parseUnary(expected)
}

func Parse(path string) *ast.Module {
	s, err := scanner.New(path)
	if err != nil {
		panic("could not open file")
	}
	module := &ast.Module{Fns: []*ast.FnDecl{}}
	p := parser{s: s}
	for !p.s.Has(scanner.EOF) {
		if p.s.Has(scanner.FN) {
			fn := p.parseFnDecl()
			module.Fns = append(module.Fns, fn)
		} else {
			st := p.parseStructDecl()
			module.Structs = append(module.Structs, st)
		}
	}
	return module
}
