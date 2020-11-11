package aws

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"path/filepath"
	"strings"
)

const (
	awsCredentialPath = ".aws/credentials"
)
func CreateAwsSessionConfig(provider *Provider) (*aws.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	credentialFile, err  := os.Open(filepath.Join(homeDir, awsCredentialPath))
	if err != nil {
		return nil, err
	}
	defer credentialFile.Close()
	scanner := bufio.NewScanner(credentialFile)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), provider.Profile) {
			if !scanner.Scan() {
				return nil, fmt.Errorf("not founnd credential")
			}
			accessKeyId := strings.Split(scanner.Text(), " ")[2]
			if !scanner.Scan() {
				return nil, fmt.Errorf("not founnd credential")
			}
			accessKeySecret := strings.Split(scanner.Text(), " ")[2]
			return &aws.Config{
				Credentials:                       credentials.NewStaticCredentialsFromCreds(
					credentials.Value{
						AccessKeyID:     accessKeyId,
						SecretAccessKey: accessKeySecret,
					}),
				Region:                            aws.String(provider.Region),
			}, nil
		}
	}
	return nil, fmt.Errorf("no such profile")
}
