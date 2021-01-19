package terraform

import (
	"fmt"
	"testing"
)

func TestToJson(t *testing.T) {
	conf, err := ToJson("../../examples/.kakoi/provider.tf")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(conf)
}
