package dex_test

import (
	"dex"
	"testing"
)

type expectToken struct {
	tok dex.Token
	lit string
	pos dex.Pos
}

// TODO: more comprehensive tests
func TestTokenizer(t *testing.T) {
	tk := dex.NewScanner("set{1 2 3} > test")
	expect := []expectToken{
		{dex.IDENT, "set", dex.Pos{0, 3}},
		{dex.LBRACE, "{", dex.Pos{3, 4}},
		{dex.IDENT, "1", dex.Pos{4, 5}},
		{dex.IDENT, "2", dex.Pos{6, 7}},
		{dex.IDENT, "3", dex.Pos{8, 9}},
		{dex.RBRACE, "}", dex.Pos{9, 10}},
		{dex.STREAM, ">", dex.Pos{11, 12}},
		{dex.IDENT, "test", dex.Pos{13, 17}},
	}

	var i int
	for {
		tok, lit, _ := tk.Next()
		if tok == dex.EOF {
			break
		}
		if lit != expect[i].lit {
			t.Errorf("Expected %v, got %v", expect[i].lit, lit)
		}
		if tok != expect[i].tok {
			t.Errorf("Expected %v, got %v", expect[i].tok, tok)
		}
		i++
	}

	if i != len(expect) {
		t.Errorf("Expected length: %v, got %v", len(expect), i)
	}
}

func TestFnMaps(t *testing.T) {
	tk := dex.NewScanner("set(fn{test} fn2{test})")
	expect := []expectToken{
		{dex.IDENT, "set", dex.Pos{0, 3}},
		{dex.LPAREN, "(", dex.Pos{3, 4}},
		{dex.IDENT, "fn", dex.Pos{4, 6}},
		{dex.LBRACE, "{", dex.Pos{6, 7}},
		{dex.IDENT, "test", dex.Pos{7, 11}},
		{dex.RBRACE, "}", dex.Pos{11, 12}},
		{dex.IDENT, "fn2", dex.Pos{13, 16}},
		{dex.LBRACE, "{", dex.Pos{16, 17}},
		{dex.IDENT, "test", dex.Pos{17, 21}},
		{dex.RBRACE, "}", dex.Pos{21, 22}},
		{dex.RPAREN, ")", dex.Pos{22, 23}},
	}

	var i int
	for {
		tok, lit, _ := tk.Next()
		if tok == dex.EOF {
			break
		}
		if lit != expect[i].lit {
			t.Errorf("Expected %v, got %v", expect[i].lit, lit)
		}
		if tok != expect[i].tok {
			t.Errorf("Expected %v, got %v", expect[i].tok, tok)
		}
		i++
	}

	if i != len(expect) {
		t.Errorf("Expected length: %v, got %v", len(expect), i)
	}
}
