package cert

import "testing"

func TestGeneratePki(t *testing.T) {
	err := GeneratePki()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	if err := GenerateKeyPair("test", "../data"); err != nil {
		t.Fatal(err)
	}
}
