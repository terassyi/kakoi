package resource

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	builderDesc string = "instance image for kakoi exercise environment."
)

type ImageBuilder struct {
	Region string
	Name string
	Scripts []string
}

type packerBuilder struct {
	Builders []awsBuilder `json:"builders"`
	Provisioners []provisioner `json:"provisioners"`
}

type awsBuilder struct {
	Type string `json:"type"`
	Region string `json:"region"`
	InstanceType string `json:"instance_type"`
	UserName string `json:"ssh_username"`
	AmiName string `json:"ami_name"`
	AmiDescription string `json:"ami_description"`
	PublicIp bool `json:"associate_public_ip_address"`
	Filters awsSourceAmiFilter `json:"source_ami_filter"`
}

type awsSourceAmiFilter struct {
	Filters awsSourceFilterImpl `json:"filters"`
	MostRecent bool `json:"most_recent"`
	Owners []string `json:"owners"`
}

type awsSourceFilterImpl struct {
	Type string `json:"virtualization_type"`
	Name string `json:"name"`
	RootDeviceType string `json:"root_device_type"`
}

type provisioner struct {
	Type string `json:"type"`
	Inline []string `json:"inline"`
	Scripts []string `json:"scripts"`
}

func NewPackerBuilder(region, name string, commands, files []string) (*packerBuilder, error) {
	builder := newAwsBuilder(region, name)
	var filenames []string
	for _, p := range files {
		if _, err := os.Stat(p); err != nil {
			return nil, err
		}
		filenames = append(filenames, filepath.Base(p))
	}
	prov, err := newProvisioner(commands, filenames)
	if err != nil {
		return nil, err
	}
	return &packerBuilder{
		Builders:     []awsBuilder{*builder},
		Provisioners: []provisioner{ *prov },
	}, nil
}

func newAwsBuilder(region, name string) *awsBuilder {
	return &awsBuilder{
		Type:           "amazon-ebs",
		Region:         region,
		InstanceType:   "t2.micro",
		UserName:       "ec2-user",
		AmiName:        "kakoi-" + name,
		AmiDescription: builderDesc,
		PublicIp:       true,
		Filters:        *newAwsSourceFilter(),
	}
}

func newAwsSourceFilter() *awsSourceAmiFilter {
	return &awsSourceAmiFilter{
		Filters:    awsSourceFilterImpl{
			Type: "hvm",
			Name: "amzn-ami*-ebs",
			RootDeviceType: "ebs",
		},
		MostRecent: true,
		Owners:     []string{"137112412989"},
	}
}

func newProvisioner(commands, files []string) (*provisioner, error) {
	if commands != nil {
		return &provisioner{
			Type:    "shell",
			Inline:  commands,
		}, nil
	}
	if files != nil {
		return &provisioner{
			Type:    "shell",
			Inline:  nil,
			Scripts: files,
		}, nil
	}
	return nil, fmt.Errorf("provisioning script is not set.")
}