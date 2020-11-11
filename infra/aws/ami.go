package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"time"
)

func WaitImportImageResult(profile string, ids map[string]string) (map[string]string, error){
	var idsPtr []*string
	taskIdImageIdMap := make(map[string]string)
	for k, _ := range ids {
		idsPtr = append(idsPtr, aws.String(ids[k]))
	}

	fmt.Println("ids ptr", idsPtr)
	errCh := make(chan error)
	counter := 0
	timer := time.NewTicker(time.Minute)
	go func(){
		for {
			<-timer.C
			fmt.Println("[INFO] checking for building status")
			output, err := checkImportResult(profile, idsPtr)
			if err != nil {
				errCh <- err
			}
			for _, o := range output.ImportImageTasks {
				id, ok, err := checkImportImageOutput(o)
				if err != nil {
					errCh <- err
					return
				}
				if ok {
					if _, ok := taskIdImageIdMap[*o.ImportTaskId]; !ok {
						counter += 1
						taskIdImageIdMap[*o.ImportTaskId] = id
						fmt.Printf("[INFO] image build completed id=%s(task id=%s)\n", id, *o.ImportTaskId)
						fmt.Printf("[DEBUG] counter=%d ids_len=%d\n", counter, len(ids))
					}
				}
			}
			if counter >= len(ids) {
				fmt.Println("[INFO] finished building image!")
				errCh <- nil
				return
			}
		}
	}()

	err := <- errCh
	if err != nil {
		return nil, err
	}
	resImageMap := make(map[string]string)
	for k, v := range ids {
		id, ok := taskIdImageIdMap[v]
		if !ok {
			fmt.Printf("[ERROR] failed to find image id(%s) for image task id\n", v)
			continue
		}
		resImageMap[k] = id
	}
	return resImageMap, nil
}

func checkImportResult(profile string, ids []*string) (*ec2.DescribeImportImageTasksOutput, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}))
	ins := ec2.New(sess)
	return ins.DescribeImportImageTasks(&ec2.DescribeImportImageTasksInput{
		ImportTaskIds: ids,
	})
}

func checkImportImageOutput(out *ec2.ImportImageTask) (string, bool, error) {
	if out.StatusMessage != nil {
		fmt.Printf("[INFO] task %s's status is %s: %s \n", *out.ImportTaskId, *out.Status, *out.StatusMessage)
	}
	if *out.Status == "completed" && out.ImageId != nil {
		return *out.ImageId, true, nil
	}
	if *out.Status == "deleting" || *out.Status == "deleted" {
		return "", false, fmt.Errorf("failed to import image (id=%s)", *out.ImportTaskId)
	}
	if out.ImageId == nil {
		return "", false, nil
	}


	return "", false, nil
}
