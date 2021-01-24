package terraform

import (
	tfjson "github.com/hashicorp/terraform-json"
	"os"
)

func ToJson(file string) (*tfjson.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data := make([]byte, 1024)
	if _, err := f.Read(data); err != nil {
		return nil, err
	}
	conf := &tfjson.Config{}
	if err := conf.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	//fmt.Println(conf)
	return conf, nil
}