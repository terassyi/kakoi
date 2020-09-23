package resource

import "testing"

func TestNetwork_BuildTemplate(t *testing.T) {
	n, err := newNetwork("test", "10.10.0.0/16", "ap-northeast-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := n.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}

func TestSubnet_BuildTemplate(t *testing.T) {
	n, err := newNetwork("test", "10.10.0.0/16", "ap-northeast-1")
	if err != nil {
		t.Fatal(err)
	}
	s, err := newSubnet("test", "10.10.10.0/24", true, n)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}

func TestVpn_BuildTemplate(t *testing.T) {
	n, err := newNetwork("test", "10.10.0.0/16", "ap-northeast-1")
	if err != nil {
		t.Fatal(err)
	}
	s, err := newSubnet("test", "10.10.10.0/24", true, n)
	if err != nil {
		t.Fatal(err)
	}
	v, err := newVpn("192.168.1.0/22", s)
	if err != nil {
		t.Fatal(err)
	}
	if err := v.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}
