package dex

type FnWrapper struct {
	fn func(args Node) Set
}

// Eval calls the wrapped function with the given arguments.
func (fn *FnWrapper) Eval(args Node) Set { return fn.fn(args) }

// NewFnWrapper creates a node which uses the given function as the Eval method.
func NewFnWrapper(fn func(args Node) Set) Node { return &FnWrapper{fn} }

// MapNode is a node which contains a map of nodes. This is the data type
// used for map literals.
type MapNode struct {
	name    string
	entries map[string]Node
}

// Name returns the name of the map.
func (m *MapNode) Name() string { return m.name }

// Eval returns itself.
func (m *MapNode) Eval(args Node) Set { return m }

func (m *MapNode) Get(key string) Set {
	if v, ok := m.entries[key]; ok {
		return v.Eval(nil)
	}
	return nil
}

func NewMapNode(name string, entries map[string]Node) Node {
	return &MapNode{name, entries}
}

// StreamStmnt nodes are used to define a streams. Streams
// function similar to pipes, but are ran asynchronously.
type StreamStmnt struct {
	name     string
	Consumer Node
	Expr     Node
}

// Name returns the name of the stream.
func (s *StreamStmnt) Name() string { return s.name }

// Eval runs the stream with the evaluated expression.
func (s *StreamStmnt) Eval(args Node) Set { return s.Consumer.Eval(s.Expr) }

// FnMap is a map which implicitly calls expression found with the given
// sets name. In essence, this is a switch statement.
type FnMap struct {
	name    string
	entries map[string]Node
}

// Name returns the name of the function map.
func (fn *FnMap) Name() string { return fn.name }

// Eval maps the given nodes name to an expression and calls it with the node
// as argument.
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

// NewFnMap creates a new function map with the given name and entries.
func NewFnMap(name string, entries map[string]Node) Node { return &FnMap{name, entries} }
