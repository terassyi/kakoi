package config

type Config struct {
	WorkDir  string
	Provider Provider `yaml:"provider"`
	Service  Service  `yaml:"service"`
}

type Service struct {
	Name    string  `yaml:"name"`
	Network Network `yaml:"network"`
	Servers []Server  `yaml:"servers"`
	KeyPair string  `yaml:"key_pair_name"`
}

type Provider struct {
	Name    string `yaml:"name"`
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
}

type Network struct {
	Name    string     `yaml:"name"`
	Range   string     `yaml:"range"`
	Subnets []Subnet   `yaml:"subnets"`
	Vpn     VpnGateway `yaml:"vpn_gateway"`
}

type Subnet struct {
	Name                 string  `yaml:"name"`
	Range                string  `yaml:"range"`
	Private              bool    `yaml:"private"`
	VpnGatewayAssociated bool    `yaml:"vpn_gateway_associated"`
	Routes               []Route `yaml:"routes"`
}

type Route struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type VpnGateway struct {
	Range            string `yaml:"range"`
	Domain           string `yaml:"domain"`
	AssociatedSubnet string `yaml:"associated_subnet"`
}

type Server struct {
	Name   string `yaml:"name"`
	Subnet string `yaml:"subnet"`
	Image ServerImage `yaml:"image"`
	Ports  []int  `yaml:"ports"`
}

type ServerImage struct {
	Custom bool `yaml:"custom"`
	ImagePath string `yaml:"image_path"`
	ScriptFilePath []string `yaml:"scripts"`
	InlineScripts []string `yaml:"inline"`
}
