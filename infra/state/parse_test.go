package state

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	p, _ := NewParser("", "../../examples/example.yml")
	c, err := p.Parse()
	if err != nil {
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
	if c.Service.Hosts.Servers[0].Size != "t2.micro" {
		t.Fatal("not match servers[0].size")
	}
	if c.Service.Hosts.Servers[1].Image.Id != "ami-0e7192738fc977648" {
		t.Fatal("not match server[1].name")
	}
	if c.Service.Hosts.Servers[1].Number != 2 {
		t.Fatal("not match server[1].number")
	}
}
