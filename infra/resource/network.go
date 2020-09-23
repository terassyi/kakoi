package resource

import (
	"net"
	"os"
	"path/filepath"
	"text/template"
)

type Network struct {
	Name   string
	Cidr   *net.IPNet
	Region string
}

func newNetwork(name, cidr, region string) (*Network, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &Network{
		Name:   name,
		Cidr:   ipNet,
		Region: region,
	}, nil
}

func (n *Network) BuildTemplate(workDir string) error {
	fileName := "vpc-" + n.Name + ".tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("vpc.tf.tmpl").ParseFiles(filepath.Join("../../templates/aws", "vpc.tf.tmpl"))
	if err != nil {
		return err
	}
	return t.Execute(file, n)
}

type Subnet struct {
	Name    string
	Network *Network
	Cidr    *net.IPNet
	Private bool
	Region  string
}

func newSubnet(name, cidr string, private bool, network *Network) (*Subnet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &Subnet{
		Name:    name,
		Network: network,
		Cidr:    ipNet,
		Private: private,
		Region:  network.Region,
	}, nil
}

func (s *Subnet) BuildTemplate(workDir string) error {
	fileName := "subnet-" + s.Name + ".tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("subnet.tf.tmpl").ParseFiles(filepath.Join("../../templates/aws", "subnet.tf.tmpl"))
	if err != nil {
		return err
	}
	return t.Execute(file, s)
}

type Route struct {
	from *net.IPNet
	to   *net.IPNet
}

func newRoute(from, to string) (*Route, error) {
	_, fromNet, err := net.ParseCIDR(from)
	if err != nil {
		return nil, err
	}
	_, toNet, err := net.ParseCIDR(to)
	if err != nil {
		return nil, err
	}
	return &Route{
		from: fromNet,
		to:   toNet,
	}, nil
}

type Vpn struct {
	Cidr             *net.IPNet
	Domain           string
	AssociatedSubnet *Subnet
	Pki              *Pki
}

func newVpn(cidr string, subnet *Subnet) (*Vpn, error) {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &Vpn{
		Cidr:             cidrNet,
		AssociatedSubnet: subnet,
	}, nil
}

func (v *Vpn) BuildTemplate(workDir string) error {
	fileName := "kakoi-vpn.tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("vpn.tf.tmpl").ParseFiles("../../templates/aws/vpn.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, v)
}
