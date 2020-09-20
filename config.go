package main

import (
	"io/ioutil"
	"net"
	"os"
	"text/template"
)

type ovpnConfig struct {
	Id     string
	Region string
	Addr   string
	Mask   string
	CaCert string
	Cert   string
	Key    string
}

func outputOvpnConfig(id, region, cidr string) error {
	// parse ip cidr
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	mask := net.IP{ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3]}

	// certificate files
	caData, err := ioutil.ReadFile("./data/ca01.crt")
	if err != nil {
		return err
	}

	clientCertData, err := ioutil.ReadFile("./data/client.crt")
	if err != nil {
		return err
	}

	clientKeyData, err := ioutil.ReadFile("./data/client.key")
	if err != nil {
		return err
	}

	config := &ovpnConfig{
		Id:     id,
		Region: region,
		Addr:   ipnet.IP.String(),
		Mask:   mask.String(),
		CaCert: string(caData),
		Cert:   string(clientCertData),
		Key:    string(clientKeyData),
	}
	t, err := template.New("kakoi.ovpn.tmpl").ParseFiles("templates/kakoi.ovpn.tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create("./kakoi.ovpn")
	if err != nil {
		return err
	}
	defer file.Close()
	if err := t.Execute(file, config); err != nil {
		return err
	}
	return nil
}
