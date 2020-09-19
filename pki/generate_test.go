package pki

import "testing"

func TestGeneratePki(t *testing.T) {
	err := GeneratePki()
	if err != nil {
		t.Fatal(err)
	}
}
