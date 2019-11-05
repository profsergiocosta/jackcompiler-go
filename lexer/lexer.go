// lexer/lexer.go
package lexer

import (
	"fmt"
	"os"

	"github.com/profsergiocosta/jackcompiler-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // current char under examination
	currToken    token.Token
	curLine      int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	l.curLine = 1
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = l.curLine
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComments()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComments()
			return l.NextToken()
		} else {
			tok = l.newToken(token.SLASH, l.ch)
		}

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.curLine
	default:
		if token.IsSymbol(l.ch) {
			tok = l.newToken(token.LookupSymbol(l.ch), l.ch)
		} else if token.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.curLine
			return tok
		} else if token.IsDigit(l.ch) {
			tok.Type = token.INTCONST
			tok.Literal = l.readNumber()
			tok.Line = l.curLine
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.curLine}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.ch == 10 {

		if l.ch == '\n' {
			l.curLine++
		}
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

		if l.ch == 0 { // eof
			fmt.Printf("%v: expected %s, got %s instead", l.curLine, "*/", "EOF")
			os.Exit(1)
		}

		if l.ch == '\n' {
			l.curLine++
		}
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
	for token.IsLetter(l.ch) || token.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for token.IsDigit(l.ch) {
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
