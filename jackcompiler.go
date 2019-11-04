package main

import (
	"os"

	"github.com/profsergiocosta/jackcompiler-go/parser"
)

// https://golang.org/doc/code.html

func main() {

	arg := os.Args[1:]

	//xmlwrite.PrintAll(arg[0])

	p := parser.New(arg[0])
	p.Compile()
	//p.CompileExpression()

}
