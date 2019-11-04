package vmwriter

import (
	"fmt"
	"os"
)

/*
enum Segment { CONST, ARG, LOCAL, STATIC, THIS, THAT, POINTER, TEMP };
enum Command { ADD, SUB, NEG, EQ, GT, LT, AND, OR, NOT };
*/

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

func (vm *VMWriter) WriteArithmetic(command string) {
	s := fmt.Sprintf("%s\n", command)
	vm.out.WriteString(s)
}

func (vm *VMWriter) CloseFile() {
	vm.out.Close()
}
