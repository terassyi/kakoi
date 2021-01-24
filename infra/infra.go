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
	dir, file := filepath.Split(path)
	if err := state.ValidateExtName(file); err != nil {
		return nil, err
	}
	conf, workDir, err := parse(dir)
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
	// TODO
	if _, err := os.Stat("kakoi.ovpn"); err != nil {
		if err := i.createCertFiles(); err != nil {
			return err
		}
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

func isFileSpecified(path string) bool {
	if err := validateExtName(path); err != nil {
		return false
	}
	return true
}

func isExistStateFile(path string) bool {
	workDir, err := IsExistWorkDir(path)
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(workDir, kakoi_state)); err != nil {
		return false
	}
	return true
}

func parse(path string) (*state.State, string, error) {
	//if isFileSpecified(path) {
	//	workDir, err := createWorkDir(filepath.Base(path))
	//	if err != nil {
	//		return nil, "", err
	//	}
	//	parser, err := state.NewParser(workDir, path)
	//	if err != nil {
	//		return nil, "", err
	//	}
	//	s, err := parser.Parse()
	//	if err != nil {
	//		return nil, "", err
	//	}
	//	return s, workDir, nil
	//}
	if !isExistStateFile(filepath.Join(path, kakoi_dir)) {
		return nil, "", fmt.Errorf("kakoi.state is not found")
	}
	workDir := filepath.Join(path, kakoi_dir)
	parser, err := state.NewParser(workDir, filepath.Join(workDir, kakoi_state))
	if err != nil {
		return nil, "", err
	}
	s, err := parser.Parse()
	if err != nil {
		return nil, "", err
	}
	return s, workDir, nil
}