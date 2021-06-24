package cert

import "testing"

func TestGeneratePki(t *testing.T) {
	err := GeneratePki("test", "test.domain")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	if err := GenerateKeyPair("test", "../data"); err != nil {
		t.Fatal(err)
	}
}
