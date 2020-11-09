package resource

import (
	"fmt"
	"strings"
)

func getOutputValue(data map[string]string, key string) (string, error) {
	v, ok := data[key]
	if !ok {
		return "", fmt.Errorf("no such key: %s", key)
	}
	return strings.Replace(v, "\"", "", -1), nil

}
