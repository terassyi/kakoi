package resource

type Resource interface {
	BuildTemplate(workDir string) error
}

const templatePathAws string = "./templates/aws"
