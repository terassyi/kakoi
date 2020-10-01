package resource

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"strings"
	"testing"
)

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
	v, err := newVpn("192.168.1.0/22", "test_domain", s)
	if err != nil {
		t.Fatal(err)
	}
	if err := v.BuildTemplate("../../data"); err != nil {
		t.Fatal(err)
	}
}

func TestGetVpnId(t *testing.T) {
	path := "../../examples/.kakoi/output/output.json"
	output := make(map[string]tfexec.OutputMeta)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 512)
	l, err := file.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf)
	fmt.Println(string(buf))
	if err := json.Unmarshal(buf[:l], &output); err != nil {
		t.Fatal(err)
	}
	vpn, ok := output["kakoi-vpn-id"]
	if !ok {
		t.Fatal("not found kakoi-vpn-id entry")
	}
	vpnId := string(vpn.Value)
	vpnId = strings.Replace(vpnId, "\"", "", -1)
	if vpnId[:4] != "cvpn" {
		t.Fatalf("actual: %v", vpnId)
	}
}