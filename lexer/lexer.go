// lexer/lexer.go
package lexer

import (
	"strconv"
	"strings"

	"github.com/profsergiocosta/jackcompiler-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // current char under examination
	currToken    token.Token
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

/*******************************
 nand2tetris api
/*******************************/
func (l *Lexer) advance() {
	l.currToken = l.NextToken()
}
func (l *Lexer) hasMoreTokens() bool {
	return l.currToken.Type != token.EOF
}

func (l *Lexer) tokenType() token.TokenType {
	return l.currToken.Type
}

func (l *Lexer) keyword() token.TokenType {
	return l.currToken.Type
}

func (l *Lexer) symbol() byte {
	return l.currToken.Literal[0]
}

func (l *Lexer) identifier() string {
	return l.currToken.Literal
}

func (l *Lexer) intVal() int {
	i1, err := strconv.Atoi(l.currToken.Literal)
	if err == nil {
		return i1
	}
	return -1
}

func (l *Lexer) stringVal() string {
	return l.currToken.Literal
}

/***********************************************/

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	case '/':
		if l.peekChar() == '/' {
			l.skipLineComments()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComments()
			return l.NextToken()
		} else {
			tok = newToken(token.SYMBOL, l.ch)
		}

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

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(c byte) bool {
	symbols := "{}()[].,;+-*/&|<>=~"
	return strings.IndexByte(symbols, c) != -1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.ch == 10 {
		l.readChar()
	}
}

func (l *Lexer) skipLineComments() {
	for l.ch != '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComments() {
	endComment := false
	for !endComment {
		l.readChar()
		if l.ch == '*' {
			for l.ch == '*' {
				l.readChar()
			}
			if l.ch == '/' {
				endComment = true
				l.readChar()
			}
		}
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
