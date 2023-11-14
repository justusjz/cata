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

func (p *parser) parseExpr(expected string) ast.ExprNode {
	pos := p.s.Pos()
	if p.s.Has(scanner.INT) {
		val := p.s.Expect(scanner.INT, "")
		return &ast.IntExpr{Pos: pos, Val: val}
	} else {
		ident := p.parseIdent(expected)
		return &ast.VarExpr{Name: ident}
	}
}

func (p *parser) parseType(expected string) ast.TypeNode {
	name := p.parseIdent(expected)
	return &ast.NamedType{Name: name}
}

func (p *parser) parseStmt(expected string) ast.StmtNode {
	pos := p.s.Pos()
	if p.s.Skip(scanner.RETURN) {
		if !p.s.Skip(scanner.SEMICOLON) {
			expr := p.parseExpr("expression or ';'")
			p.s.Expect(scanner.SEMICOLON, "';'")
			return &ast.ReturnStmt{Pos: pos, Expr: expr}
		} else {
			return &ast.ReturnStmt{Pos: pos, Expr: nil}
		}
	}
	expr := p.parseExpr(expected)
	p.s.Expect(scanner.SEMICOLON, "';'")
	return &ast.ExprStmt{Expr: expr}
}

func (p *parser) parseBlock(expected string) []ast.StmtNode {
	p.s.Expect(scanner.LBRACE, expected)
	block := []ast.StmtNode{}
	for !p.s.Skip(scanner.RBRACE) {
		stmt := p.parseStmt("statement or '}'")
		block = append(block, stmt)
	}
	return block
}

func (p *parser) parseFnDecl() *ast.FnDecl {
	p.s.Expect(scanner.FN, "declaration")
	name := p.parseIdent("identifier")
	p.s.Expect(scanner.LPAREN, "'('")
	params := []ast.Param{}
	if !p.s.Has(scanner.RPAREN) {
		for {
			paramName := p.parseIdent("identifier")
			paramType := p.parseType("type")
			params = append(params, ast.Param{Name: paramName, Type: paramType})
			if !p.s.Skip(scanner.COMMA) {
				break
			}
		}
	}
	p.s.Expect(scanner.RPAREN, "')' or ','")
	var returnType ast.TypeNode = nil
	if !p.s.Has(scanner.LBRACE) {
		returnType = p.parseType("type or '{'")
	}
	body := p.parseBlock("type or '{'")
	return &ast.FnDecl{Name: name, Params: params, ReturnType: returnType, Body: body, Scanner: p.s}
}

func Parse(path string) (*ast.FnDecl, error) {
	scanner, err := scanner.New(path)
	if err != nil {
		return nil, err
	}
	p := parser{s: scanner}
	fn := p.parseFnDecl()
	return fn, nil
}
