package file

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// ReadFromYAML reads the YAML file and pass to the object
// args:
//
//	path: file path location
//	target: object which will hold the value
//
// returns:
//
//	error: operation state error
func ReadFromYAML(path string, target interface{}) error {
	yf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yf, target)
}
