package xmlwrite

import (
	"fmt"
	"io/ioutil"

	"github.com/profsergiocosta/jackcompiler-go/token"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
)

func tagToken(tok token.Token) string {
	value := tok.Literal
	ttype := tok.Type

	if tok.Type == token.SYMBOL {
		switch tok.Literal {
		case "<":
			value = "&lt;"
		case ">":
			value = "&gt;"
		case "\\":
			value = "quot;"
		case "&":
			value = "&amp;"
		default:
			value = tok.Literal
		}
	}

	if token.IsSymbol(tok.Literal[0]) {
		ttype = "SYMBOL"
	} else if token.IsKeyword(tok.Literal) {
		ttype = "KEYWORD"
	}

	return fmt.Sprintf("<%s>%s</%s>", ttype, value, ttype)
}

func PrintTerminal(tok token.Token, toPrint bool) {
	if tok.Type == token.EOF {
		fmt.Println("</EOF>")
	} else if toPrint {
		fmt.Println(tagToken(tok))
	}
}

func PrintNonTerminal(nonTerminal string, toPrint bool) {
	if toPrint {
		fmt.Println("<" + nonTerminal + ">")
	}
}

func Imprime() {
	/*
		input := `let five = 5;
		let ten = "10";
		 hoje  é + casa
		de sair para casa
		let add = x / y;
		let result = add(five, ten);`
	*/

	input, err := ioutil.ReadFile("xmlwrite/Main.jack")
	if err != nil {
		panic("erro")
	}
	l := lexer.New(string(input))
	fmt.Println("<tokens>")
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		//fmt.Printf("%+f\n", tok)
		fmt.Println(tagToken(tok))
	}
	fmt.Println("</tokens>")
}
