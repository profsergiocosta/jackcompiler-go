package main

import (
	"fmt"
	"strings"

	"github.com/profsergiocosta/jackcompiler-go/token"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
)

// https://golang.org/doc/code.html

func isSymbol(c string) bool {
	symbols := "{}()[].,;+-*/&|<>=~"
	return strings.Index(symbols, c) != -1
}

func tagToken(tok token.Token) string {
	ttype := tok.Type
	if isSymbol(string(ttype)) {
		ttype = token.SYMBOL
	}

	return fmt.Sprintf("<%s>%s<%s>", ttype, tok.Literal, ttype)
}
func main() {
	input := `let five = 5;
	let ten = 10;
	let add = x + y;
	let result = add(five, ten);`

	l := lexer.New(input)
	fmt.Println("<tokens>")
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		//fmt.Printf("%+v\n", tok)
		fmt.Println(tagToken(tok))
	}
	fmt.Println("</tokens>")
}
