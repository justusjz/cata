// Copyright (c) 2023 Justus Zorn

package scanner

import (
	"fmt"
	"os"
)

type Scanner struct {
	path    string
	content string
	pos     int
	line    int
	lbegin  int
	tok_pos Pos
	tok     Token
	tok_val string
}

func New(path string) (*Scanner, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := &Scanner{path: path, content: string(content), pos: 0, line: 1, lbegin: 0}
	// scan first token
	s.next()
	return s, nil
}

func (s *Scanner) Diagnose(pos Pos, format string, a ...any) {
	fmt.Printf("%s:%d:%d: error: ", s.path, pos.line, pos.column)
	fmt.Printf(format, a...)
	fmt.Println()
	os.Exit(1)
}

func (s *Scanner) Pos() Pos {
	return s.tok_pos
}

func (s *Scanner) Has(tok Token) bool {
	return s.tok == tok
}

func (s *Scanner) Skip(tok Token) bool {
	if s.tok == tok {
		s.next()
		return true
	} else {
		return false
	}
}

func (s *Scanner) Expect(tok Token, expected string) string {
	if s.tok == tok {
		// save current token, scan next one
		val := s.tok_val
		s.next()
		return val
	} else {
		// expected different token
		if s.tok == IDENT {
			s.Diagnose(s.tok_pos, "expected %s but got '%s'", expected, s.tok_val)
		} else {
			s.Diagnose(s.tok_pos, "expected %s but got %s", expected, s.tok)
		}
		return ""
	}
}

func (s *Scanner) next() {
	s.scan()
	// skip comments for now
	for s.tok == COMMENT {
		s.scan()
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func isEscape(c byte) bool {
	return c == '"' || c == '\'' || c == 'n' || c == 't' || c == '\\'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) scan() {
	// skip whitespace
	for s.pos < len(s.content) && isWhitespace(s.content[s.pos]) {
		if s.content[s.pos] == '\n' {
			s.line++
			s.lbegin = s.pos + 1
		}
		s.pos++
	}
	// compute token position
	s.tok_pos = Pos{line: s.line, column: s.pos - s.lbegin + 1}
	if s.pos == len(s.content) {
		// end of file
		s.tok = EOF
	} else if isLetter(s.content[s.pos]) {
		// identifier
		begin := s.pos
		for s.pos < len(s.content) && (isLetter(s.content[s.pos]) || isDigit(s.content[s.pos])) {
			s.pos++
		}
		val := s.content[begin:s.pos]
		if keyword, ok := keywords[val]; ok {
			// keyword
			s.tok = keyword
		} else {
			// identifier
			s.tok_val = s.content[begin:s.pos]
			s.tok = IDENT
		}
	} else if isDigit(s.content[s.pos]) {
		// integer literal
		begin := s.pos
		for s.pos < len(s.content) && isDigit(s.content[s.pos]) {
			s.pos++
		}
		s.tok_val = s.content[begin:s.pos]
		s.tok = INT
	} else if s.content[s.pos] == '"' {
		// string literal
		s.pos++
		begin := s.pos
		escaped := false
		for s.pos < len(s.content) && (s.content[s.pos] != '"' || escaped) {
			if escaped && !isEscape(s.content[s.pos]) {
				s.Diagnose(s.tok_pos, "'\\%c' is not a valid escape sequence", s.content[s.pos])
			}
			escaped = !escaped && s.content[s.pos] == '\\'
			if s.content[s.pos] == '\n' {
				// string literal must not contain newline
				s.Diagnose(s.tok_pos, "unterminated string literal")
			}
			s.pos++
		}
		// no ending " found
		if s.pos == len(s.content) {
			s.Diagnose(s.tok_pos, "unterminated string literal")
		}
		s.tok_val = s.content[begin:s.pos]
		s.tok = STRING
		s.pos++
	} else {
		c := s.content[s.pos]
		s.pos++
		switch c {
		case '(':
			s.tok = LPAREN
		case ')':
			s.tok = RPAREN
		case '[':
			s.tok = LBRACKET
		case ']':
			s.tok = RBRACKET
		case '{':
			s.tok = LBRACE
		case '}':
			s.tok = RBRACE
		case '.':
			s.tok = PERIOD
		case ',':
			s.tok = COMMA
		case ':':
			s.tok = COLON
		case ';':
			s.tok = SEMICOLON
		case '+':
			s.tok = PLUS
		case '-':
			s.tok = MINUS
		case '*':
			s.tok = ASTERISK
		case '/':
			if s.pos < len(s.content) && s.content[s.pos] == '/' {
				// comment
				for s.pos < len(s.content) && s.content[s.pos] != '\n' {
					s.pos++
				}
				s.tok = COMMENT
			} else {
				s.tok = SLASH
			}
		case '^':
			s.tok = CARET
		case '=':
			s.tok = ASSIGN
		default:
			s.Diagnose(s.tok_pos, "invalid character '%c'", s.content[s.pos-1])
		}
	}
}

var keywords = map[string]Token{
	"else":   ELSE,
	"extern": EXTERN,
	"fn":     FN,
	"if":     IF,
	"return": RETURN,
	"struct": STRUCT,
	"var":    VAR,
	"while":  WHILE,
}
