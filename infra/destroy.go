package infra

import (
	"context"
	"os"

	//"context"
	"fmt"
	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/terassyi/kakoi/infra/aws"
	"github.com/terassyi/kakoi/infra/state"
	"github.com/terassyi/kakoi/infra/terraform"
	"path/filepath"
)

type Destroyer interface {
	Destroy() error
}

type destroyer struct {
	workDir string
	conf    *state.State
}

func NewDestroyer(path string) (Destroyer, error) {
	if !isExistStateFile(path) {
		return nil, fmt.Errorf("[ERROR] state file is not found")
	}
	base := filepath.Join(filepath.Dir(path), kakoi_dir)
	parser, err := state.NewParser(base, filepath.Join(base, kakoi_state))
	if err != nil {
		return nil, err
	}
	s, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &destroyer{
		workDir: base,
		conf:    s,
	}, nil
}

func (d *destroyer) Destroy() error {
	tf, err := terraform.Prepare(d.workDir)
	if err != nil {
		return err
	}
	fmt.Println(tf)
	//destroy terraform resource
	if err := tf.Destroy(context.Background()); err != nil {
		return err
	}
	// destroy image
	if err := d.destroyImage(); err != nil {
		return err
	}
	if err := d.destroyWorkDir(); err != nil {
		return err
	}
	return nil
}

func (d *destroyer) destroyImage() error {
	var ids []*string
	for _, s := range d.conf.Service.Hosts.Servers {
		if s.Image.Custom && s.Image.Id != "" {
			ids = append(ids, awsSdk.String(s.Image.Id))
		}
	}
	if err := aws.DeleteImage(d.conf.Provider.Profile, ids); err != nil {
		return err
	}
	return nil
}

func (d *destroyer) destroyWorkDir() error {
	return os.Remove(d.workDir)
}
