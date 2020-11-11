package aws

import "testing"

func TestS3_UploadFile(t *testing.T) {
	s := NewS3("kakoi.example")
	//if err := s.UploadFile("kakoi", "../../data/upload-test.txt", "upload-test.txt"); err != nil {
	if err := s.UploadFile("kakoi", "/Users/spectre/VirtualMachine/nezukochan.ova", "images/nezukochan.ova"); err != nil {
		t.Fatal(err)
	}
}