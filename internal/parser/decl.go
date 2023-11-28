// Copyright (c) 2023 Justus Zorn

package parser

import (
	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

func (p *parser) parseParams(end scanner.Token) []ast.Param {
	params := []ast.Param{}
	if !p.s.Has(end) {
		for {
			paramName := p.parseIdent("identifier")
			p.s.Expect(scanner.COLON, "':'")
			paramType := p.parseType("type")
			params = append(params, ast.Param{Name: paramName, Type: paramType})
			if !p.s.Skip(scanner.COMMA) {
				break
			}
		}
	}
	p.s.Expect(end, "')', '}', or ','")
	return params
}

func (p *parser) parseStructDecl() *ast.StructDecl {
	p.s.Expect(scanner.STRUCT, "declaration")
	name := p.parseIdent("identifier")
	params := []ast.Ident{}
	if p.s.Skip(scanner.LBRACKET) {
		// generic parameters
		for {
			param := p.parseIdent("identifier")
			params = append(params, param)
			if !p.s.Skip(scanner.COMMA) {
				break
			}
		}
		p.s.Expect(scanner.RBRACKET, "']'")
	}
	p.s.Expect(scanner.LBRACE, "'{'")
	fields := p.parseParams(scanner.RBRACE)
	return &ast.StructDecl{Name: name, Params: params, Fields: fields, Scanner: p.s, Started: false, Done: false}
}

func (p *parser) parseFnDecl() *ast.FnDecl {
	extern := p.s.Skip(scanner.EXTERN)
	p.s.Expect(scanner.FN, "declaration")
	name := p.parseIdent("identifier")
	p.s.Expect(scanner.LPAREN, "'('")
	params := p.parseParams(scanner.RPAREN)
	var ret *ast.NamedType = nil
	if extern {
		if !p.s.Skip(scanner.SEMICOLON) {
			ret = p.parseType("type or ';'")
			p.s.Expect(scanner.SEMICOLON, "';'")
		}
		return &ast.FnDecl{Name: name, Params: params, Return: ret, Scanner: p.s}
	} else {
		if !p.s.Has(scanner.LBRACE) {
			ret = p.parseType("type or '{'")
		}
		body := p.parseBlock("type or '{'")
		return &ast.FnDecl{Name: name, Params: params, Return: ret, Body: body, Scanner: p.s}
	}
}
