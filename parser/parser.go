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
	p.output = XML
	return p
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Compile() {
	p.nextToken()
	//p.CompileClass()
	p.CompileExpression()

}

func (p *Parser) CompileClass() {

	xmlwrite.TagNonTerminal("class", p.output == XML)

	p.expectPeek(token.CLASS)
	xmlwrite.PrintTerminal(p.curToken, p.output == XML)

	p.expectPeek(token.IDENT)
	xmlwrite.PrintTerminal(p.curToken, p.output == XML)

	p.expectPeek(token.LBRACE)
	xmlwrite.PrintTerminal(p.curToken, p.output == XML)

	xmlwrite.UntagNonTerminal("class", p.output == XML)

}

func (p *Parser) CompileExpression() {
	xmlwrite.TagNonTerminal("EXPRESSION", p.output == XML)
	p.CompileTerm()
	for !p.peekTokenIs(token.EOF) && token.IsOperator(p.peekToken.Literal[0]) {
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		p.CompileTerm()
	}
	xmlwrite.UntagNonTerminal("EXPRESSION", p.output == XML)
}

func (p *Parser) CompileTerm() {
	xmlwrite.TagNonTerminal("TERM", p.output == XML)
	switch p.peekToken.Type {
	case token.INTCONST, token.KEYWORD, token.STRING:
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	case token.IDENT:
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		switch p.peekToken.Type {
		case token.LBRACKET:
			p.nextToken()
			xmlwrite.PrintTerminal(p.curToken, p.output == XML)
			p.CompileExpression()
			p.expectPeek(token.RBRACKET)
			xmlwrite.PrintTerminal(p.curToken, p.output == XML)

		case token.LPAREN, token.DOT:
			p.CompileSubroutineCall()

		default:

		}

	case token.LPAREN:
		p.expectPeek(token.LPAREN)
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)

		p.CompileExpression()

		p.expectPeek(token.RPAREN)
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	case token.MINUS, token.NOT:
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		p.CompileTerm()

	default:

	}
	xmlwrite.TagNonTerminal("/TERM", p.output == XML)
}

func (p *Parser) CompileSubroutineCall() {
	// ainda vou precisar saber o nome da funcao
	p.nextToken()
	if p.curTokenIs(token.LPAREN) {
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		p.CompileExpression()
		p.expectPeek(token.RPAREN)
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	} else {
		xmlwrite.PrintTerminal(p.curToken, p.output == XML) // DOT

		p.expectPeek(token.IDENT)
		xmlwrite.PrintTerminal(p.curToken, p.output == XML) // ident

		p.expectPeek(token.LPAREN)
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)

		p.CompileExpression()
		p.expectPeek(token.RPAREN)

	}

}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) {
	if p.peekTokenIs(t) {
		p.nextToken()
	} else {
		p.peekError(t)
		fmt.Println(p.errors)
		os.Exit(1)
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

/*

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
*/
