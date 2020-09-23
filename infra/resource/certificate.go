package resource

import (
	"os"
	"path/filepath"
	"text/template"
)

type Pki struct {
	Path       string
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
		domain:     domain,
		CaCert:     "ca." + domain + ".crt",
		CaKey:      "ca." + domain + ".key",
		ServerCert: "server." + domain + ".crt",
		ServerKey:  "server." + domain + ".key",
		ClientCert: "client." + domain + ".crt",
		ClientKey:  "client" + domain + ".key",
	}
}

func (p *Pki) BuildTemplate(workDir string) error {
	pkiFile, err := os.Create(filepath.Join(workDir, p.domain+".tf"))
	if err != nil {
		return err
	}
	defer pkiFile.Close()
	t, err := template.New("pki.tf.tmpl").ParseFiles("../../templates/aws/pki.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(pkiFile, p)
}

type KeyPair struct {
	Path string
	Name string
	Pem  string
	Pub  string
}

func newKeyPair(path, name string) *KeyPair {
	return &KeyPair{
		Path: path,
		Name: name,
		Pem:  name + ".pem",
		Pub:  name + ".pub",
	}
}

func (k *KeyPair) BuildTemplate(workDir string) error {
	keyPairFile, err := os.Create(filepath.Join(workDir, k.Name+"-keypair.tf"))
	if err != nil {
		return err
	}
	defer keyPairFile.Close()
	t, err := template.New("keypair.tf.tmpl").ParseFiles("../../templates/aws/keypair.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(keyPairFile, k)
}
