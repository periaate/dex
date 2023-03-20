package dex

import "errors"

var (
	ErrMutation = errors.New("attempt to mutate a value")
)

// scopeMap is used as the scope for the parser.
// Variables in scope are immutable.
type scopeMap struct {
	entries map[string]Node
}

// Get returns the value for the given key.
func (m *scopeMap) Get(key string) Node { return m.entries[key] }

// Set sets the value for the given key if it doesn't exist.
// Values are immutable, trying to mutate a value will return an error.
func (m *scopeMap) Set(key string, n Node) error {
	if _, ok := m.entries[key]; ok {
		return ErrMutation
	}
	m.entries[key] = n
	return nil
}
