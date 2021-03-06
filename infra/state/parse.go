package state

import (
	yaml "github.com/goccy/go-yaml"
	"io/ioutil"
)

type Parser struct {
	workDir    string
	targetPath string
}

func NewParser(workDir, path string) (*Parser, error) {
	return &Parser{
		workDir:    workDir,
		targetPath: path,
	}, nil
}

func (p *Parser) Parse() (*State, error) {
	data, err := ioutil.ReadFile(p.targetPath)
	if err != nil {
		return nil, err
	}
	config := &State{
		WorkDir: p.workDir,
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
