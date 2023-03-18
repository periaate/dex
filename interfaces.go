package dex

type Node interface {
	Eval(args Node) Set
}

type Set interface {
	Node
	Get(key string) Set
	Name() string
}
