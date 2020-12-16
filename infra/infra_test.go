package infra

import (
	"bufio"
	"fmt"
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

func TestGetImageIdFromLog(t *testing.T) {
	i := &initializer{
		workDir: "../examples/.kakoi",
		conf:    nil,
	}

	file, err := os.Open("../examples/.kakoi/log/example-host1.log")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	id, err := i.getImageIdFromLog(lines)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
}
