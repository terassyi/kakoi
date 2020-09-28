package infra

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/terassyi/kakoi/config"
	"github.com/terassyi/kakoi/infra/resource"
	"github.com/terassyi/kakoi/infra/terraform"
)

type builder struct {
	tf        *tfexec.Terraform
	resources []resource.Resource
}

func newBuilder(workDir string, conf *config.Config) (*builder, error) {
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

func loadConfig(conf *config.Config) ([]resource.Resource, error) {
	return resource.New(conf)
}
