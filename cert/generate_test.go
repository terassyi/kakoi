package cert

import "testing"

func TestGeneratePki(t *testing.T) {
	err := GeneratePki("../data", "test.kakoi.terassyi.net")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	if err := GenerateKeyPair("test", "../data"); err != nil {
		t.Fatal(err)
	}
}
