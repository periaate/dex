package dex

type EvalType int

const (
	_          EvalType = iota
	Expression          // Evaluate
	Statement           // Do not evaluate
)

// Parser implements the necessary logic to parse dex.
// The parser is implemented as recursive descent parser.
// Parser implements the Node interface.
type Parser struct {
	t   *Scanner
	ast *AST
	tar *AST

	S *scopeMap

	tok     Token
	lastTok Token

	lit     string
	lastLit string
}

// NewParser returns a new Parser from the given scopeMap.
// If nil is passed the function will instantiate a new scopeMap.
func NewParser(s *scopeMap) *Parser {
	if s == nil {
		s = &scopeMap{make(map[string]Node)}
	}
	return &Parser{
		S: s,
	}
}

// AST is both the internal representation of the parsed expression,
// as well as the final result of parsing a line.
// AST implements the Node interface, and is typically the value which
// is saved to the scope or evaluated by the interpreter.
type AST struct {
	In    Node
	To    *AST
	From  *AST
	Name  string
	Token Token
	Type  EvalType
}

// addToAst adds a new node to the given parsers AST.
func addToAst(p *Parser, n Node) {
	p.tar.In = n
	p.tar.To = &AST{
		Token: p.tok,
		From:  p.tar,
	}
	p.tar = p.tar.To
}

// Eval evaluates the AST. ASTs are evaluated top down; the root is evaluated
// first and its result is passed to the next node. If the next node does not
// contain a Node, the result is returned.
func (a *AST) Eval(args Node) Set {
	if a.Type == Statement {
		a.Type = 0
		return nil
	}
	res := a.In.Eval(args)

	if a.To.In == nil {
		return res
	}

	return a.To.Eval(res)
}

// nextToken advances the parser to the next token and updates the
// last token and last literal.
func (p *Parser) nextToken() {
	p.lastTok = p.tok
	p.lastLit = p.lit
	p.tok, p.lit, _ = p.t.Next()
}

// Parse creates a new scanner from the argument and then parses it.
func (p *Parser) Parse(src string) Node {
	p.t = NewScanner(src)
	return p.parse()
}

// Run creates a new scanner from the argument and then parses and evaluates it.
func (p *Parser) Run(src string, arg Node) Set {
	p.t = NewScanner(src)
	return p.parse().Eval(arg)
}

// parse builds an AST from its scanner.
func (p *Parser) parse() Node {
	p.ast = &AST{}
	p.tar = p.ast

	for {
		p.nextToken()
		switch p.tok {
		case IDENT:
			switch p.t.Peek() {
			case LBRACE:
				p.parseLiteral()
			case LPAREN:
				p.parseFunctionMap()
			case STREAM:
				return p.parseStream()
			case APPLY:
				p.parseApply()
				p.ast.Type = Statement
			default:
				p.parseFunction()
			}
		case LPAREN:
			p.parseFunctionMap()
		case EOF:
			return p.ast
		}
	}
}

func (p *Parser) parseFunction() {
	if fn := p.S.Get(p.lit); fn != nil {
		addToAst(p, fn)
	}
}

func (p *Parser) parseApply() {
	p.nextToken()
	expr := &AST{}
	applied := &AST{To: p.ast, In: expr}

	err := p.S.Set(p.lastLit, applied)
	if err != nil {
		panic(err)
	}
	p.ast = applied
	p.tar = expr
}

func (p *Parser) parseStream() Node {
	fn := p.S.Get(p.lit)
	if fn == nil {
		panic("Stream to undefined function")
	}

	p.nextToken()
	expr := p.parse()
	newNode := &StreamStmnt{
		name:     p.lit,
		Consumer: fn,
		Expr:     expr,
	}
	return newNode
}

func (p *Parser) parseLiteral() {
	p.nextToken() // Consume the identifier token
	newNode := p.parseRecursiveLiteral()
	if newNode == nil {
		panic("null literal: literals can not be empty")
	}
	addToAst(p, newNode)
}

func (p *Parser) parseRecursiveLiteral() Node {
	var name string
	if p.lastTok == IDENT {
		name = p.lastLit
	}

	entries := make(map[string]Node)
	ns := NewMapNode(name, entries)
	for {
		// This either consumes the first "{", moves to next identifier,
		// or "}". If there is any other token, it is an error.
		p.nextToken()
		name := p.lit // Assuming the current token is the map name

		switch p.tok {
		case LBRACE:
			name = p.lastLit
			e := p.parseRecursiveLiteral()
			entries[name] = e
		case RBRACE:
			fallthrough
		case EOF:
			return ns
		default:
			if p.tok != IDENT {
				panic("Unexpected token in literal")
			}
			entries[name] = NewMapNode(name, make(map[string]Node))
		}
	}
}

// parseFunctionMap scans ahead until it finds a closing parenthesis or identifiers.
// Multiple identifiers are allowed, and each identifier is treated as the key in
// the resulting map. Multiple identifiers are allowed, but recursive fnmap literals
// are not. Using existing fnmap identifiers is possible.
func (p *Parser) parseFunctionMap() {
	name := "fnMap"
	if p.tok == IDENT {
		name = p.lit
		p.nextToken()
	}

	entries := make(map[string]Node)

	for {
		p.nextToken()
		switch p.tok {
		case IDENT:
			if p.t.Peek() == LBRACE {
				mapName := p.lit
				p.nextToken() // Consume the identifier token
				entries[mapName] = p.parseFunctionMapEntry()
			}
		case RPAREN:
			fallthrough
		case EOF:
			fnmap := NewFnMap(name, entries)
			// Scope.Set(name, NewFunctionMapNode(name, entries))

			addToAst(p, fnmap)
			return
		}
	}
}

func (p *Parser) parseFunctionMapEntry() Node {
	start := p.t.chPos
	p.nextToken() // Consume brace
	for {
		p.nextToken()
		switch p.tok {
		case RBRACE:
			return NewParser(p.S).Parse(p.t.src[start : p.t.chPos-1])
		case EOF:
			panic("Unexpected EOF in function map entry")
		}
	}
}
