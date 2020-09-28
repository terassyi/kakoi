package infra

import (
	"context"
	"github.com/terassyi/kakoi/config"
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
	dir, file := filepath.Split(path)
	if err := config.ValidateExtName(file); err != nil {
		return nil, err
	}
	workDir, err := config.CreateWorkDir(dir)
	if err != nil {
		return nil, err
	}
	parser, err := config.NewParser(path, workDir)
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
	return i.tf.Apply(context.Background())
}
