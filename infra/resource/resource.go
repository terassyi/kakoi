package resource

import "net"

type Network struct {
	name string
	cidr net.IPNet
}
