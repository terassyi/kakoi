package infra

import (
	"fmt"
	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/terassyi/kakoi/infra/aws"
	"github.com/terassyi/kakoi/infra/resource"
	"github.com/terassyi/kakoi/infra/state"
	"os"
	"path/filepath"
)

const (
	kakoi_dir string = ".kakoi"
	ext_yaml  string = ".yaml"
	ext_yml   string = ".yml"
	ext_json string = ".json"
)

type Initializer interface {
	Init() error
}

type initializer struct {
	workDir string
	conf *state.State
}

func NewInitializer(path string) (Initializer, error) {
	dir, file := filepath.Split(path)
	if err := state.ValidateExtName(file); err != nil {
		return nil, err
	}
	workDir, err := createWorkDir(dir)
	if err != nil {
		return nil, err
	}
	parser, err := state.NewParser(workDir, path)
	if err != nil {
		return nil, err
	}
	s, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &initializer{workDir: workDir, conf: s}, nil
}

func (i *initializer) init() error {

	var resources []resource.Resource
	switch i.conf.Provider.Name {
	case "aws":
		resources = append(resources, aws.NewProvider(i.conf.Provider.Profile, i.conf.Provider.Region))
		//resources = append(resources, aws.NewS3(conf.Service.Name))
	default:
		return fmt.Errorf("unknown provider")
	}
	// create storage
	storage := aws.NewS3(i.conf.Service.Name)
	resources = append(resources, storage)
	r, err := i.createImageUploader()
	if err != nil {
		return err
	}
	resources = append(resources, r...)
	builder, err := newBuilderFromResources(i.workDir, resources)
	if err != nil {
		return err
	}
	if err := builder.buildTemplate(); err != nil {
		return err
	}
	if err := builder.create(); err != nil {
		return err
	}

	ids, err := i.importImage(storage)
	if err != nil {
		return err
	}
	fmt.Println(ids)
	// create state file
	if err := i.conf.CreateState(); err != nil {
		return err
	}
	return nil
}

func (i *initializer) importImage(storage *aws.S3) (map[string]string, error) {
	const imagesBase = "images/"
	importTaskIds := make(map[string]string)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           i.conf.Provider.Profile,
		SharedConfigState: session.SharedConfigEnable,
	}))
	ins := ec2.New(sess)
	for _, s := range i.conf.Service.Hosts.Servers {
		if s.Image.ImagePath != "" {
			// convert to ami
			ext := filepath.Ext(s.Image.ImagePath)
			fmt.Println("s3 path=", filepath.Join(storage.BucketName(), imagesBase, s.Name))
			input := &ec2.ImportImageInput{
				Architecture: awsSdk.String("x86_64"),
				Description:  awsSdk.String("kakoi vulnerable image"),
				DiskContainers: []*ec2.ImageDiskContainer{
					{
						Description: awsSdk.String(s.Name),
						//DeviceName:  awsSdk.String(""),
						Format: awsSdk.String(convertExtImageFormat(ext)),
						//Url:        awsSdk.String(filepath.Join(imagesBase, s.Name)),
						UserBucket: &ec2.UserBucket{
							S3Bucket: awsSdk.String(storage.BucketName()),
							S3Key:    awsSdk.String(filepath.Join(imagesBase, filepath.Base(s.Image.ImagePath))),
						},
					},
				},
				Platform: awsSdk.String("Linux"),
				RoleName: nil,
			}
			output, err := ins.ImportImage(input)
			if err != nil {
				return nil, err
			}
			fmt.Printf("[INFO] import task for %s is %s\n", s.Name, *(output.ImportTaskId))
			importTaskIds[s.Name] = *output.ImportTaskId
		}
	}
	// wait for importing image

	idMap, err := aws.WaitImportImageResult(i.conf.Provider.Profile, importTaskIds)
	if err != nil {
		return nil, err
	}
	for _, s := range i.conf.Service.Hosts.Servers {
		id, ok := idMap[s.Name]
		if ok {
			if s.Image.Id == "" {
				s.Image.Id = id
			}
		}
	}
	return idMap, nil
}

func (i *initializer) createImageUploader() ([]resource.Resource, error) {
	const imagesBase = "images/"
	var imageResources []resource.Resource
	for _, s := range i.conf.Service.Hosts.Servers {
		if s.Image.ImagePath != "" {
			imagePath, err := i.buildAbsPath(s.Image.ImagePath)
			if err != nil {
				return nil, err
			}
			i := aws.NewS3Uploader(imagePath, filepath.Join(imagesBase, filepath.Base(s.Image.ImagePath)))
			imageResources = append(imageResources, i)
		}
		// image builder files
		if s.Image.ScriptFilePath != nil {
			base, err := absWorkDir(i.workDir)
			if err != nil {
				return nil, err
			}
			ib, err := resource.NewImageBuilder(s.Name, i.conf.Provider.Region, base,nil, s.Image.ScriptFilePath)
			if err != nil {
				return nil, err
			}
			imageResources = append(imageResources, ib)
			fmt.Printf("[info] custom image build for %v\n", s.Name)
		}
	}
	return imageResources, nil
}

func convertExtImageFormat(ext string) string {
	switch ext {
	case ".ova":
		return "OVA"
	case ".vmdk":
		return "VMDK"
	case ".vhd":
		return "VHD"
	case ".vhdx":
		return "VHDX"
	default:
		return ""
	}
}

func (i *initializer) buildImage(storage *aws.S3) error {

	return nil
}

func (i *initializer) Init() error {
	return i.init()
}

func createWorkDir(path string) (string, error) {
	// if work on current dir, path = ""
	workPath := filepath.Join(path, kakoi_dir)
	if err := os.MkdirAll(workPath, 0755); err != nil {
		return "", err
	}
	// pki cert files
	if err := os.MkdirAll(filepath.Join(workPath, "pki"), 0755); err != nil {
		return "", err
	}
	// server key pair
	if err := os.MkdirAll(filepath.Join(workPath, "keys"), 0755); err != nil {
		return "", err
	}
	// output files
	if err := os.MkdirAll(filepath.Join(workPath, "output"), 0755); err != nil {
		return "", err
	}
	// image files
	if err := os.MkdirAll(filepath.Join(workPath, "images"), 0755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(workPath, "storage"), 0755); err != nil {
		return "", err
	}
	return workPath, nil
}

func isAbsPath(path string) bool {
	if path[0] == '/' {
		return true
	}
	return false
}

func (i *initializer) buildAbsPath(path string) (string, error) {
	if isAbsPath(path) {
		return path, nil
	}
	base := i.workDir[:len(i.workDir)-7]
	if isAbsPath(i.workDir) {
		return filepath.Join(base, path), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, base, path), nil
}

func absWorkDir(workDir string) (string, error) {
	if isAbsPath(workDir) {
		return workDir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, workDir), nil
}