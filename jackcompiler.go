package main

import (
	"io/ioutil"

	"github.com/profsergiocosta/jackcompiler-go/lexer"
	"github.com/profsergiocosta/jackcompiler-go/parser"
)

// https://golang.org/doc/code.html

func main() {
	//fmt.Println("Hello, world.")
	input, err := ioutil.ReadFile("xmlwrite/Main.jack")
	if err != nil {
		panic("erro")
	}
	l := lexer.New(string(input))
	p := parser.New(l)
	p.Compile()

}
