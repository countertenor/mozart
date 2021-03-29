package yaml

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

func ParseCommonFolder(config map[string]interface{}, dirName string) error {
	err := static.Walk(static.ResourceType, dirName, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileName := info.Name()
			if strings.Contains(fileName, "-") {
				return fmt.Errorf("filename %v cannot contain '-', kindly remove or replace with '_'", fileName)
			}
			// fmt.Println("file : ", fileName)
			staticFile, err := static.OpenFileFromStaticFS(static.ResourceType, path)
			if err != nil {
				return fmt.Errorf("error opening file %v err: %v", fileName, err)
			}
			fileData, err := ioutil.ReadAll(staticFile)
			if err != nil {
				return fmt.Errorf("could not read file %v err: %v", fileName, err)
			}
			config[strings.TrimSuffix(fileName, filepath.Ext(fileName))] = string(fileData)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error getting static files : %v", err)
	}
	return nil
}
