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
	ARG     Segment = "argument"
	LOCAL   Segment = "local"
	CONST   Segment = "constant"
	THIS    Segment = "this"
	THAT    Segment = "that"
	POINTER Segment = "pointer"
	TEMP    Segment = "temp"
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

func (vm *VMWriter) WriteCall(name string, nArgs int) {
	s := fmt.Sprintf("call %s %d\n", name, nArgs)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteFunction(name string, nLocals int) {
	s := fmt.Sprintf("function %s %d\n", name, nLocals)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteReturn() {
	s := fmt.Sprintf("return\n")
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteLabel(label string) {
	s := fmt.Sprintf("label %s\n", label)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteGoto(label string) {
	s := fmt.Sprintf("goto %s\n", label)
	vm.out.WriteString(s)
}

func (vm *VMWriter) WriteIf(label string) {
	s := fmt.Sprintf("if-goto %s\n", label)
	vm.out.WriteString(s)
}

func (vm *VMWriter) CloseFile() {
	vm.out.Close()
}
