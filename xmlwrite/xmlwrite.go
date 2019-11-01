package xmlwrite

import (
	"fmt"
	"io/ioutil"

	"github.com/profsergiocosta/jackcompiler-go/token"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
)

func tagToken(tok token.Token) string {
	value := tok.Literal
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
	return fmt.Sprintf("<%s>%s</%s>", tok.Type, value, tok.Type)
}

func PrintTerminal(tok token.Token) {
	fmt.Println(tagToken(tok))
}

func TagNonTerminal(nonTerminal string) {
	fmt.Println("<" + nonTerminal + ">")
}

func UntagNonTerminal(nonTerminal string) {
	fmt.Println("</" + nonTerminal + ">")
}

func imprime() {
	/*
		input := `let five = 5;
		let ten = "10";
		 hoje  é + casa
		de sair para casa
		let add = x / y;
		let result = add(five, ten);`
	*/

	input, err := ioutil.ReadFile("Main.jack")
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
