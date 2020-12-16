package aws

import (
	"fmt"
	"testing"
)

func TestGetLog(t *testing.T) {
	s, err := GetLog("../../examples/.kakoi", "kakoi", "kakoi-example-host1", "build", "kakoi-example-host1-build/cc8f8525-0a26-496b-8a00-b98f5f1b59ba")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range s {
		fmt.Println(v)
	}
}
