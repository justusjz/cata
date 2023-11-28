// Copyright (c) 2023 Justus Zorn

package scanner

type Token int

const (
	EOF Token = iota
	COMMENT
	IDENT
	INT
	STRING

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
	CARET

	ASSIGN

	ELSE
	EXTERN
	FN
	IF
	RETURN
	STRUCT
	VAR
	WHILE
)

func (t Token) String() string {
	return tokens[t]
}

var tokens = map[Token]string{
	EOF:     "end of file",
	COMMENT: "comment",
	IDENT:   "identifier",
	INT:     "integer literal",
	STRING:  "string literal",

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
	CARET:    "'^'",

	ASSIGN: "'='",

	ELSE:   "'else'",
	EXTERN: "'extern'",
	FN:     "'fn'",
	IF:     "'if'",
	RETURN: "'return'",
	STRUCT: "'struct'",
	VAR:    "'var'",
	WHILE:  "'while'",
}

type Pos struct {
	line   int
	column int
}
