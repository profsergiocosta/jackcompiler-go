package main

import (
	"fmt"
	"io/ioutil"

	"github.com/profsergiocosta/jackcompiler-go/token"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
)

// https://golang.org/doc/code.html

/*
    if (symbol() == '<')
        return "<symbol> &lt; </symbol>";
      else if (symbol() == '>')
        return "<symbol> &gt; </symbol>";
      else if (symbol() == '\"')
        return "<symbol> &quot; </symbol>";
      else if (symbol() == '&')
		return "<symbol> &amp; </symbol>";
*/

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
func main() {
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
