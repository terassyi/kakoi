package infra

import (
	"fmt"
	"os"
	"testing"
)

func TestDestroyer_DestroyKakoiVpnFile(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	res, err := os.Stat("kakoi.ovpn")
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}

}
