// Copyright (c) 2023 Justus Zorn

package scanner

type Token int

const (
	EOF Token = iota
	COMMENT
	IDENT
	INT

	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE

	PERIOD
	COMMA
	COLON
	SEMICOLON

	PLUS
	MINUS
	ASTERISK
	SLASH
)

func (t Token) String() string {
	return tokens[t]
}

var tokens = map[Token]string{
	EOF:     "end of file",
	COMMENT: "comment",
	IDENT:   "identifier",
	INT:     "integer literal",

	LPAREN:   "'('",
	RPAREN:   "')'",
	LBRACKET: "'['",
	RBRACKET: "']'",
	LBRACE:   "'{'",
	RBRACE:   "'}'",

	PERIOD:    "'.'",
	COMMA:     "','",
	COLON:     "':'",
	SEMICOLON: "';'",

	PLUS:     "'+'",
	MINUS:    "'-'",
	ASTERISK: "'*'",
	SLASH:    "'/'",
}

type Pos struct {
	line   int
	column int
}
