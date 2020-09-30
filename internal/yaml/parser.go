package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/prashantgupta24/mozart/internal/config"
	"github.com/prashantgupta24/mozart/statik"
	"gopkg.in/yaml.v2"
)

//ParseDefault parses default file into config
func ParseDefault(config *config.Instance, defaultFileName string) error {
	defaultFile, err := statik.OpenFile("/" + defaultFileName)
	if err != nil {
		return err
	}
	defer defaultFile.Close()
	fileData, err := ioutil.ReadAll(defaultFile)
	if err != nil {
		return fmt.Errorf("could not read file %v err: %v", defaultFileName, err)
	}

	err = yaml.Unmarshal([]byte(fileData), &config)
	if err != nil {
		return fmt.Errorf("error while unmarshalling yaml %v: %v", defaultFileName, err)
	}
	return nil
}

//ParseFile parses file into config
func ParseFile(config *config.Instance, filename string) error {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no yaml file called %v found. Run 'init' to generate the sample config yaml file", filename)
		}
		return fmt.Errorf("error while reading conf file %v : %v", filename, err)
	}

	err = yaml.Unmarshal([]byte(fileData), &config)
	if err != nil {
		return fmt.Errorf("error while unmarshalling yaml %v: %v", filename, err)
	}
	return nil
}

//PreCheck handles validation and pre-configuration
func PreCheck(config *config.Instance) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	config.PreconfigureFields()
	return nil
}
