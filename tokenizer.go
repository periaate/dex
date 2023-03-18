package dex

// Token is the set of tokens that the parser recognizes.
type Token int

// The set of tokens
const (
	// Special tokens
	_ Token = iota
	ILLEGAL
	EOF
	NEWLINE
	WHITESPACE

	IDENT

	// Delimiters
	LBRACE // {
	RBRACE // }

	LPAREN // (
	RPAREN // )

	// Operators
	STREAM // >
	APPLY  // <
)

// String returns a string representation of the token.
func (t Token) String() string {
	s := tokenToString[t]
	if s == "" {
		return "ILLEGAL"
	}
	return s
}

var tokenToString = [...]string{
	EOF:        "EOF",
	WHITESPACE: "WHITESPACE",

	IDENT: "IDENT",

	LBRACE: "{",
	RBRACE: "}",
	LPAREN: "(",
	RPAREN: ")",

	STREAM: ">",
	APPLY:  "<",
}

var runeToToken = map[rune]Token{
	'{': LBRACE,
	'}': RBRACE,
	'(': LPAREN,
	')': RPAREN,
	'>': STREAM,
	'<': APPLY,

	' ':  WHITESPACE,
	'\t': WHITESPACE,
	'\n': NEWLINE,

	'\r': ILLEGAL,
	'\b': ILLEGAL,
}

// Scanner is a simple tokenizer for the Dex dexguage.
type Scanner struct {
	src   string
	chPos int
}
type Pos struct {
	Start int
	End   int
}

func (s *Scanner) Next() (tok Token, lit string, pos Pos) {
	var (
		start int = s.chPos
		r     rune
		t     Token
		ok    bool
	)

	for !ok {
		if s.chPos >= len(s.src) {
			t = EOF
			break
		}
		r = rune(s.src[s.chPos])
		t, ok = runeToToken[r]
		s.chPos++
	}

	switch t {
	case WHITESPACE:
		if s.chPos-start > 1 {
			return IDENT, s.src[start : s.chPos-1], Pos{start, s.chPos - 1}
		}
		return s.Next()
	case EOF:
		if s.chPos-start > 1 {
			return IDENT, s.src[start:s.chPos], Pos{start, s.chPos}
		}
		return t, "", Pos{start, s.chPos}
	default:
		if s.chPos-start > 1 {
			s.chPos--
			return IDENT, s.src[start:s.chPos], Pos{start, s.chPos}
		}
		return t, s.src[start:s.chPos], Pos{start, s.chPos}
	}
}

func (s *Scanner) Peek() Token {
	chPos := s.chPos
	tok, _, _ := s.Next()
	s.chPos = chPos
	return tok
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src: src,
	}
}
