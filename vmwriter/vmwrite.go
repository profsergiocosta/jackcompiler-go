package vmwriter

import (
	"fmt"
	"os"
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

func (vm *VMWriter) WritePush(segment string, index int) {
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
