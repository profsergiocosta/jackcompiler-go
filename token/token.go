package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// TOKEN TYPE
	IDENT    = "IDENTIFIER" // identifiers: add, foobar, x, y, ...
	INTCONST = "INTCONST"
	STRING   = "STRINGCONST"
	SYMBOL   = "SYMBOL"
	KEYWORD  = "KEYWORD"

	// Keywords
	METHOD      = "METHOD"
	STATIC      = "STATIC"
	INT         = "INT"
	BOOLEAN     = "BOOLEAN"
	TRUE        = "TRUE"
	NULL        = "NULL"
	LET         = "LET"
	IF          = "IF"
	WHILE       = "WHILE"
	CONSTRUCTOR = "CONSTRUCTOR"
	FIELD       = "FIELD"
	VAR         = "VAR"
	CHAR        = "CHAR"
	VOID        = "VOID"
	CLASS       = "CLASS"
	FALSE       = "FALSE"
	DO          = "DO"
	ELSE        = "ELSE"
	RETURN      = "RETURN"
	FUNCTION    = "FUNCTION"
	THIS        = "THIS"
)

var keywords = map[string]TokenType{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

func LookupIdent(ident string) TokenType {
	if _, ok := keywords[ident]; ok {
		return KEYWORD //tok
	}
	return IDENT
}
