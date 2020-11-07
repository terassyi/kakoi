package state

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const kakoi_state_file string = "kakoi.state"
func (s *State) CreateState() error {
	//if _, err := os.Stat(filepath.Join(s.WorkDir, kakoi_state_file)); err != nil {
	//	return err
	//}
	data, err := s.output()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.WorkDir, kakoi_state_file), data, 0755)
}

func (s *State) output() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}