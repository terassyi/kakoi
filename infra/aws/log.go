package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/terassyi/kakoi/infra/log"
)

const (
	log_path = "/var/log/kakoi"
)

func GetLog(workDir, profile, name, prefix, id string) ([]string, error) {
	streamName := "kakoi-" + name + "-" + prefix + "/" + id
	fmt.Println(streamName)
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	cw := cloudwatchlogs.New(sess)
	out, err := cw.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		EndTime:       nil,
		Limit:         nil,
		LogGroupName:  aws.String("kakoi"),
		LogStreamName: aws.String(streamName),
		NextToken:     nil,
		StartFromHead: nil,
		StartTime:     nil,
	})
	if err != nil {
		return nil, err
	}
	l, err := formatLogs(out)
	if err != nil {
		return nil, err
	}
	if err := log.OutputLogs(workDir, name, l); err != nil {
		return nil, err
	}
	return l, nil
}

func formatLogs(logs *cloudwatchlogs.GetLogEventsOutput) ([]string, error) {
	var l []string
	for _, event := range logs.Events {
		l = append(l, *event.Message)
	}
	return l, nil
}
