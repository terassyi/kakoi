package log

import (
	"io/ioutil"
	"path/filepath"
)

const (
	log_ext = ".log"
)

func OutputLogs(workDir, name string, logs []string) error {
	fileName := name + log_ext
	path := filepath.Join(workDir, "log", fileName)
	var byteLogs []byte
	for _, l := range logs {
		byteLogs = append(byteLogs, []byte(l)...)
	}
	if err := ioutil.WriteFile(path, byteLogs, 0666); err != nil {
		return err
	}
	return nil
}
