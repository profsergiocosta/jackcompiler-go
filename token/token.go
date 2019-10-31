package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT    = "IDENT" // add, foobar, x, y, ...
	INTCONST = "INTCOST"
	STRING   = "STRING"
	SYMBOL   = "SYMBOL"

	//  Symbols
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	AND      = "&"
	OR       = "|"
	NOT      = "~"
	DOT      = "."
	LT       = "<"
	GT       = ">"
	EQ       = "="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

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
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
