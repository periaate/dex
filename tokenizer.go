package lan

// TokenType is the set of tokens that the parser recognizes.
type TokenType int

func (t TokenType) String() string {
	return tokenToString[t]
}

// The set of tokens
const (
	// Special tokens
	EOF TokenType = iota
	WHITESPACE

	IDENT

	// Delimiters
	LBRACE // {
	RBRACE // }

	// Operators
	STREAM // >
	APPLY  // <
)

var tokenToString = [...]string{
	EOF:        "EOF",
	WHITESPACE: " ",

	IDENT: "IDENT",

	LBRACE: "{",
	RBRACE: "}",

	STREAM: ">",
	APPLY:  "<",
}

var runeToToken = map[rune]TokenType{
	'{': LBRACE,
	'}': RBRACE,
	'>': STREAM,
	'<': APPLY,

	' ':  WHITESPACE,
	'\n': EOF,
	'\t': EOF,
	'\r': EOF,
}

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	src string
	pos int
}

func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		src: src,
	}
}

func (t *Tokenizer) NextToken() Token {
	var (
		start int = t.pos
		r     rune
		tok   TokenType
		ok    bool
	)

	for !ok && t.pos < len(t.src) {
		r = rune(t.src[t.pos])
		tok, ok = runeToToken[r]
		t.pos++
	}

	switch tok {
	case WHITESPACE:
		if t.pos-start > 1 {
			return Token{IDENT, t.src[start : t.pos-1]}
		}
		return t.NextToken()
	case EOF:
		if t.pos-start > 1 {
			return Token{IDENT, t.src[start:t.pos]}
		}
		return Token{tok, ""}
	default:
		if t.pos-start > 1 {
			t.pos--
			return Token{IDENT, t.src[start:t.pos]}
		}
		return Token{tok, string(r)}
	}
}

func (t *Tokenizer) PeekToken() Token {
	pos := t.pos
	tok := t.NextToken()
	t.pos = pos
	return tok
}
