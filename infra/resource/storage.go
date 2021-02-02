package resource

import (
	"html/template"
	"os"
	"path/filepath"
)

type Storage struct {
	Id string
}

func NewStorage(id string) *Storage {
	return &Storage{Id: id}
}

func (s *Storage) BuildTemplate(workDir string) error {
	fileName := "storage.tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("s3.tf.tmpl").ParseFiles("/etc/kakoi/templates/aws/s3.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, s)
}
