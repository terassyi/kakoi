package infra

import (
	"context"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/terassyi/kakoi/infra/resource"
	"github.com/terassyi/kakoi/infra/state"
	"github.com/terassyi/kakoi/infra/terraform"
)

type builder struct {
	tf        *tfexec.Terraform
	resources []resource.Resource
}

func newBuilder(workDir string, conf *state.State) (*builder, error) {
	tf, err := terraform.Prepare(workDir)
	if err != nil {
		return nil, err
	}
	resources, err := loadConfig(conf)
	if err != nil {
		return nil, err
	}
	return &builder{
		tf:        tf,
		resources: resources,
	}, nil
}

func newBuilderFromResources(workDir string, resources []resource.Resource) (*builder, error) {
	tf, err := terraform.Prepare(workDir)
	if err != nil {
		return nil, err
	}
	return &builder{
		tf:        tf,
		resources: resources,
	}, nil
}

func loadConfig(conf *state.State) ([]resource.Resource, error) {
	return resource.New(conf)
}

func (b *builder) buildTemplate() error {
	for _, r := range b.resources {
		if err := r.BuildTemplate(b.tf.WorkingDir()); err != nil {
			return err
		}
	}
	return nil
}

func (b *builder) create() error {
	if err := b.tf.Init(context.Background()); err != nil {
		return err
	}
	if err := b.tf.Apply(context.Background()); err != nil {
		return err
	}
	return nil
}

func (b *builder) output() (map[string]string, error) {
	data, err := b.tf.Output(context.Background())
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for k, v := range data {
		res[k] = string(v.Value)
	}
	return res, nil
}

func (b *builder) destroy() error {
	return b.tf.Destroy(context.Background())
}
