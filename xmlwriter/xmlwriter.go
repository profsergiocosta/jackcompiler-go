package xmlwriter

import (
	"fmt"
	"io/ioutil"

	"jackcompiler-go/token"

	"jackcompiler-go/lexer"
)

func tagToken(tok token.Token) string {
	value := tok.Literal
	ttype := tok.Type

	if token.IsSymbol(tok.Literal[0]) {
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

func PrintAll(path string) {
	input, err := ioutil.ReadFile(path)
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
