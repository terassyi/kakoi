package resource

import (
	"html/template"
	"net"
	"os"
	"path/filepath"
)

const (
	aws_default_image_id = "ami-01748a72bed07727c"
	aws_default_size = "t2.micro"
)

type Server struct {
	Name string
	Image *Image
	Size string
	Subnet *Subnet
	Ports []int
	KeyPair *KeyPair
	PrivateIp *net.IP
	ImageBuilder ImageBuilder
}

func newServer(name, size string, subnet *Subnet, key *KeyPair, ports []int) (*Server, error) {
	s := &Server{
		Name:      name,
		Size: size,
		Subnet:    subnet,
		Ports:     ports,
		KeyPair:   key,
		PrivateIp: nil,
	}
	if size == "" {
		s.Size = aws_default_size
	} else {
		s.Size = size
	}
	return s, nil
}
func (s *Server) SetImage(image *Image) {
	s.Image = image
}

func (s *Server) SetImageBuilder(imageBuilder *ImageBuilder) {
	s.ImageBuilder = *imageBuilder
}

func (s *Server) BuildTemplate(workDir string) error {
	fileName := "kakoi-" + s.Name + ".tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("server.tf.tmpl").ParseFiles("templates/aws/server.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, s)
}

type Image struct {
	Id string
	Path string
}

func NewImage(id, path string) *Image {
	image := &Image{}
	if id == "" {
		image.Id = aws_default_image_id
	} else {
		image.Id = id
	}
	image.Path = path
	return image
}