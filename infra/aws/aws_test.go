package aws

import (
	"testing"
)

func TestCreateAwsSession(t *testing.T) {
	provider := &Provider{
		Profile: "kakoi",
		Region:  "ap-northeast-1",
		Name:    "test-project",
	}
	c, err := CreateAwsSessionConfig(provider)
	if err != nil {
		t.Fatal(err)
	}
	if *c.Region != provider.Region {
		t.Fatal(*c.Region)
	}
}

func TestWaitImageBuildResult(t *testing.T) {
	_, err := WaitImageBuildResult("kakoi", []string{"kakoi-example-host1:ace01887-31c7-47ac-b288-9d07df25b564", "kakoi-example-host1:936af007-84c1-48dd-89c2-1222ef7e3f15"})
	if err != nil {
		t.Fatal(err)
	}
}
