package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTerraformTest(t *testing.T) {
	if err := terraformTest(); err != nil {
		t.Fatal(err)
	}
}

func TestTerraformApply(t *testing.T) {
	tf, err := Prepare("../../test")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	fmt.Println("initialize finished")
	if err := tf.Apply(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestTerraformDestroy(t *testing.T) {
	tf, err := Prepare("../../test")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	fmt.Println("initialize finished")
	if err := tf.Destroy(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestTerraformOutput(t *testing.T) {
	tf, err := Prepare("../../examples/.kakoi")
	if err != nil {
		t.Fatal(err)
	}
	outputDir := "../../examples/.kakoi/output"
	outputPath := filepath.Join(outputDir, "output.json")
	file, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	data, err := tf.Output(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(bytes)
	if _, err := file.Write(bytes); err != nil {
		fmt.Printf("here")
		t.Fatal(err)
	}
}