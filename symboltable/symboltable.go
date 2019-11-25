package symboltable

import (
	"fmt"
	"os"
)

type SymbolScope string

const (
	STATIC SymbolScope = "STATIC"
	FIELD  SymbolScope = "FIELD"
	ARG    SymbolScope = "ARG"
	VAR    SymbolScope = "VAR"
)

type Symbol struct {
	Name  string
	Type  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	classScope      map[string]Symbol
	subRoutineScope map[string]Symbol
	//numDefinitions  int
	numDefinitions map[SymbolScope]int
}

func NewSymbolTable() *SymbolTable {
	s1 := make(map[string]Symbol)
	s2 := make(map[string]Symbol)
	defs := map[SymbolScope]int{
		"STATIC": 0,
		"FIELD":  0,
		"ARG":    0,
		"VAR":    0,
	}
	return &SymbolTable{classScope: s1, subRoutineScope: s2, numDefinitions: defs}
}
func (s *SymbolTable) Lookup(name string) (Symbol, bool) {
	sym, hasSub := s.subRoutineScope[name]
	if hasSub {
		return sym, hasSub
	} else {
		sym, hasSub := s.classScope[name]
		return sym, hasSub
	}
}

func (s *SymbolTable) StartSubroutine() {
	s.subRoutineScope = make(map[string]Symbol)
	s.numDefinitions["ARG"] = 0
	s.numDefinitions["VAR"] = 0
}

func (s *SymbolTable) Define(name string, ttype string, scope SymbolScope) {

	if scope == STATIC || scope == FIELD {
		symbol := Symbol{Name: name, Type: ttype, Scope: scope, Index: s.numDefinitions[scope]}
		s.classScope[name] = symbol
	} else {
		symbol := Symbol{Name: name, Type: ttype, Scope: scope, Index: s.numDefinitions[scope]}
		s.subRoutineScope[name] = symbol
	}
	s.numDefinitions[scope]++

}

func (s *SymbolTable) Resolve(name string) Symbol {
	sym, hasDefined := s.Lookup(name)

	if !hasDefined {
		fmt.Printf("identifier %s not defined \n", name)
		os.Exit(1)
	}

	return sym

}

func (s *SymbolTable) VarCount(scope SymbolScope) int {
	return s.numDefinitions[scope]
}
