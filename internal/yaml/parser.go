package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//ParseFile parses file into config
func ParseFile(config map[string]interface{}, filename string) error {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error while reading conf file %v : %v", filename, err)
	}

	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		return fmt.Errorf("error while unmarshalling yaml %v: %v", filename, err)
	}
	return nil
}
