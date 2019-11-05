package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/profsergiocosta/jackcompiler-go/symboltable"
	"github.com/profsergiocosta/jackcompiler-go/vmwriter"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
	"github.com/profsergiocosta/jackcompiler-go/token"
	xmlwriter "github.com/profsergiocosta/jackcompiler-go/xmlwrite"
)

const (
	XML = "XML"
	VM  = "VM"
)

type Parser struct {
	l             *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	output        string
	errors        []string
	st            *symboltable.SymbolTable
	className     string
	vm            *vmwriter.VMWriter
	fileName      string
	whileLabelNum int
	ifLabelNum    int
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func New(pathName string) *Parser {

	input, err := ioutil.ReadFile(pathName)
	if err != nil {
		panic("erro")
	}
	l := lexer.New(string(input))

	p := &Parser{l: l}
	p.fileName = FilenameWithoutExtension(path.Base(pathName))

	p.output = VM
	p.st = symboltable.NewSymbolTable()
	p.vm = vmwriter.New(FilenameWithoutExtension(pathName) + ".vm1")
	p.nextToken()
	p.whileLabelNum = 0
	p.ifLabelNum = 0
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

	xmlwriter.PrintNonTerminal("class", p.output == XML)

	p.expectPeek(token.CLASS)

	p.expectPeek(token.IDENT)
	p.className = p.curToken.Literal

	// criar um pacote para erro
	if p.className != p.fileName {
		fmt.Printf("In %s.jack (line %d): The class name doesn't match the file name\n", p.fileName, p.curToken.Line)
		os.Exit(1)
	}

	p.expectPeek(token.LBRACE)

	for p.peekTokenIs(token.STATIC) || p.peekTokenIs(token.FIELD) {
		p.CompileClassVarDec()
	}

	for p.peekTokenIs(token.FUNCTION) || p.peekTokenIs(token.CONSTRUCTOR) || p.peekTokenIs(token.METHOD) {
		p.CompileSubroutine()
	}

	p.expectPeek(token.RBRACE)

	xmlwriter.PrintNonTerminal("/class", p.output == XML)

}

func (p *Parser) CompileClassVarDec() {
	xmlwriter.PrintNonTerminal("classVarDec", p.output == XML)

	var scope symboltable.SymbolScope

	if p.peekTokenIs(token.FIELD) {
		p.expectPeek(token.FIELD)
		scope = token.FIELD
	} else {
		p.expectPeek(token.STATIC)
		scope = token.STATIC
	}

	p.CompileType()
	ttype := p.curToken.Literal

	p.expectPeek(token.IDENT)

	name := p.curToken.Literal

	p.st.Define(name, ttype, scope)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)
		p.expectPeek(token.IDENT)

		name := p.curToken.Literal
		p.st.Define(name, ttype, scope)
	}

	p.expectPeek(token.SEMICOLON)

	xmlwriter.PrintNonTerminal("/classVarDec", p.output == XML)
}

func (p *Parser) CompileSubroutine() {
	p.st.StartSubroutine()

	xmlwriter.PrintNonTerminal("subroutineDec", p.output == XML)

	if p.peekTokenIs(token.CONSTRUCTOR) {
		p.expectPeek(token.CONSTRUCTOR)

	} else if p.peekTokenIs(token.FUNCTION) {
		p.expectPeek(token.FUNCTION)

	} else {
		p.expectPeek(token.METHOD)
		p.st.Define("this", p.className, symboltable.ARG)
	}

	subroutineType := p.curToken.Type

	if p.peekTokenIs(token.VOID) {
		p.expectPeek(token.VOID)

	} else {
		p.CompileType()
	}

	p.expectPeek(token.IDENT)

	functionName := p.className + "." + p.curToken.Literal

	p.expectPeek(token.LPAREN)

	if !p.peekTokenIs(token.RPAREN) {
		p.CompileParameterList()
	} else {
		// because of compare xml
		xmlwriter.PrintNonTerminal("parameterList", p.output == XML)
		xmlwriter.PrintNonTerminal("/parameterList", p.output == XML)
	}

	p.expectPeek(token.RPAREN)

	p.CompileSubroutineBody(functionName, subroutineType)

	xmlwriter.PrintNonTerminal("/subroutineDec", p.output == XML)
}

func (p *Parser) CompileSubroutineBody(functionName string, subroutineType token.TokenType) {

	xmlwriter.PrintNonTerminal("subroutineBody", p.output == XML)

	p.expectPeek(token.LBRACE)

	for p.peekTokenIs(token.VAR) {
		p.CompileVarDec()
	}

	nLocals := p.st.VarCount(symboltable.VAR)

	p.vm.WriteFunction(functionName, nLocals)

	switch subroutineType {
	case token.CONSTRUCTOR:
		p.vm.WritePush(vmwriter.CONST, p.st.VarCount(symboltable.FIELD))
		p.vm.WriteCall("Memory.alloc", 1)
		p.vm.WritePop(vmwriter.POINTER, 0)

	case token.METHOD:
		p.vm.WritePush(vmwriter.ARG, 0)
		p.vm.WritePop(vmwriter.POINTER, 0)
	}

	p.CompileStatements()

	p.expectPeek(token.RBRACE)

	xmlwriter.PrintNonTerminal("/subroutineBody", p.output == XML)
}

func (p *Parser) CompileVarDec() {
	xmlwriter.PrintNonTerminal("varDec", p.output == XML)

	p.expectPeek(token.VAR)

	var scope symboltable.SymbolScope = symboltable.VAR

	p.CompileType()
	ttype := p.curToken.Literal

	p.expectPeek(token.IDENT)

	name := p.curToken.Literal

	p.st.Define(name, ttype, scope)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)

		p.expectPeek(token.IDENT)
		name := p.curToken.Literal

		p.st.Define(name, ttype, scope)
	}

	p.expectPeek(token.SEMICOLON)

	xmlwriter.PrintNonTerminal("/varDec", p.output == XML)
}

func (p *Parser) CompileParameterList() {
	xmlwriter.PrintNonTerminal("parameterList", p.output == XML)

	var scope symboltable.SymbolScope = symboltable.ARG

	p.CompileType()
	ttype := p.curToken.Literal

	p.expectPeek(token.IDENT)
	name := p.curToken.Literal

	p.st.Define(name, ttype, scope)

	for p.peekTokenIs(token.COMMA) {
		p.expectPeek(token.COMMA)

		p.CompileType()
		ttype := p.curToken.Literal

		p.expectPeek(token.IDENT)

		name := p.curToken.Literal
		p.st.Define(name, ttype, scope)
	}

	xmlwriter.PrintNonTerminal("/parameterList", p.output == XML)
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
	xmlwriter.PrintNonTerminal("EXPRESSION", p.output == XML)
	p.CompileTerm()

	for !p.peekTokenIs(token.EOF) && token.IsOperator(p.peekToken.Literal[0]) {
		p.nextToken()
		xmlwriter.PrintTerminal(p.curToken, p.output == XML)

		op := p.curToken.Type

		p.CompileTerm()

		p.CompileOperators(op)
	}

	xmlwriter.PrintNonTerminal("/EXPRESSION", p.output == XML)
}

func (p *Parser) CompileOperators(op token.TokenType) {

	switch op {
	case token.PLUS:
		p.vm.WriteArithmetic(vmwriter.ADD)
	case token.MINUS:
		p.vm.WriteArithmetic(vmwriter.SUB)
	case token.ASTERISK:
		p.vm.WriteCall("Math.multiply", 2)
	case token.SLASH:
		p.vm.WriteCall("Math.divide", 2)
	case token.AND:
		p.vm.WriteArithmetic(vmwriter.AND)
	case token.OR:
		p.vm.WriteArithmetic(vmwriter.OR)
	case token.LT:
		p.vm.WriteArithmetic(vmwriter.LT)
	case token.GT:
		p.vm.WriteArithmetic(vmwriter.GT)
	case token.EQ:
		p.vm.WriteArithmetic(vmwriter.EQ)
	case token.NOT:
		p.vm.WriteArithmetic(vmwriter.NOT)

	}

}

func (p *Parser) CompileKeywordConst(v token.TokenType) {
	switch v {
	case token.TRUE:
		p.expectPeek(token.TRUE)
		p.vm.WritePush(vmwriter.CONST, 0)
		p.vm.WriteArithmetic(vmwriter.NOT) // false is -1
	case token.FALSE, token.NULL:
		p.nextToken()
		xmlwriter.PrintTerminal(p.curToken, p.output == XML)
		p.vm.WritePush(vmwriter.CONST, 0)
	case token.THIS:
		p.expectPeek(token.THIS)
		p.vm.WritePush(vmwriter.POINTER, 0)

	}
}
func (p *Parser) CompileTerm() {
	xmlwriter.PrintNonTerminal("TERM", p.output == XML)
	switch p.peekToken.Type {
	case token.INTCONST:
		p.nextToken()
		p.vm.WritePush(vmwriter.CONST, p.curTokenAsInt())
		xmlwriter.PrintTerminal(p.curToken, p.output == XML)
	case token.TRUE, token.FALSE, token.NULL, token.THIS:
		p.CompileKeywordConst(p.peekToken.Type)

	case token.STRING:
		p.expectPeek(token.STRING)
		xmlwriter.PrintTerminal(p.curToken, p.output == XML)

		str := p.curToken.Literal
		p.vm.WritePush(vmwriter.CONST, len(str))
		p.vm.WriteCall("String.new", 1)
		for i := 0; i < len(str); i++ {
			p.vm.WritePush(vmwriter.CONST, int(str[i]))
			p.vm.WriteCall("String.appendChar", 2)
		}

	case token.IDENT:
		p.expectPeek(token.IDENT)
		identName := p.curToken.Literal

		switch p.peekToken.Type {
		case token.LBRACKET: // is array
			p.expectPeek(token.LBRACKET)

			p.CompileExpression()

			sym := p.st.Resolve(identName)
			p.vm.WritePush(scopeToSegment(sym.Scope), sym.Index)
			p.vm.WriteArithmetic(vmwriter.ADD)

			p.expectPeek(token.RBRACKET)
			p.vm.WritePop(vmwriter.POINTER, 1)
			p.vm.WritePush(vmwriter.THAT, 0)

		case token.LPAREN, token.DOT: //is subroutine
			p.CompileSubroutineCall()
			//os.Exit(0)

		default: // is variable
			sym := p.st.Resolve(identName)
			p.vm.WritePush(scopeToSegment(sym.Scope), sym.Index)
		}

	case token.LPAREN:
		p.expectPeek(token.LPAREN)

		p.CompileExpression()

		p.expectPeek(token.RPAREN)

	case token.MINUS, token.NOT:
		p.nextToken()
		op := p.curToken.Type

		xmlwriter.PrintTerminal(p.curToken, p.output == XML)
		p.CompileTerm()
		if op == token.MINUS {
			p.vm.WriteArithmetic(vmwriter.NEG)
		} else {
			p.vm.WriteArithmetic(vmwriter.NOT)
		}

	default:
		fmt.Printf("%v %v unario operator not recognized\n", p.peekToken, p.curToken)
		os.Exit(1)

	}
	xmlwriter.PrintNonTerminal("/TERM", p.output == XML)
}

func (p *Parser) CompileSubroutineCall() {
	// ainda vou precisar saber o nome da funcao
	ident := p.curToken.Literal
	numArgs := 0
	if p.peekTokenIs(token.LPAREN) { // é um método
		p.expectPeek(token.LPAREN)

		numArgs = p.CompileExpressionList()
		p.vm.WritePush(vmwriter.POINTER, 0)
		p.expectPeek(token.RPAREN)

		ident = p.className + "." + ident
		numArgs++
		p.vm.WriteCall(ident, numArgs)

	} else {
		p.expectPeek(token.DOT)
		sym, has := p.st.Lookup(ident)

		if has { // method
			p.vm.WritePush(scopeToSegment(sym.Scope), sym.Index)
			p.expectPeek(token.IDENT)
			ident = sym.Type + "." + p.curToken.Literal
			p.expectPeek(token.LPAREN)

			numArgs = p.CompileExpressionList()
			numArgs++
		} else { // static
			p.expectPeek(token.IDENT)
			ident = ident + "." + p.curToken.Literal
			p.expectPeek(token.LPAREN)

			numArgs = p.CompileExpressionList()
		}

		p.expectPeek(token.RPAREN)
		p.vm.WriteCall(ident, numArgs)
	}

}

func (p *Parser) CompileExpressionList() int {
	xmlwriter.PrintNonTerminal("ExpressionList", p.output == XML)
	numArgs := 0

	if !p.peekTokenIs(token.RPAREN) {
		p.CompileExpression()
		numArgs++
	}

	for p.peekTokenIs(token.COMMA) {

		p.expectPeek(token.COMMA)

		p.CompileExpression()
		numArgs++
	}

	xmlwriter.PrintNonTerminal("/ExpressionList", p.output == XML)
	return numArgs
}

func (p *Parser) CompileDo() {
	xmlwriter.PrintNonTerminal("doStatement", p.output == XML)
	p.expectPeek(token.DO)

	p.expectPeek(token.IDENT)

	p.CompileSubroutineCall()
	p.expectPeek(token.SEMICOLON)
	p.vm.WritePop(vmwriter.TEMP, 0)

	xmlwriter.PrintNonTerminal("/doStatement", p.output == XML)
}

func (p *Parser) CompileWhile() {
	xmlwriter.PrintNonTerminal("whileStatement", p.output == XML)

	labelWhileExp := fmt.Sprintf("WHILE_EXP%d", p.whileLabelNum)
	labelWhileEnd := fmt.Sprintf("WHILE_END%d", p.whileLabelNum)
	p.whileLabelNum++

	p.vm.WriteLabel(labelWhileExp)
	p.expectPeek(token.WHILE)

	p.expectPeek(token.LPAREN)

	p.CompileExpression()
	p.vm.WriteArithmetic(vmwriter.NOT)
	p.vm.WriteIf(labelWhileEnd)

	p.expectPeek(token.RPAREN)

	p.expectPeek(token.LBRACE)

	p.CompileStatements()

	p.vm.WriteGoto(labelWhileExp)
	p.vm.WriteLabel(labelWhileEnd)

	p.expectPeek(token.RBRACE)

	xmlwriter.PrintNonTerminal("/whileStatement", p.output == XML)
}

/*
push constant 60
pop local 1
goto IF_END0
label IF_FALSE0
push constant 10
pop local 2
label IF_END0
*/

func (p *Parser) CompileIf() {
	xmlwriter.PrintNonTerminal("ifStatement", p.output == XML)

	labelTrue := fmt.Sprintf("IF_TRUE%d", p.ifLabelNum)
	labelFalse := fmt.Sprintf("IF_FALSE%d", p.ifLabelNum)
	labelEnd := fmt.Sprintf("IF_END%d", p.ifLabelNum)
	p.ifLabelNum++

	p.expectPeek(token.IF)
	p.expectPeek(token.LPAREN)
	p.CompileExpression()
	p.expectPeek(token.RPAREN)

	p.vm.WriteIf(labelTrue)
	p.vm.WriteGoto(labelFalse)
	p.vm.WriteLabel(labelTrue)

	p.expectPeek(token.LBRACE)

	p.CompileStatements()

	p.vm.WriteGoto(labelEnd)
	p.expectPeek(token.RBRACE)

	p.vm.WriteLabel(labelFalse)
	if p.peekTokenIs(token.ELSE) {
		p.expectPeek(token.ELSE)

		p.expectPeek(token.LBRACE)

		p.CompileStatements()

		p.expectPeek(token.RBRACE)

	}
	p.vm.WriteLabel(labelEnd)

	xmlwriter.PrintNonTerminal("/ifStatement", p.output == XML)
}

func (p *Parser) CompileReturn() {
	xmlwriter.PrintNonTerminal("returnStatement", p.output == XML)

	p.expectPeek(token.RETURN)

	if !p.peekTokenIs(token.SEMICOLON) {
		p.CompileExpression()
		p.vm.WriteReturn()
	} else {
		p.vm.WritePush(vmwriter.CONST, 0)
		p.vm.WriteReturn()
	}

	p.expectPeek(token.SEMICOLON)

	xmlwriter.PrintNonTerminal("/returnStatement", p.output == XML)
}

func (p *Parser) CompileLet() {

	isArray := false
	xmlwriter.PrintNonTerminal("letStatement", p.output == XML)

	p.expectPeek(token.LET)

	p.expectPeek(token.IDENT)

	varName := p.curToken.Literal
	sym := p.st.Resolve(varName)

	if p.peekTokenIs(token.LBRACKET) {
		p.expectPeek(token.LBRACKET)

		p.CompileExpression()
		p.vm.WritePush(scopeToSegment(sym.Scope), sym.Index)

		p.vm.WriteArithmetic(vmwriter.ADD)
		p.expectPeek(token.RBRACKET)
		isArray = true
	}

	p.expectPeek(token.EQ)

	p.CompileExpression()

	if isArray {
		p.vm.WritePop(vmwriter.TEMP, 0)
		p.vm.WritePop(vmwriter.POINTER, 1)
		p.vm.WritePush(vmwriter.TEMP, 0)
		p.vm.WritePop(vmwriter.THAT, 0)
	} else {
		p.vm.WritePop(scopeToSegment(sym.Scope), sym.Index)
	}

	p.expectPeek(token.SEMICOLON)

	xmlwriter.PrintNonTerminal("/letStatement", p.output == XML)
}
func (p *Parser) CompileStatements() {
	xmlwriter.PrintNonTerminal("statements", p.output == XML)
	p.CompileStatement()
	xmlwriter.PrintNonTerminal("/statements", p.output == XML)
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
		xmlwriter.PrintTerminal(p.curToken, p.output == XML)
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

//// auxiliars

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *Parser) curTokenAsInt() int {
	i1, err := strconv.Atoi(p.curToken.Literal)
	check(err)
	return i1
}

func scopeToSegment(scope symboltable.SymbolScope) vmwriter.Segment {
	switch scope {
	case symboltable.STATIC:
		return vmwriter.STATIC
	case symboltable.FIELD:
		return vmwriter.THIS
	case symboltable.VAR:
		return vmwriter.LOCAL
	case symboltable.ARG:
		return vmwriter.ARG
	default:

	}
	panic("scope undefined")
}
