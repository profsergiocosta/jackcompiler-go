package symboltable

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
	numDefinitions map[string]int
}

func NewSymbolTable() *SymbolTable {
	s1 := make(map[string]Symbol)
	s2 := make(map[string]Symbol)
	numDef = map[string]int {
		STATIC : 0,
		FIELD : 0,
		ARG : 0,
		VAR : 0,
	}
	return &SymbolTable{classScope: s1, subRoutineScope: s2, numDefinitions: numDef}
}

func (s *SymbolTable) StartSubroutine() {
	s.subRoutineScope = make(map[string]Symbol)
	s.numDefinitions[ARG] = 0
	s.numDefinitions[VAR] = 0
}

func (s *SymbolTable) Define(string name, string ttype, SymbolScope scope) {

	if scope == STATIC || scope == FIELD {
		symbol := Symbol{Name: name, Type: ttype, Index: s.numDefinitions[scope], Scope: GlobalScope}
		s.classScope[name] = symbol
	} else {
		symbol := Symbol{Name: name, Type: ttype, Index: s.numDefinitions[scope], Scope: GlobalScope}
		s.subRoutineScope[name] = symbol
	}
	s.numDefinitions[SymbolScope]++
	fmt.Println(s.numDefinitions[SymbolScope])
}



func (s *SymbolTable) VarCount( SymbolScope scope) {

}

func (s *SymbolTable) kindOf(string name) {

}

func (s *SymbolTable) TypeOf(string name) {

}

func (s *SymbolTable) IndexOf(string name) {

}
/
