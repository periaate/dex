package lan

type Node interface {
	Eval(args ...Node) Set
}

type Set interface {
	Node
	Get(key string) Set
}

type Parser struct {
	tk        *Tokenizer
	tok       Token
	lastToken Token
	ast       *AST
	tar       *AST
}

type AST struct {
	In Node
	To *AST
}

func (a *AST) Eval(args ...Node) Set {
	res := a.In.Eval(args...)

	if a.To.In == nil {
		return res
	}

	return a.To.Eval(res)
}

func NewParser(src string) *Parser {
	return &Parser{
		tk: NewTokenizer(src),
	}
}

var Scope = &scopeMap{make(map[string]Node)}

func (p *Parser) nextToken() {
	p.lastToken = p.tok
	p.tok = p.tk.NextToken()
}

func (p *Parser) Parse() Node {
	p.ast = &AST{}
	p.tar = p.ast
	// defer recoverParse()

	for {
		p.nextToken()
		switch p.tok.Type {
		case IDENT:
			switch p.tk.PeekToken().Type {
			case LBRACE:
				p.parseLiteral()
			case STREAM:
				return p.parseStream()
			case APPLY:
				p.parseApply()
			default:
				p.parseFunction()
			}
		case EOF:
			return p.ast
		}
	}
}

func (p *Parser) parseFunction() {
	if fn := Scope.Get(p.tok.Value); fn != nil {
		nn := fn
		p.tar.In = nn
		p.tar.To = &AST{}
		p.tar = p.tar.To
	}
}

func (p *Parser) parseApply() {
	name := p.tok.Value

	p.nextToken()
	expr := &AST{}
	nast := &AST{To: p.ast, In: expr}

	err := Scope.Set(name, nast)
	if err != nil {
		panic(err)
	}
	p.ast = nast
	p.tar = expr
}

func (p *Parser) parseStream() Node {
	fn := Scope.Get(p.tok.Value)
	if fn == nil {
		panic("Stream to undefined function")
	}

	p.nextToken()
	expr := p.Parse()
	newNode := &StreamNode{
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
	p.tar.In = newNode
	p.tar.To = &AST{}
	p.tar = p.tar.To
}

func (p *Parser) parseRecursiveLiteral() Node {
	var name string
	if p.lastToken.Type == IDENT {
		name = p.lastToken.Value
	}

	entries := make(map[string]Node)
	ns := NewMapNode(name, entries)
	for {
		// This either consumes the first "{", moves to next identifier,
		// or "}". If there is any other token, it is an error.
		p.nextToken()
		name := p.tok.Value // Assuming the current token is the map name

		switch p.tok.Type {
		case LBRACE:
			// Like at the start of the function, we need to use the
			// last token, as the current one isn't an identifier.
			name = p.lastToken.Value
			e := p.parseRecursiveLiteral()
			entries[name] = e
		case RBRACE:
			fallthrough
		case EOF:
			return ns
		default:
			if p.tok.Type != IDENT {
				panic("Unexpected token in literal")
			}
			entries[name] = NewMapNode(name, make(map[string]Node))
		}
	}
}
