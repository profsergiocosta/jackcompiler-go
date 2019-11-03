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
	p.CompileClass()
}

func (p *Parser) CompileClass() {

	xmlwrite.PrintNonTerminal("class", p.output == XML)

	p.expectPeek(token.CLASS)

	p.expectPeek(token.IDENT)

	p.expectPeek(token.LBRACE)

	for p.peekTokenIs(token.STATIC) || p.peekTokenIs(token.FIELD) {
		p.CompileClassVarDec()
	}

	for p.peekTokenIs(token.FUNCTION) || p.peekTokenIs(token.CONSTRUCTOR) || p.peekTokenIs(token.METHOD) {
		p.CompileSubroutine()
	}

	p.expectPeek(token.RBRACE)

	xmlwrite.PrintNonTerminal("/class", p.output == XML)

}

func (p *Parser) CompileClassVarDec() {
	xmlwrite.PrintNonTerminal("classVarDec", p.output == XML)

	if p.peekTokenIs(token.FIELD) {
		p.expectPeek(token.FIELD)

	} else {
		p.expectPeek(token.STATIC)

	}

	p.CompileType()

	p.expectPeek(token.IDENT)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)

		p.expectPeek(token.IDENT)

	}

	p.expectPeek(token.SEMICOLON)

	xmlwrite.PrintNonTerminal("/classVarDec", p.output == XML)
}

func (p *Parser) CompileSubroutine() {
	xmlwrite.PrintNonTerminal("subroutineDec", p.output == XML)

	if p.peekTokenIs(token.CONSTRUCTOR) {
		p.expectPeek(token.CONSTRUCTOR)

	} else if p.peekTokenIs(token.FUNCTION) {
		p.expectPeek(token.FUNCTION)

	} else {
		p.expectPeek(token.METHOD)

	}

	if p.peekTokenIs(token.VOID) {
		p.expectPeek(token.VOID)

	} else {
		p.CompileType()
	}

	p.expectPeek(token.IDENT)

	p.expectPeek(token.LPAREN)

	if !p.peekTokenIs(token.RPAREN) {
		p.CompileParameterList()
	} else {
		// because of compare xml
		xmlwrite.PrintNonTerminal("parameterList", p.output == XML)
		xmlwrite.PrintNonTerminal("/parameterList", p.output == XML)
	}

	p.expectPeek(token.RPAREN)

	p.CompileSubroutineBody()

	xmlwrite.PrintNonTerminal("/subroutineDec", p.output == XML)
}

func (p *Parser) CompileSubroutineBody() {

	xmlwrite.PrintNonTerminal("subroutineBody", p.output == XML)

	p.expectPeek(token.LBRACE)

	for p.peekTokenIs(token.VAR) {
		p.CompileVarDec()
	}

	p.CompileStatements()

	p.expectPeek(token.RBRACE)

	xmlwrite.PrintNonTerminal("/subroutineBody", p.output == XML)
}

func (p *Parser) CompileVarDec() {
	xmlwrite.PrintNonTerminal("varDec", p.output == XML)

	p.expectPeek(token.VAR)

	p.CompileType()

	p.expectPeek(token.IDENT)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)

		p.expectPeek(token.IDENT)

	}

	p.expectPeek(token.SEMICOLON)

	xmlwrite.PrintNonTerminal("/varDec", p.output == XML)
}

func (p *Parser) CompileParameterList() {
	xmlwrite.PrintNonTerminal("parameterList", p.output == XML)
	p.CompileType()
	p.expectPeek(token.IDENT)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)

		p.CompileType()

		p.expectPeek(token.IDENT)
	}

	xmlwrite.PrintNonTerminal("/parameterList", p.output == XML)
}

func (p *Parser) CompileType() {
	switch p.peekToken.Type {
	case token.INT:
		p.expectPeek(token.INT)

	case token.CHAR:
		p.expectPeek(token.CHAR)

	case token.BOOLEAN:
		p.expectPeek(token.BOOLEAN)

	case token.IDENT:
		p.expectPeek(token.IDENT)

	}
}

func (p *Parser) CompileExpression() {
	xmlwrite.PrintNonTerminal("EXPRESSION", p.output == XML)
	p.CompileTerm()

	for !p.peekTokenIs(token.EOF) && token.IsOperator(p.peekToken.Literal[0]) {
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		p.CompileTerm()
	}

	xmlwrite.PrintNonTerminal("/EXPRESSION", p.output == XML)
}

func (p *Parser) CompileTerm() {
	xmlwrite.PrintNonTerminal("TERM", p.output == XML)
	switch p.peekToken.Type {
	case token.INTCONST, token.TRUE, token.FALSE, token.NULL, token.THIS, token.STRING:
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	case token.IDENT:
		p.expectPeek(token.IDENT)

		switch p.peekToken.Type {
		case token.LBRACKET:
			p.expectPeek(token.LBRACKET)

			p.CompileExpression()

			p.expectPeek(token.RBRACKET)

		case token.LPAREN, token.DOT:
			p.CompileSubroutineCall()

		default:

		}

	case token.LPAREN:
		p.expectPeek(token.LPAREN)

		p.CompileExpression()

		p.expectPeek(token.RPAREN)

	case token.MINUS, token.NOT:
		p.nextToken()
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
		p.CompileTerm()

	default:
		//fmt.Println(p.peekToken)
		fmt.Println("operdor unario nao reconhecido")
		os.Exit(1)

	}
	xmlwrite.PrintNonTerminal("/TERM", p.output == XML)
}

func (p *Parser) CompileSubroutineCall() {
	// ainda vou precisar saber o nome da funcao
	if p.peekTokenIs(token.LPAREN) {
		p.expectPeek(token.LPAREN)

		p.CompileExpressionList()

		p.expectPeek(token.RPAREN)

	} else {
		p.expectPeek(token.DOT)
		p.expectPeek(token.IDENT)

		p.expectPeek(token.LPAREN)

		p.CompileExpressionList()

		p.expectPeek(token.RPAREN)

	}

}

func (p *Parser) CompileExpressionList() {
	xmlwrite.PrintNonTerminal("ExpressionList", p.output == XML)
	if !p.peekTokenIs(token.RPAREN) {
		p.CompileExpression()
	}

	for p.peekTokenIs(token.COMMA) {

		p.expectPeek(token.COMMA)

		p.CompileExpression()
	}

	xmlwrite.PrintNonTerminal("/ExpressionList", p.output == XML)

}

func (p *Parser) CompileDo() {
	xmlwrite.PrintNonTerminal("doStatement", p.output == XML)
	p.expectPeek(token.DO)

	p.expectPeek(token.IDENT)

	p.CompileSubroutineCall()
	p.expectPeek(token.SEMICOLON)

	xmlwrite.PrintNonTerminal("/doStatement", p.output == XML)
}

func (p *Parser) CompileWhile() {
	xmlwrite.PrintNonTerminal("whileStatement", p.output == XML)
	p.expectPeek(token.WHILE)

	p.expectPeek(token.LPAREN)

	p.CompileExpression()

	p.expectPeek(token.RPAREN)

	p.expectPeek(token.LBRACE)

	p.CompileStatements()

	p.expectPeek(token.RBRACE)

	xmlwrite.PrintNonTerminal("/whileStatement", p.output == XML)
}

func (p *Parser) CompileIf() {
	xmlwrite.PrintNonTerminal("ifStatement", p.output == XML)
	p.expectPeek(token.IF)

	p.expectPeek(token.LPAREN)

	p.CompileExpression()

	p.expectPeek(token.RPAREN)

	p.expectPeek(token.LBRACE)

	p.CompileStatements()

	p.expectPeek(token.RBRACE)

	if p.peekTokenIs(token.ELSE) {
		p.expectPeek(token.ELSE)

		p.expectPeek(token.LBRACE)

		p.CompileStatements()

		p.expectPeek(token.RBRACE)

	}

	xmlwrite.PrintNonTerminal("/ifStatement", p.output == XML)
}

func (p *Parser) CompileReturn() {
	xmlwrite.PrintNonTerminal("returnStatement", p.output == XML)

	p.expectPeek(token.RETURN)

	if !p.peekTokenIs(token.SEMICOLON) {
		p.CompileExpression()
	}

	p.expectPeek(token.SEMICOLON)

	xmlwrite.PrintNonTerminal("/returnStatement", p.output == XML)
}

func (p *Parser) CompileLet() {

	xmlwrite.PrintNonTerminal("letStatement", p.output == XML)

	p.expectPeek(token.LET)

	p.expectPeek(token.IDENT)

	if p.peekTokenIs(token.LBRACKET) {
		p.expectPeek(token.LBRACKET)

		p.CompileExpression()

		p.expectPeek(token.RBRACKET)

	}

	p.expectPeek(token.EQ)

	p.CompileExpression()

	p.expectPeek(token.SEMICOLON)

	xmlwrite.PrintNonTerminal("/letStatement", p.output == XML)
}
func (p *Parser) CompileStatements() {
	xmlwrite.PrintNonTerminal("statements", p.output == XML)
	p.CompileStatement()
	xmlwrite.PrintNonTerminal("/statements", p.output == XML)
}
func (p *Parser) CompileStatement() {

	switch p.peekToken.Type {
	case token.LET:
		p.CompileLet()
		p.CompileStatement()
	case token.DO:
		p.CompileDo()
		p.CompileStatement()
	case token.IF:
		p.CompileIf()
		p.CompileStatement()
	case token.WHILE:
		p.CompileWhile()
		p.CompileStatement()
	case token.RETURN:
		p.CompileReturn()
		p.CompileStatement()
	default:
		return
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
		xmlwrite.PrintTerminal(p.curToken, p.output == XML)
	} else {
		p.peekError(t, p.peekToken.Line)
		fmt.Println(p.errors)
		os.Exit(1)
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType, line int) {
	msg := fmt.Sprintf(" %v: expected next token to be %s, got %s instead",
		line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
