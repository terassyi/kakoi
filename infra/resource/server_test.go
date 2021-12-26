package resource

import "testing"

func TestServer_BuildTemplate(t *testing.T) {
	n, err := newNetwork("test", "10.10.0.0/16", "ap-northeast-1")
	if err != nil {
		t.Fatal(err)
	}
	s, err := newSubnet("test", "10.10.10.0/24", true, n)
	if err != nil {
		t.Fatal(err)
	}
	k := newKeyPair("./test", "test-key")

	server, err := newServer("test", "test-size-type", 0, s, k, []int{80})
	if err != nil {
		t.Fatal(err)
	}
	server.SetImage(NewImage("test-image-id", "test-image-path"))
	if err := server.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}
