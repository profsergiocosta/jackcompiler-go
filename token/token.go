package token

import "strings"

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// TOKEN TYPE
	IDENT    = "IDENTIFIER" // identifiers: add, foobar, x, y, ...
	INTCONST = "INTCONST"
	STRING   = "STRINGCONST"

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

var symbols = map[byte]TokenType{
	'=': EQ,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'&': AND,
	'|': OR,
	'~': NOT,
	'.': DOT,
	'<': LT,
	'>': GT,
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func LookupSymbol(symbol byte) TokenType {
	if tok, ok := symbols[symbol]; ok {
		return tok
	}
	return ILLEGAL
}

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsSymbol(c byte) bool {
	symbols := "{}()[].,;+-*/&|<>=~"
	return strings.IndexByte(symbols, c) != -1
}

func IsOperator(c byte) bool {
	operators := "+-*/&|<>=~"
	return strings.IndexByte(operators, c) != -1
}

func IsKeyword(k string) bool {
	if _, ok := keywords[k]; ok {
		return true
	}
	return false
}
