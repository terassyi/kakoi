package infra

import (
	"fmt"
	"github.com/terassyi/kakoi/config"
	"github.com/terassyi/kakoi/infra/resource"
	"path/filepath"
)

type Infrastructure interface {
	WorkDir() string
	Provider() string
	Build() error
	Create() error
}

type infra struct {
	*builder
	workDir  string
	provider string
}

func New(path string) (Infrastructure, error) {
	fmt.Println("path: ", path)
	dir, file := filepath.Split(path)
	if err := config.ValidateExtName(file); err != nil {
		return nil, err
	}
	workDir, err := config.CreateWorkDir(dir)
	if err != nil {
		return nil, err
	}
	parser, err := config.NewParser(workDir, path)
	if err != nil {
		return nil, err
	}
	conf, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	builder, err := newBuilder(workDir, conf)
	if err != nil {
		return nil, err
	}
	return &infra{
		builder:  builder,
		workDir:  workDir,
		provider: conf.Provider.Name,
	}, nil
}

func (i *infra) WorkDir() string {
	return i.workDir
}

func (i *infra) Provider() string {
	return i.provider
}

func (i *infra) Build() error {
	// build template file
	for _, b := range i.resources {
		if err := b.BuildTemplate(i.workDir); err != nil {
			return err
		}
	}
	return nil
}

func (i *infra) Create() error {
	if err := i.Build(); err != nil {
		return err
	}
	if err := i.createCertFiles(); err != nil {
		return err
	}
	//return i.tf.Apply(context.Background())
	return nil
}

func (i *infra) createCertFiles() error {
	for _, r := range i.resources {
		switch c := r.(type) {
		case *resource.Pki:
			if err := c.Create(); err != nil {
				return err
			}
		case *resource.KeyPair:
			if err := c.Create(); err != nil {
				return err
			}
		}
	}
	return nil
}
