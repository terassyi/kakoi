package infra

import (
	"github.com/terassyi/kakoi/infra/aws"
	"github.com/terassyi/kakoi/infra/resource"
	"github.com/terassyi/kakoi/infra/state"
	"os"
	"path/filepath"
)

const (
	kakoi_dir string = ".kakoi"
	ext_yaml  string = ".yaml"
	ext_yml   string = ".yml"
	ext_json string = ".json"
)

type Initializer interface {
	Init() error
}

type initializer struct {
	workDir string
	conf *state.State
}

func NewInitializer(path string) (Initializer, error) {
	dir, file := filepath.Split(path)
	if err := state.ValidateExtName(file); err != nil {
		return nil, err
	}
	workDir, err := createWorkDir(dir)
	if err != nil {
		return nil, err
	}
	parser, err := state.NewParser(workDir, path)
	if err != nil {
		return nil, err
	}
	s, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &initializer{workDir: workDir, conf: s}, nil
}

func (i *initializer) init() error {
	// create state file
	if err := i.conf.CreateState(); err != nil {
		return err
	}
	// create storage
	storage := aws.NewS3(i.conf.Service.Name + kakoi_dir)
	builder, err := newBuilderFromResources(i.workDir, []resource.Resource{storage})
	if err != nil {
		return err
	}
	if err := builder.buildTemplate(); err != nil {
		return err
	}
	if err := builder.create(); err != nil {
		return err
	}
	return nil
}

func (i *initializer) Init() error {
	return i.init()
}

func createWorkDir(path string) (string, error) {
	// if work on current dir, path = ""
	workPath := filepath.Join(path, kakoi_dir)
	if err := os.MkdirAll(workPath, 0755); err != nil {
		return "", err
	}
	// pki cert files
	if err := os.MkdirAll(filepath.Join(workPath, "pki"), 0755); err != nil {
		return "", err
	}
	// server key pair
	if err := os.MkdirAll(filepath.Join(workPath, "keys"), 0755); err != nil {
		return "", err
	}
	// output files
	if err := os.MkdirAll(filepath.Join(workPath, "output"), 0755); err != nil {
		return "", err
	}
	// image files
	if err := os.MkdirAll(filepath.Join(workPath, "images"), 0755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(workPath, "storage"), 0755); err != nil {
		return "", err
	}
	return workPath, nil
}

