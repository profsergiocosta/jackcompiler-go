package parser

import (
	"fmt"
	"os"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
	"github.com/profsergiocosta/jackcompiler-go/token"
	"github.com/profsergiocosta/jackcompiler-go/xmlwrite"
)

const (
	XML = "XML"
	VM  = "VM"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	output    string
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	p.output = XML
	return p
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Compile() {
	p.CompileClass()

}

func (p *Parser) CompileClass() {

	p.expectTokenByLiteral("class")
	xmlwrite.TagNonTerminal("class", p.output == XML)
	p.nextToken()

	p.expectTokenByType(token.IDENT)
	xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	p.nextToken()

	p.expectTokenByLiteral("{")
	xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	p.nextToken()

	xmlwrite.UntagNonTerminal("class", p.output == XML)

}

// book write an interpreter in go
/*
func (p *Parser) curTokenIs(t string) bool {
	return p.curToken.Literal == t || string(p.curToken.Type) == t // I am not sure about this
}

func (p *Parser) peekTokenIs(t string) bool {
	return p.peekToken.Literal == t || string(p.peekToken.Type) == t // I am not sure about this
}

func (p *Parser) expectPeek(token string) bool {
	if p.peekTokenIs(token) {
		p.nextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

*/

func (p *Parser) expectTokenByType(token token.TokenType) {
	if p.curToken.Type == token {
		return
	} else {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead",
			token, p.curToken.Literal)
		fmt.Println(msg)
		os.Exit(1)
	}
}

func (p *Parser) expectTokenByLiteral(token string) {
	if p.curToken.Literal == token {
		return
	} else {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead",
			token, p.curToken.Literal)
		fmt.Println(msg)
		os.Exit(1)
	}
}
