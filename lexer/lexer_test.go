package lexer

import (
	"fmt"
	"testing"

	"github.com/profsergiocosta/jackcompiler-go/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = x + y;
	let result = add(five, ten);`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.KEYWORD, "let"},
		{token.IDENT, "five"},
		{token.SYMBOL, "="},
		{token.INTCONST, "5"},
		{token.SYMBOL, ";"},
		{token.KEYWORD, "let"},
		{token.IDENT, "ten"},
		{token.SYMBOL, "="},
		{token.INTCONST, "10"},
		{token.SYMBOL, ";"},
		{token.KEYWORD, "let"},
		{token.IDENT, "add"},
		{token.SYMBOL, "="},
		{token.IDENT, "x"},
		{token.SYMBOL, "+"},
		{token.IDENT, "y"},
		{token.SYMBOL, ";"},
		{token.KEYWORD, "let"},
		{token.IDENT, "result"},
		{token.SYMBOL, "="},
		{token.IDENT, "add"},
		{token.SYMBOL, "("},
		{token.IDENT, "five"},
		{token.SYMBOL, ","},
		{token.IDENT, "ten"},
		{token.SYMBOL, ")"},
		{token.SYMBOL, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Println(tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
