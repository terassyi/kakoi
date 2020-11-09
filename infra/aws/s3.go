package aws

import (
	"html/template"
	"os"
	"path/filepath"
)

type S3 struct {
	Name string
}

func NewS3(name string) *S3 {
	return &S3{Name: name}
}

func (s *S3) BuildTemplate(workDir string) error {
	fileName := "s3.tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("s3.tf.tmpl").ParseFiles("templates/aws/s3.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, s)
}
