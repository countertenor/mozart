package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/countertenor/mozart/static"
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

//ParseFileFromStatic parses file from static into config
func ParseFileFromStatic(config map[string]interface{}, filename string) error {
	staticFile, err := static.OpenFileFromStaticFS(static.ResourceType, filename)
	if err != nil {
		return err
	}
	defer staticFile.Close()

	fileData, err := ioutil.ReadAll(staticFile)
	if err != nil {
		return fmt.Errorf("could not read file %v err: %v", filename, err)
	}
	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		return fmt.Errorf("error while unmarshalling yaml %v: %v", filename, err)
	}
	return nil
}
