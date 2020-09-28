package aws

import (
	"os"
	"path/filepath"
	"text/template"
)

type Provider struct {
	Profile string
	Region  string
	Name    string
}

func NewProvider(profile, region string) *Provider {
	return &Provider{
		Profile: profile,
		Region:  region,
		Name:    "aws",
	}
}

func (p *Provider) BuildTemplate(workDir string) error {
	fileName := "provider.tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("provider.tf.tmpl").ParseFiles("templates/aws/provider.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, p)
}
