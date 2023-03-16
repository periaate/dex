package lan

type FunctionNode struct {
	fn func(args ...Node) Set
}

func (fn *FunctionNode) Eval(args ...Node) Set {
	return fn.fn(args...)
}
func NewFn(fn func(args ...Node) Set) Node {
	return &FunctionNode{fn}
}

type MapNode struct {
	name    string
	entries map[string]Node
}

func (m *MapNode) Eval(args ...Node) Set {
	return m
}

func (m *MapNode) Get(key string) Set {
	if v, ok := m.entries[key]; ok {
		return v.Eval()
	}
	return nil
}

func NewMapNode(name string, entries map[string]Node) Node {
	return &MapNode{name, entries}
}

type StreamNode struct {
	Consumer Node
	Expr     Node
}

func (s *StreamNode) Eval(args ...Node) Set {
	return s.Consumer.Eval(s.Expr)
}
