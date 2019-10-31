// lexer/lexer.go
package lexer

import (
	"strings"

	"github.com/profsergiocosta/jackcompiler-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// nand2tetris api

/*

void advance();
  bool hasMoreTokens();

  TokenType tokenType();

  Keyword keyword();

  char symbol();

  string identifier();

  int intVal();

  string stringVal();

  bool isSymbol(char t);
  bool isSymbol(string t);
  bool isKeyword(string t);

  bool isStringConst(string t);

  bool isIntConst(string t);

  bool isIdentifier(string t);

*/

// privates

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isSymbol(l.ch) {
			tok = newToken(token.SYMBOL, l.ch)
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INTCONST
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(c byte) bool {
	symbols := "{}()[].,;+-*/&|<>=~"
	return strings.IndexByte(symbols, c) != -1
}
