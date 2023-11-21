// Copyright (c) 2023 Justus Zorn

package parser

import (
	"github.com/justusjz/cata/internal/ast"
	"github.com/justusjz/cata/internal/scanner"
)

func (p *parser) parseIfStmt() ast.StmtNode {
	pos := p.s.Pos()
	p.s.Expect(scanner.IF, "'if'")
	p.s.Expect(scanner.LPAREN, "'('")
	cond := p.parseExpr("expression")
	p.s.Expect(scanner.RPAREN, "')'")
	body := p.parseBlock("'{'")
	if p.s.Skip(scanner.ELSE) {
		if p.s.Has(scanner.IF) {
			elif := p.parseIfStmt()
			return &ast.IfStmt{Pos: pos, Cond: cond, Body: body, Else: elif}
		} else {
			pos := p.s.Pos()
			els := p.parseBlock("'if' or '{'")
			return &ast.IfStmt{Pos: pos, Cond: cond, Body: body, Else: &ast.BlockStmt{Pos: pos, Block: els}}
		}
	} else {
		return &ast.IfStmt{Pos: pos, Cond: cond, Body: body, Else: nil}
	}
}

func (p *parser) parseWhileStmt() ast.StmtNode {
	pos := p.s.Pos()
	p.s.Expect(scanner.WHILE, "'while'")
	p.s.Expect(scanner.LPAREN, "'('")
	cond := p.parseExpr("expression")
	p.s.Expect(scanner.RPAREN, "')'")
	body := p.parseBlock("'{'")
	return &ast.WhileStmt{Pos: pos, Cond: cond, Body: body}
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
	} else if p.s.Has(scanner.IF) {
		return p.parseIfStmt()
	} else if p.s.Has(scanner.WHILE) {
		return p.parseWhileStmt()
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
