package resource

import (
	"html/template"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	Name string
	Id string
	Subnet *Subnet
	Ports []int
	KeyPair *KeyPair
	PrivateIp *net.IP
	ImageBuilder ImageBuilder
}

func newServer(name string, subnet *Subnet, key *KeyPair, ports []int) (*Server, error) {
	return &Server{
		Name:      name,
		Subnet:    subnet,
		Ports:     ports,
		KeyPair:   key,
		PrivateIp: nil,
	}, nil
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