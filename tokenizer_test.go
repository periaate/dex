package lan_test

import (
	"lan"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tk := lan.NewTokenizer("set{1 2 3} > test")
	expect := []lan.Token{
		{lan.IDENT, "set"},
		{lan.LBRACE, "{"},
		{lan.IDENT, "1"},
		{lan.IDENT, "2"},
		{lan.IDENT, "3"},
		{lan.RBRACE, "}"},
		{lan.STREAM, ">"},
		{lan.IDENT, "test"},
	}

	var tok lan.Token
	var i int
	for {
		tok = tk.NextToken()
		if tok.Type == lan.EOF {
			break
		}
		// fmt.Printf("%v: '%v'\n", tok.Type.String(), tok.Value)
		if tok.Value != expect[i].Value {
			t.Errorf("Expected %v, got %v", expect[i], tok.Value)
		}
		if tok.Type != expect[i].Type {
			t.Errorf("Expected %v, got %v", expect[i], tok.Type)
		}
		i++
	}

	if i != len(expect) {
		t.Errorf("Expected %v tokens, got %v", len(expect), i)
	}
}
