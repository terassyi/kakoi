package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/terassyi/kakoi/infra/resource"
	"github.com/terassyi/kakoi/infra/state"
	"os"
	"path/filepath"
)

type Infrastructure interface {
	Name() string
	WorkDir() string
	Provider() string
	//Build() error
	Create() error
}

type infra struct {
	*builder
	name string
	workDir  string
	provider string
}

func New(path string) (Infrastructure, error) {
	fmt.Println("path: ", path)
	dir, file := filepath.Split(path)
	if err := validateExtName(file); err != nil {
		return nil, err
	}
	workDir, err := state.CreateWorkDir(dir)
	if err != nil {
		return nil, err
	}
	parser, err := state.NewParser(workDir, path)
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
		name: conf.Service.Name,
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

func (i *infra) Name() string {
	return i.name
}

func (i *infra) Create() error {
	if err := i.buildTemplate(); err != nil {
		return err
	}
	if err := i.createCertFiles(); err != nil {
		return err
	}
	ctx :=  context.Background()
	err := i.tf.Init(ctx, tfexec.Upgrade(true), tfexec.LockTimeout("60s"))
	if err != nil {
		return err
	}
	if err := i.tf.Apply(ctx); err != nil {
		return err
	}

	if err := i.Output(filepath.Join(i.workDir, "output")); err != nil {
		return err
	}

	vpn := i.findVpn()
	if vpn == nil {
		return fmt.Errorf("vpn resource is not found")
	}
	if err := vpn.Create(); err != nil {
		return err
	}
	return nil
}

func (i *infra) findVpn() *resource.Vpn {
	for _, r := range i.resources {
		switch r := r.(type) {
		case *resource.Vpn:
			return r
		}
	}
	return nil
}

func (i *infra) Output(outputDir string) error {
	outputPath := filepath.Join(outputDir, "output.json")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := i.output()
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if _, err := file.Write(bytes); err != nil {
		return err
	}
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

func (i *infra) Destroy() error {
	return i.destroy()
}

func IsExistWorkDir(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	base := filepath.Dir(path)
	workDir := filepath.Join(base, ".kakoi")
	if _, err := os.Stat(workDir); err != nil {
		return "", err
	}
	return workDir, nil
}

func validateExtName(file string) error {
	extName := filepath.Ext(file)
	if extName != ext_yaml && extName != ext_yml && extName != ext_json {
		return fmt.Errorf("config file must be .yaml of .yml format: %s", file)
	}
	return nil
}