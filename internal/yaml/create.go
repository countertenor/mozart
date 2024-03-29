package yaml

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/countertenor/mozart/static"
	"gopkg.in/yaml.v2"
)

//CreateSampleConfigFile creates sample config file
func CreateSampleConfigFile(sampleFileName string) error {

	//check if file exists already
	filesInDir, err := ioutil.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error while reading current dir, err : %v", err)
	}
	for _, file := range filesInDir {
		if file.Name() == sampleFileName {
			return fmt.Errorf("sample file already exists in the current directory. Kindly delete it then run command again")
		}
	}

	//Read sample file
	sampleFile, err := static.OpenFileFromStaticFS(static.ResourceType, sampleFileName)
	if err != nil {
		return err
	}

	defer sampleFile.Close()
	fileData, err := ioutil.ReadAll(sampleFile)
	if err != nil {
		return fmt.Errorf("could not read file %v err: %v", sampleFileName, err)
	}

	err = ioutil.WriteFile(sampleFileName, fileData, 0644)
	if err != nil {
		return fmt.Errorf("could not write sample file %v, err : %v", sampleFileName, err)
	}
	return nil
}

//CreateFromConfig creates the yaml file from config
func CreateFromConfig(config map[string]interface{}, filename string) error {

	filename = strings.TrimSpace(filename)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create %v file, err : %v", filename, err)
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error while marshalling data to file %v: %v", filename, err)
	}

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error while writing to file %v: %v", filename, err)
	}

	return nil
}
