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
