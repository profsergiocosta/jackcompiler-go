package main

import (
	"fmt"

	"github.com/profsergiocosta/jackcompiler-go/token"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
)

// https://golang.org/doc/code.html

func tagToken(tok token.Token) string {
	return fmt.Sprintf("<%s>%s<%s>", tok.Type, tok.Literal, tok.Type)
}
func main() {
	input := `let five = 5;
	let ten = "10";
	let add = x / y;  
	let result = add(five, ten);`

	l := lexer.New(input)
	fmt.Println("<tokens>")
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
		fmt.Println(tagToken(tok))
	}
	fmt.Println("</tokens>")
}
