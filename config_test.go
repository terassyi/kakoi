package main

import "testing"

func TestOutputOvpnConfig(t *testing.T) {
	err := outputOvpnConfig("test", "fuk", "192.168.1.0/22")
	if err != nil {
		t.Fatal(err)
	}
}
