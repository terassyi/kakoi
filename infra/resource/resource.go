package resource

import (
	"fmt"
	"github.com/terassyi/kakoi/config"
	"github.com/terassyi/kakoi/infra/aws"
	"path/filepath"
)

type Resource interface {
	BuildTemplate(workDir string) error
}

const template_path_aws string = "./templates/aws"

func New(conf *config.Config) ([]Resource, error) {
	resources := make([]Resource, 0, 100)
	switch conf.Provider.Name {
	case "aws":
		resources = append(resources, aws.NewProvider(conf.Provider.Profile, conf.Provider.Region))
	default:
		return nil, fmt.Errorf("unknown provider")
	}

	// vpc
	vpc, err := newNetwork(conf.Service.Network.Name, conf.Service.Network.Range, conf.Provider.Region)
	if err != nil {
		return nil, err
	}
	resources = append(resources, vpc)

	// subnet
	var vpnAssociatedSubnet *Subnet
	for _, s := range conf.Service.Network.Subnets {
		subnet, err := newSubnet(s.Name, s.Range, s.Private, vpc)
		if err != nil {
			return nil, err
		}
		if s.VpnGatewayAssociated {
			if conf.Service.Network.Vpn.AssociatedSubnet != s.Name {
				return nil, fmt.Errorf("not match vpn associating definition")
			}
			vpnAssociatedSubnet = subnet
		}
		resources = append(resources, subnet)
		// TODO route table settings
	}

	// TODO firewall settings

	// pki settings
	resources = append(resources, newPki(filepath.Join(conf.WorkDir, "pki"), conf.Service.Network.Vpn.Domain))

	// vpn
	if vpnAssociatedSubnet == nil {
		return nil, fmt.Errorf("vpn must be required")
	}
	vpn, err := newVpn(conf.Service.Network.Vpn.Range, conf.Service.Network.Vpn.Domain, vpnAssociatedSubnet)
	if err != nil {
		return nil, err
	}
	resources = append(resources, vpn)

	// server
	// key pair
	resources = append(resources, newKeyPair(filepath.Join(conf.WorkDir, "keys"), conf.Service.KeyPair))

	// TODO hosts
	//for _, host := range conf.Service.Servers {
	//
	//}

	return resources, nil
}

func findSubnet(resources []Resource, name string) *Subnet {
	for _, resource := range resources {
		switch r := resource.(type) {
		case *Subnet:
			if r.Name == name {
				return r
			}
		default:
		}
	}
	return nil
}
