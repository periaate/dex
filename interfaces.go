package dex

// Node describes a node in the expression tree. Nodes can be either
// expressions or sets.
type Node interface {
	Eval(args Node) Set
}

// Set describes an immutable datastructure.
// Todo: wrapper(?)
type Set interface {
	Node
	Get(key string) Set
	Name() string
}
