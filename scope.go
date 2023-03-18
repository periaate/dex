package dex

import "errors"

var (
	ErrMutation = errors.New("attempt to mutate a value")
)

type scopeMap struct {
	entries map[string]Node
}

func (m *scopeMap) Get(key string) Node {
	return m.entries[key]
}

func (m *scopeMap) Set(key string, n Node) error {
	if _, ok := m.entries[key]; ok {
		return ErrMutation
	}
	m.entries[key] = n
	return nil
}
