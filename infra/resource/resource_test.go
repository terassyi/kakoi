package resource

import (
	"github.com/terassyi/kakoi/infra/state"
	"testing"
)

func TestNew(t *testing.T) {
	p, _ := state.NewParser("", "../../examples/example.yml")
	c, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	re, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, r := range re {
		switch r.(type) {
		case *Server:
			count += 1
		}
	}
	if count != 3 {
		t.Fatal("not  match number of servers: ", count)
	}
}
