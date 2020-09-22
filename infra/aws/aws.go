package aws

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/terassyi/kakoi/infra/terraform"
)

type Infra struct {
	workDir string
	profile string
	region  string
	tf      *tfexec.Terraform
}

func New(workDir, profile, region string) (*Infra, error) {
	tf, err := terraform.Prepare(workDir)
	if err != nil {
		return nil, err
	}
	return &Infra{
		workDir: workDir,
		profile: profile,
		region:  region,
		tf:      tf,
	}, nil
}

func (i *Infra) BuildTemplate() error {

	return nil
}
