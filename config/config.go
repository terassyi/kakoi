package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	kakoi_dir string = ".kakoi"
	ext_yaml  string = ".yaml"
	ext_yml   string = ".yml"
)

func CreateWorkDir(path string) (string, error) {
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
	return workPath, nil
}

func ValidateExtName(file string) error {
	extName := filepath.Ext(file)
	if extName != ext_yaml && extName != ext_yml {
		return fmt.Errorf("config file must be .yaml of .yml format: %s", file)
	}
	return nil
}
