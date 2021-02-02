package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type S3 struct {
	Name string
}

func NewS3(name string) *S3 {
	return &S3{Name: name}
}

func (s *S3) BucketName() string {
	return fmt.Sprintf("kakoi.%s", s.Name)
}

func (s *S3) BuildTemplate(workDir string) error {
	fileName := "s3.tf"
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

func (s *S3) UploadFile(profile, src, dst string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
		SharedConfigState: session.SharedConfigEnable,
	}))
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	uploader := s3manager.NewUploader(sess)
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:                  aws.String(s.Name),
		Key:                     aws.String(dst),
		Body:                    file,
	}); err != nil {
		return err
	}
	return nil
}

type S3Uploader struct {
	Src string
	Key string
	Name string
}

func NewS3Uploader(src, key string) *S3Uploader {
	return &S3Uploader{
		Src: src,
		Key: key,
		Name: strings.Replace(filepath.Base(key), ".", "-", -1),
	}
}

func (su *S3Uploader) BuildTemplate(workDir string) error {
	fileName := fmt.Sprintf("s3_upload-%s.tf", su.Name)
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("s3_upload.tf.tmpl").ParseFiles("/etc/kakoi/templates/aws/s3_upload.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, su)
}