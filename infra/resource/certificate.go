package resource

import (
	"github.com/terassyi/kakoi/cert"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type Pki struct {
	Path       string
	WorkPath string
	domain     string
	CaCert     string
	CaKey      string
	ServerCert string
	ServerKey  string
	ClientCert string
	ClientKey  string
}

func newPki(path, domain string) *Pki {
	return &Pki{
		Path:       path,
		WorkPath: "pki",
		domain:     domain,
		CaCert:     "ca." + domain + ".crt",
		CaKey:      "ca." + domain + ".key",
		ServerCert: "server." + domain + ".crt",
		ServerKey:  "server." + domain + ".key",
		ClientCert: "client." + domain + ".crt",
		ClientKey:  "client." + domain + ".key",
	}
}

func (p *Pki) BuildTemplate(workDir string) error {
	pkiFile, err := os.Create(filepath.Join(workDir, p.domain+".tf"))
	if err != nil {
		return err
	}
	defer pkiFile.Close()
	t, err := template.New("pki.tf.tmpl").ParseFiles("/etc/kakoi/templates/aws/pki.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(pkiFile, p)
}

func (p *Pki) Create() error {
	return cert.GeneratePki(p.Path, p.domain)
}

func (p *Pki) ReadCaCert() (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(p.Path, p.CaCert))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *Pki) ReadClientCert() (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(p.Path, p.ClientCert))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *Pki) ReadClientKey() (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(p.Path, p.ClientKey))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type KeyPair struct {
	Path string
	WorkPath string
	Name string
	Pem  string
	Pub  string
}

func newKeyPair(path, name string) *KeyPair {
	return &KeyPair{
		Path: path,
		WorkPath: "keys",
		Name: name,
		Pem:  name + ".pem",
		Pub:  name + ".pub",
	}
}

func (k *KeyPair) BuildTemplate(workDir string) error {
	keyPairFile, err := os.Create(filepath.Join(workDir, k.Name + ".tf"))
	if err != nil {
		return err
	}
	defer keyPairFile.Close()
	t, err := template.New("keypair.tf.tmpl").ParseFiles("/etc/kakoi/templates/aws/keypair.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(keyPairFile, k)
}

func (k *KeyPair) Create() error {
	if err := cert.GenerateKeyPair(k.Name, k.Path); err != nil {
		return err
	}
	return nil
}