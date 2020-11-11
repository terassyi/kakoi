package state

type State struct {
	WorkDir  string
	Provider Provider `json:"provider"`
	Service  Service  `json:"service"`
}

type Service struct {
	Name    string  `json:"name"`
	Storage *Storage `json:"storage"`
	Network *Network `json:"network"`
	Hosts *Hosts  `json:"hosts"`
}

type Provider struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Region  string `json:"region"`
}

type Storage struct {
	Id string `json:"id"`
}

type Network struct {
	Name    string     `json:"name"`
	Range   string     `json:"range"`
	Subnets []Subnet   `json:"subnets"`
	Vpn     VpnGateway `json:"vpn_gateway"`
}

type Subnet struct {
	Name                 string  `json:"name"`
	Range                string  `json:"range"`
	Private              bool    `json:"private"`
	VpnGatewayAssociated bool    `json:"vpn_gateway_associated"`
	Routes               []Route `json:"routes"`
}

type Route struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type VpnGateway struct {
	Id string `json:"id"`
	Range            string `json:"range"`
	Domain           string `json:"domain"`
	AssociatedSubnet string `json:"associated_subnet"`  // TODO 重複してる
	Cert Cert `json:"cert"`
}

type Cert struct {
	Status string `json:"status"`
	Path string `json:"path"`
}

type KeyPair struct {
	Name string `json:"name"`
	Status string `json:"status"` // automatically generated
}

type Hosts struct {
	KeyPair KeyPair `json:"key"`
	Servers []*Server `json:"servers"`
}

type Server struct {
	Name   string `json:"name"`
	Size string `json:"size"`
	Number int `json:"number"` // default 1
	Subnet string `json:"subnet"`
	Image *ServerImage `json:"image"`
	Status string `json:"status"` // automatically generated
	Ports  []int  `json:"ports"`
}

type ServerImage struct {
	Custom bool `yaml:"custom"`
	Id string `json:"id"`
	Status string `json:"status"`
	ImagePath string `json:"image_path"`
	ScriptFilePath []string `json:"scripts"`
	InlineScripts []string `json:"inline"`
}
