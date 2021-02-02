package resource

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"text/template"
)

const (
	kakoi_vpn_id string = "kakoi-vpn-id"
	kakoi_ovpn_config_name string = "kakoi.ovpn"
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
	t, err := template.New("vpc.tf.tmpl").ParseFiles(filepath.Join("/etc/kakoi/templates/aws", "vpc.tf.tmpl"))
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
	AZ  string
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
		AZ:  network.Region + "a",
	}, nil
}

func (s *Subnet) BuildTemplate(workDir string) error {
	fileName := "subnet-" + s.Name + ".tf"
	file, err := os.Create(filepath.Join(workDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("subnet.tf.tmpl").ParseFiles(filepath.Join("/etc/kakoi/templates/aws", "subnet.tf.tmpl"))
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

func newVpn(cidr, domain string, subnet *Subnet) (*Vpn, error) {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &Vpn{
		Cidr:             cidrNet,
		Domain:           domain,
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
	t, err := template.New("vpn.tf.tmpl").ParseFiles("/etc/kakoi/templates/aws/vpn.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, v)
}

func (v *Vpn) SetPki(pki *Pki) {
	v.Pki = pki
}

func (v *Vpn) Create() error {
	return v.createOvpnConfig()
}


func (v *Vpn) createOvpnConfig() error {
	type ovpnConfig struct {
		Id     string
		Region string
		Addr   string
		Mask   string
		CaCert string
		Cert   string
		Key    string
	}
	// build path
	outputPath := filepath.Join(filepath.Join(v.Pki.Path[:len(v.Pki.Path)-3], "output"), "output.json")
	if _, err := os.Stat(outputPath); err != nil {
		return err
	}
	outputFile, err := os.Open(outputPath)
	if err != nil {
		return err
	}
	bytes := make([]byte, 512)
	l, err := outputFile.Read(bytes)
	if err != nil {
		return err
	}
	outputData := make(map[string]string)
	if err := json.Unmarshal(bytes[:l], &outputData); err != nil {
		return err
	}
	vpnId, err := getOutputValue(outputData, "kakoi-vpn-id")
	if err != nil {
		return err
	}
	caCertString, err := v.Pki.ReadCaCert()
	if err != nil {
		return err
	}
	clientCertString, err := v.Pki.ReadClientCert()
	if err != nil {
		return err
	}
	clientKeyString, err := v.Pki.ReadClientKey()
	if err != nil {
		return err
	}
	ipnet := v.AssociatedSubnet.Network.Cidr
	config := &ovpnConfig{
		Id:     vpnId,
		Region: v.AssociatedSubnet.Network.Region,
		Addr:   v.AssociatedSubnet.Network.Cidr.IP.String(),
		Mask:   net.IP{ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3]}.String(),
		CaCert: caCertString,
		Cert:   clientCertString,
		Key:    clientKeyString,
	}
	t, err := template.New("kakoi.ovpn.tmpl").ParseFiles("/etc/kakoi/templates/kakoi.ovpn.tmpl")
	if err != nil {
		return err
	}
	ovpnConfigFile, err := os.Create(kakoi_ovpn_config_name)
	if err != nil {
		return err
	}
	defer ovpnConfigFile.Close()
	if err := t.Execute(ovpnConfigFile, config); err != nil {
		return err
	}
	return nil
}

