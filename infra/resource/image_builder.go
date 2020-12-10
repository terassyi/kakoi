package resource

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	builderDesc      string = "instance image for kakoi exercise environment."
	packer_file_name string = "image_builder.json"
)

var (
	build_spec_pre_build_commands = []string{
		"echo \"installing hashicorp packer\"",
		"curl -qL -o packer.zip https://releases.hashicorp.com/packer/1.6.4/packer_1.6.4_linux_amd64.zip && unzip packer.zip",
		"echo \"installing jq\"",
		"curl -qL -o jq https://stedolan.github.io/jq/download/linux64/jq && chmod +x ./jq",
		"echo \"validate packer cofiguration file\"",
		"./packer validate {{ .Path }}",
	}
	build_spec_pre_build_path_index = 5

	build_spec_build_commands = []string{
		"curl -qL -o aws_credentials.json http://169.254.170.2/$AWS_CONTAINER_CREDENTIALS_RELATIVE_URI > aws_credentials.json",
		"aws configure set region $AWS_REGION",
		"aws configure set aws_access_key_id `./jq -r '.AccessKeyId' aws_credentials.json`",
		"aws configure set aws_secret_access_key `./jq -r '.SecretAccessKey' aws_credentials.json`",
		"aws configure set aws_session_token `./jq -r '.Token' aws_credentials.json`",
		"echo \"building image\"",
		"./packer build {{ .Path }}",
	}
	build_spec_build_path_index = 6

	buildSpecTemplate = make(map[string]interface{})
)

func createBuildSpec(path, name string) error {
	dstPath := filepath.Join(name, packer_file_name)
	placeHolder := "{{ .Path }}"
	build_spec_pre_build_commands[build_spec_pre_build_path_index] = strings.Replace(build_spec_pre_build_commands[build_spec_pre_build_path_index], placeHolder, dstPath, -1)
	build_spec_build_commands[build_spec_build_path_index] = strings.Replace(build_spec_build_commands[build_spec_build_path_index], placeHolder, dstPath, -1)

	type phase struct {
		Commands []string
	}

	type buildSpecPhases struct {
		PreBuild phase `yaml:"pre_build"`
		Build    phase `yaml:"build"`
	}
	type buildTemplate struct {
		Version float32
		Phases  buildSpecPhases
	}

	buildTempl := buildTemplate{
		Version: 0.2,
		Phases: buildSpecPhases{
			PreBuild: phase{Commands: build_spec_pre_build_commands},
			Build:    phase{Commands: build_spec_build_commands},
		},
	}

	data, err := yaml.Marshal(buildTempl)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, "buildspec.yml"), data, 0666)
}

type ImageBuilder struct {
	Region   string
	Name     string
	Files    []string
	Commands []string
}

func NewImageBuilder(name, region string, commands, files []string) *ImageBuilder {
	return &ImageBuilder{
		Region:   region,
		Name:     name,
		Files:    files,
		Commands: commands,
	}
}

func (i *ImageBuilder) BuildTemplate(workDir string) error {
	fileName := "image_builder-" + i.Name + ".tf"

	// create buildspec and packer config file
	packer, err := i.createPackerBuilder()
	if err != nil {
		return err
	}
	path := filepath.Join(workDir, "images", i.Name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	if err := packer.outputJson(filepath.Join(path, packer_file_name)); err != nil {
		fmt.Println("json")
		return err
	}
	if err := createBuildSpec(path, i.Name); err != nil {
		fmt.Println("buildspec")
		return err
	}
	file, err := os.Create(filepath.Join(workDir, "images", fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	t, err := template.New("image_builder.tf.tmpl").ParseFiles("../../templates/aws/image_builder.tf.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(file, i)
}

func (i *ImageBuilder) createPackerBuilder() (*packerBuilder, error) {
	return newPackerBuilder(i.Name, i.Region, i.Commands, i.Files)
}

type packerBuilder struct {
	Builders     []awsBuilder  `json:"builders"`
	Provisioners []provisioner `json:"provisioners"`
}

func (p *packerBuilder) outputJson(path string) error {
	//data, err := json.Marshal(p)
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

type awsBuilder struct {
	Type           string             `json:"type"`
	Region         string             `json:"region"`
	InstanceType   string             `json:"instance_type"`
	UserName       string             `json:"ssh_username"`
	AmiName        string             `json:"ami_name"`
	AmiDescription string             `json:"ami_description"`
	PublicIp       bool               `json:"associate_public_ip_address"`
	Filters        awsSourceAmiFilter `json:"source_ami_filter"`
}

type awsSourceAmiFilter struct {
	Filters    awsSourceFilterImpl `json:"filters"`
	MostRecent bool                `json:"most_recent"`
	Owners     []string            `json:"owners"`
}

type awsSourceFilterImpl struct {
	Type           string `json:"virtualization_type"`
	Name           string `json:"name"`
	RootDeviceType string `json:"root_device_type"`
}

type provisioner struct {
	Type    string   `json:"type"`
	Inline  []string `json:"inline"`
	Scripts []string `json:"scripts"`
}

func newPackerBuilder(name, region string, commands, files []string) (*packerBuilder, error) {
	builder := newAwsBuilder(region, name)
	var filenames []string
	for _, p := range files {
		//if _, err := os.Stat(p); err != nil {
		//	return nil, err
		//}
		filenames = append(filenames, filepath.Base(p))
	}
	prov, err := newProvisioner(commands, filenames)
	if err != nil {
		return nil, err
	}
	return &packerBuilder{
		Builders:     []awsBuilder{*builder},
		Provisioners: []provisioner{*prov},
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
		Filters: awsSourceFilterImpl{
			Type:           "hvm",
			Name:           "amzn-ami*-ebs",
			RootDeviceType: "ebs",
		},
		MostRecent: true,
		Owners:     []string{"137112412989"},
	}
}

func newProvisioner(commands, files []string) (*provisioner, error) {
	if commands != nil {
		return &provisioner{
			Type:   "shell",
			Inline: commands,
		}, nil
	}
	if files != nil {
		return &provisioner{
			Type:    "shell",
			Scripts: files,
		}, nil
	}
	return nil, fmt.Errorf("provisioning script is not set.")
}
