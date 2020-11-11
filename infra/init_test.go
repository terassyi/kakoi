package infra

import (
	"path/filepath"
	"testing"
)

func TestConvertExtFormat(t *testing.T) {
	path := "~/test_dir/test.ova"
	ext := filepath.Ext(path)

	res := convertExtImageFormat(ext)
	if res != "OVA" {
		t.Fatalf("result = %s(%s)", res, ext)
	}
}
