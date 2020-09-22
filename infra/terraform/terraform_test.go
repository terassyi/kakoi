package terraform

import (
	"context"
	"testing"
)

const inline_terraform_template string = `
resource "null_resource" "example1" {
  provisioner "local-exec" {
    command = "echo hogehoge > ../data/test.txt"
  }
}
`

func TestTerraformTest(t *testing.T) {
	if err := terraformTest(); err != nil {
		t.Fatal(err)
	}
}

func TestTerraformApply(t *testing.T) {
	tf, err := Prepare("../data")
	if err != nil {
		t.Fatal(err)
	}
	if err := tf.Apply(context.Background()); err != nil {
		t.Fatal(err)
	}
}
