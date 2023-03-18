package dex

type EvalType int

const (
	_          EvalType = iota
	Expression          // Evaluate
	Statement           // Do not evaluate
)

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

func NewParser(s *scopeMap) *Parser {
	if s == nil {
		s = &scopeMap{make(map[string]Node)}
	}
	return &Parser{
		S: s,
	}
}

type AST struct {
	In    Node
	To    *AST
	From  *AST
	Name  string
	Token Token
	Type  EvalType
}

func help(p *Parser, n Node) {
	p.tar.In = n
	p.tar.To = &AST{
		Token: p.tok,
		From:  p.tar,
	}
	p.tar = p.tar.To
}

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

func (p *Parser) nextToken() {
	p.lastTok = p.tok
	p.lastLit = p.lit
	p.tok, p.lit, _ = p.t.Next()
}

func (p *Parser) Parse(src string) Node {
	p.t = NewScanner(src)
	return p.parse()
}
func (p *Parser) Run(src string, arg Node) Set {
	p.t = NewScanner(src)
	return p.parse().Eval(arg)
}

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
		help(p, fn)
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
	help(p, newNode)
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

			help(p, fnmap)
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
