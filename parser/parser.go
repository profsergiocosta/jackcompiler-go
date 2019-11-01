package parser

import (
	"fmt"

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
	fmt.Println(p.errors)
}

func (p *Parser) CompileClass() {

	if !p.curTokenIs("class") {
		msg := fmt.Sprintf("expected token to be %s, got %s instead",
			"class", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return
	}

	if p.output == XML {
		xmlwrite.TagNonTerminal("class")
	}

	if p.output == XML {
		xmlwrite.UntagNonTerminal("class")
	}

}

// book write an interpreter in go

func (p *Parser) curTokenIs(t string) bool {
	return p.curToken.Literal == t
}

func (p *Parser) peekTokenIs(token string) bool {
	return p.peekToken.Literal == token
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

func (p *Parser) peekError(t string) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
