package lan_test

import (
	"lan"
	"net/http"
	"testing"
)

var counter int
var countFn = lan.NewFn(func(args ...lan.Node) lan.Set {
	counter++
	if len(args) == 0 {
		return nil
	}
	if args[0] == nil {
		return nil
	}
	return args[0].Eval(args...)
})

func TestParser(t *testing.T) {
	reset()
	np := lan.NewParser("set{1 2 3{4 5}}}")

	ast := np.Parse()
	if ast == nil {
		t.Errorf("ast is nil")
		return
	}
	c := ast.Eval()
	if c == nil {
		t.Errorf("ast does not contain set")
		return
	}
	if c.Get("1") == nil {
		t.Errorf("expected to find 1, got nil")
		return
	}
	if c.Get("2") == nil {
		t.Errorf("expected to find 2, got nil")
		return
	}
	if c.Get("3") == nil {
		t.Errorf("expected to find 3, got nil")
		return
	}
	if c.Get("3").Get("4") == nil {
		t.Errorf("expected to find 4, got nil")
		return
	}
	if c.Get("3").Get("5") == nil {
		t.Errorf("expected to find 5, got nil")
		return
	}
}

func TestFunctions(t *testing.T) {
	reset()
	np := lan.NewParser("set{1 2 3} count count count count")

	lan.Scope.Set("count", countFn)

	ast := np.Parse()

	ast.Eval()

	if counter != 4 {
		t.Errorf("Expected 4 counts, got %v", counter)
	}
}

func TestApply(t *testing.T) {
	reset()
	np := lan.NewParser("count4 < count count count count")
	npp := lan.NewParser("set{1 2 3} count4")

	lan.Scope.Set("count", countFn)

	np.Parse()
	past := npp.Parse()
	past.Eval()

	if counter != 4 {
		t.Errorf("Expected 4 counts, got %v", counter)
	}
}

func TestStream(t *testing.T) {
	reset()

	httpfn := lan.NewFn(func(args ...lan.Node) lan.Set {
		http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			args[0].Eval()
			if counter == 4 {
				return
			}
			w.Write([]byte("Hello, world!"))
		}))
		return nil
	})
	np := lan.NewParser("count4 < count count count count")
	npp := lan.NewParser("serveHTTP > count4")

	lan.Scope.Set("count", countFn)
	lan.Scope.Set("serveHTTP", httpfn)

	// Automate this. Declarations are to not be evaluated.
	// They are evaluated when they are applied.
	np.Parse()
	past := npp.Parse()
	go past.Eval()
	http.Get("http://localhost:8080")

	if counter != 4 {
		t.Errorf("Expected 4 counts, got %v", counter)
	}
}

func reset() {
	counter = 0
	lan.ResetScope()
}
