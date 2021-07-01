package resource

import (
	"os"
	"path/filepath"
	"testing"
)

func newImageBuilderTest() *ImageBuilder {
	ib, _ := NewImageBuilder("test2", "test-region", "~/example/.kakoi", "test-custom-base-image", "test-user", "test-owner", nil, []string{"hoge1.sh", "hoge2.sh"})
	return ib
}

func TestNewImageBuilder(t *testing.T) {
	ib := newImageBuilderTest()
	if ib.Name != "test2" {
		t.Fatalf("not match image builder name")
	}
	if ib.Region != "test-region" {
		t.Fatalf("not match image builder region")
	}
	wd, _ := os.Getwd()
	if ib.BuildSpecPath != filepath.Join(wd, ".kakoi", "images") {
		t.Fatalf("not match work path")
	}
}

func TestPackerBuilder_outputJson(t *testing.T) {
	i := newImageBuilderTest()
	p, err := i.createPackerBuilder()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.outputJson("../../test/output/packer_image_builder.test.json"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateBuildSpec(t *testing.T) {
	if err := createBuildSpec("../../test/output/buildspec.yml", "test"); err != nil {
		t.Fatal(err)
	}
}

func TestImageBuilder_BuildTemplate(t *testing.T) {
	ib := newImageBuilderTest()
	if err := ib.BuildTemplate("../../examples/.kakoi"); err != nil {
		t.Fatal(err)
	}
}
