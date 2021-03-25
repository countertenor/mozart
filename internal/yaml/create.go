package yaml

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed mozart-sample.yaml
var embedSampleFile embed.FS

//CreateSampleConfigFile creates sample config file
func CreateSampleConfigFile(sampleFileName string) error {

	//Read sample file
	// sampleFile, err := statik.OpenFileFromStaticFS(statik.Template, "/"+sampleFileName)
	sampleFile, err := embedSampleFile.Open("mozart-sample.yaml")
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
