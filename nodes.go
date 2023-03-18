package dex

type FnWrapper struct {
	fn func(args Node) Set
}

func (fn *FnWrapper) Eval(args Node) Set {
	return fn.fn(args)
}
func NewFnWrapper(fn func(args Node) Set) Node {
	return &FnWrapper{fn}
}

type MapNode struct {
	name    string
	entries map[string]Node
}

func (m *MapNode) Name() string {
	return m.name
}

func (m *MapNode) Eval(args Node) Set {
	return m
}

func (m *MapNode) Get(key string) Set {
	if v, ok := m.entries[key]; ok {
		return v.Eval(nil)
	}
	return nil
}

func NewMapNode(name string, entries map[string]Node) Node {
	return &MapNode{name, entries}
}

type StreamStmnt struct {
	name     string
	Consumer Node
	Expr     Node
}

func (s *StreamStmnt) Name() string {
	return s.name
}

func (s *StreamStmnt) Eval(args Node) Set {
	return s.Consumer.Eval(s.Expr)
}

type FnMap struct {
	name    string
	entries map[string]Node
}

func (fn *FnMap) Name() string {
	return fn.name
}

func (fn *FnMap) Eval(args Node) Set {
	if args == nil {
		panic("Function map called without arguments")
	}
	s := args.Eval(nil)
	if v, ok := fn.entries[s.Name()]; ok {
		return v.Eval(args)
	}
	return nil
}

func NewFnMap(name string, entries map[string]Node) Node {
	return &FnMap{name, entries}
}
