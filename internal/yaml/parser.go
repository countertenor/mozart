package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/prashantgupta24/mozart/internal/config"
	"github.com/prashantgupta24/mozart/statik"
	"gopkg.in/yaml.v2"
)

//Parse function first parses the default configuration file, which loads the
//default values. After that, it parses the file passed in as a param,
//thereby preserving default values if they are not over-written.
func Parse(config *config.Instance, filename, defaultFileName string) error {

	//Read default file
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

	//Read file passed as param
	fileData, err = ioutil.ReadFile(filename)
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

	err = config.Validate()
	if err != nil {
		return fmt.Errorf("error while validating YAML %v: %v", filename, err)
	}

	return nil
}
