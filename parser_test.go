package dex_test

import (
	"dex"
	"net/http"
	"strings"
	"testing"
)

type counter struct {
	count int
}

func getCountFn() (dex.Node, *counter) {
	count := new(counter)
	return dex.NewFnWrapper(func(args dex.Node) dex.Set {
		count.count++
		return nil
	}), count
}

func TestParser(t *testing.T) {
	c := dex.NewParser(nil).Parse("set{1 2 3{4 5}}}").Eval(nil)

	if c == nil {
		t.Errorf("ast does not contain set")
		return
	}
	if c.Name() != "set" {
		t.Errorf("expected set, got %v", c.Name())
	}
	if c.Get("1") == nil {
		t.Errorf("expected to find 1, got nil")
	}
	if c.Get("2") == nil {
		t.Errorf("expected to find 2, got nil")
	}
	if c.Get("3") == nil {
		t.Errorf("expected to find 3, got nil")
		return
	}
	if c.Get("3").Get("4") == nil {
		t.Errorf("expected to find 4, got nil")
	}
	if c.Get("3").Get("5") == nil {
		t.Errorf("expected to find 5, got nil")
	}
}

func TestFunctions(t *testing.T) {
	countFn, counter := getCountFn()
	p := dex.NewParser(nil)
	p.S.Set("count", countFn)
	p.Parse("set{1 2 3} count count count count").Eval(nil)

	if counter.count != 4 {
		t.Errorf("Expected 4 counts, got %v ", *counter)
	}
}

func TestApply(t *testing.T) {
	countFn, counter := getCountFn()
	p := dex.NewParser(nil)
	p.S.Set("count", countFn)

	p.Run("count4 < count count count count", nil)
	p.Run("set{1 2 3} count4", nil)
	// past.Eval(nil)

	if counter.count != 4 {
		t.Errorf("Expected 4 counts, got %v ", *counter)
	}
}

func TestStream(t *testing.T) {
	countFn, counter := getCountFn()
	httpfn := dex.NewFnWrapper(func(args dex.Node) dex.Set {
		http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			args.Eval(nil)
			if counter.count == 4 {
				return
			}
			w.Write([]byte("Hello, world!"))
		}))
		return nil
	})

	p := dex.NewParser(nil)
	p.S.Set("count", countFn)
	p.S.Set("serveHTTP", httpfn)

	// Automate this. Declarations are to not be evaluated.
	// They are evaluated when they are applied.
	p.Run("count4 < count count count count", nil)
	go p.Run("serveHTTP > count4", nil)
	http.Get("http://localhost:8080")

	if counter.count != 4 {
		t.Errorf("Expected 4 counts, got %v ", *counter)
	}
}

func TestFunctionMap(t *testing.T) {
	countFn, counter := getCountFn()
	p := dex.NewParser(nil)

	p.S.Set("count", countFn)
	p.Run("count2 < count count", nil)

	if p.S.Get("count2") == nil || p.S.Get("count") == nil {
		t.Errorf("Expected to find count2 and count")
	}

	doOnce := dex.NewMapNode("doOnce", map[string]dex.Node{})
	doThrice := dex.NewMapNode("doThrice", map[string]dex.Node{})

	dx := "set(doOnce{count} doThrice{count2 count})"
	p.Run(dx, doOnce)

	if counter.count != 1 {
		t.Errorf("Expected 1 count, got %v ", *counter)
	}

	p.Run(dx, doThrice)
	if counter.count != 4 {
		t.Errorf("Expected 4 counts, got %v ", *counter)
	}
}

func TestMappedStream(t *testing.T) {
	countFn, counter := getCountFn()
	var call int
	httpfn := dex.NewFnWrapper(func(args dex.Node) dex.Set {
		http.ListenAndServe("localhost:8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			pathArray := strings.Split(r.URL.Path[1:], "/")

			ar := dex.NewMapNode(pathArray[0], map[string]dex.Node{})
			args.Eval(ar)
			call++
			w.Write([]byte("Hello, world!"))
		}))
		return nil
	})

	p := dex.NewParser(nil)

	p.S.Set("count", countFn)
	p.S.Set("serveHTTP", httpfn)

	p.Run("count2 < count count", nil)
	p.Run("httpMap < (doOnce{count} doTwice{count count})", nil)

	if p.S.Get("count2") == nil || p.S.Get("count") == nil {
		t.Errorf("Expected to find count2 and count")
	}

	go p.Run("serveHTTP > httpMap", nil)
	r, err := http.Get("http://localhost:8081/doOnce")
	if counter.count != 1 {
		t.Errorf("Expected 1 count, got %v ", *counter)
	}
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if r.StatusCode != 200 {
		t.Errorf("Expected 200, got %v", r.StatusCode)
	}

	r, err = http.Get("http://localhost:8081/doTwice")
	if counter.count != 3 {
		t.Errorf("Expected 3 counts, got %v ", *counter)
	}
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if r.StatusCode != 200 {
		t.Errorf("Expected 200, got %v", r.StatusCode)
	}
}
