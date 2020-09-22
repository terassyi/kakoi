package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"io/ioutil"
	"os"
)

func terraformTest() error {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	tmp := tfinstall.LatestVersion(tmpDir, false)
	execPath, err := tfinstall.Find(context.Background(), tmp)
	if err != nil {
		return err
	}

	workingDir := "../data/"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		return err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(state.FormatVersion) // "0.1"
	return nil
}

func Prepare(workDir string) (*tfexec.Terraform, error) {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		return nil, err
	}
	//defer os.RemoveAll(tmpDir)
	tmp := tfinstall.LatestVersion(tmpDir, false)
	execPath, err := tfinstall.Find(context.Background(), tmp)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}
	tf, err := tfexec.NewTerraform(workDir, execPath)
	if err != nil {
		return nil, err
	}
	if err := tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.LockTimeout("60s")); err != nil {
		return nil, err
	}
	return tf, nil
}
