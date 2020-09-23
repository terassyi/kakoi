package resource

import "testing"

func TestPki_BuildTemplate(t *testing.T) {
	p := newPki("path/to/key", "test.kakoi.terassyi.net")
	if err := p.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}

func TestKeyPair_BuildTemplate(t *testing.T) {
	k := newKeyPair("path/to/key", "test")
	if err := k.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}
