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

func (p *parser) parseStructExpr(ident ast.Ident) ast.ExprNode {
	fields := []ast.ExprNode{}
	if !p.s.Skip(scanner.RBRACE) {
		for {
			field := p.parseExpr("expression")
			fields = append(fields, field)
			if !p.s.Skip(scanner.COMMA) {
				break
			}
		}
		p.s.Expect(scanner.RBRACE, "',' or '}'")
	}
	return &ast.StructExpr{Struct: ident, Fields: fields}
}

func (p *parser) parseOperand(expected string) ast.ExprNode {
	pos := p.s.Pos()
	if p.s.Has(scanner.INT) {
		val := p.s.Expect(scanner.INT, "")
		return &ast.IntExpr{Pos: pos, Val: val}
	} else {
		ident := p.parseIdent(expected)
		if p.s.Skip(scanner.LBRACE) {
			return p.parseStructExpr(ident)
		}
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

func (p *parser) parseExpr(expected string) ast.ExprNode {
	return p.parsePrimary(expected)
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
	} else if p.s.Skip(scanner.VAR) {
		name := p.parseIdent("identifier")
		p.s.Expect(scanner.COLON, "':'")
		ty := p.parseType("type")
		p.s.Expect(scanner.ASSIGN, "'='")
		expr := p.parseExpr("expression")
		p.s.Expect(scanner.SEMICOLON, "';'")
		return &ast.VarStmt{Pos: pos, Name: name, Type: ty, Expr: expr}
	}
	expr := p.parseExpr(expected)
	if p.s.Skip(scanner.ASSIGN) {
		right := p.parseExpr("expression")
		p.s.Expect(scanner.SEMICOLON, "';'")
		return &ast.AssignStmt{Left: expr, Right: right}
	} else {
		p.s.Expect(scanner.SEMICOLON, "';' or '='")
		return &ast.ExprStmt{Expr: expr}
	}
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
