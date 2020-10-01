package resource

import (
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"strings"
)

func getOutputValue(data map[string]tfexec.OutputMeta, key string) (string, error) {
	v, ok := data[key]
	if !ok {
		return "", fmt.Errorf("no such key: %s", key)
	}
	want := string(v.Value)
	return strings.Replace(want, "\"", "", -1), nil

}
