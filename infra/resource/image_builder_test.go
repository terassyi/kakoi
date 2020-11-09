package resource

import "testing"

func newImageBuilderTest() *ImageBuilder {
	return newImageBuilder("test", "test-region", nil, []string{"hoge1.sh", "hoge2.sh"})
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