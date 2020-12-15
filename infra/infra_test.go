package infra

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildAbsPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	i := &initializer{
		workDir: "example/.kakoi",
		conf:    nil,
	}
	base := filepath.Join(wd, i.workDir[:len(i.workDir)-7])
	path, _ := i.buildAbsPath("./hoge")
	if path != filepath.Join(base, "./hoge") {
		t.Fatalf("actual: %v wanted: %v\n", path, filepath.Join(base, "./base"))
	}
	path, _ = i.buildAbsPath("/hoge")
	if path != "/hoge" {
		t.Fatalf("actual: %v wanted: /hoge\n", path)
	}
	i = &initializer{
		workDir: "/example/.kakoi",
		conf:    nil,
	}
	base = "/example"
	path, _ = i.buildAbsPath("./hoge")
	if path != filepath.Join(base, "./hoge") {
		t.Fatalf("actual: %v wanted: %v\n", path, filepath.Join(base, "./base"))
	}

}
