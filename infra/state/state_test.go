package state

import (
	"encoding/json"
	"testing"
)

func TestSate_output(t *testing.T) {
	p, _ := NewParser("", "../../examples/example.yml")
	s, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	data, err := s.output()
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateState(); err != nil {
		t.Fatal(err)
	}
	c := &State{}
	if err := json.Unmarshal(data, c); err != nil {
		t.Fatal(err)
	}
	if c.Provider.Name != "aws" {
		t.Fatal("not match")
	}
	if c.Service.Name != "example" {
		t.Fatal("not match service.name")
	}
	if c.Service.Network.Range != "10.10.0.0/16" {
		t.Fatal("not match network.range")
	}
	if len(c.Service.Network.Subnets) != 2 {
		t.Fatal("not match number of subnet")
	}
	if c.Service.Network.Subnets[0].Name != "subnet1" {
		t.Fatal("not match network.subnet[0].name")
	}
	if len(c.Service.Hosts.Servers) != 2 {
		t.Fatalf("not match number of servers: actual %d", len(c.Service.Hosts.Servers))
	}
	if c.Service.Hosts.Servers[0].Name != "example-host1" {
		t.Fatal("not match servers[0].name")
	}
	if c.Service.Network.Vpn.Cert.Status != "" {
		t.Fatalf("invalid value: cert status %v", c.Service.Network.Vpn.Cert.Status)
	}
}
