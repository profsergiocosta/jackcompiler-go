package vmwriter

import (
	"fmt"
	"os"
)

type Command string

const (
	ADD Command = "add"
	SUB Command = "sub"
	NEG Command = "neg"
	EQ  Command = "eq"
	GT  Command = "gt"
	LT  Command = "lt"
	AND Command = "and"
	OR  Command = "or"
	NOT Command = "not"
)

type Segment string

const (
	STATIC  Segment = "static"
	FIELD   Segment = "field"
	ARG     Segment = "arg"
	LOCAL   Segment = "local"
	CONST   Segment = "const"
	THIS    Segment = "this"
	THAT    Segment = "that"
	POINTER Segment = "pointer"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type VMWriter struct {
	out *os.File
}

func New(pathName string) *VMWriter {
	f, err := os.Create(pathName)
	check(err)

	vm := &VMWriter{out: f}

	return vm
}

func (vm *VMWriter) WritePush(segment Segment, index int) {
	s := fmt.Sprintf("push %s %d\n", segment, index)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WritePop(segment Segment, index int) {
	s := fmt.Sprintf("pop %s %d\n", segment, index)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteArithmetic(command Command) {
	s := fmt.Sprintf("%s\n", command)
	vm.out.WriteString(s)
}

func (vm *VMWriter) CloseFile() {
	vm.out.Close()
}
