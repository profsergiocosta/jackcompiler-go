package main

import (
	"io/ioutil"
	"os"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
	"github.com/profsergiocosta/jackcompiler-go/parser"
)

// https://golang.org/doc/code.html

func main() {

	//xmlwrite.Imprime()

	arg := os.Args[1:]

	input, err := ioutil.ReadFile(arg[0])
	if err != nil {
		panic("erro")
	}
	l := lexer.New(string(input))
	p := parser.New(l)
	p.Compile()
	//p.CompileExpression()

}
