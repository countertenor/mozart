package yaml

import (
	"fmt"
	"io/ioutil"
	"os"
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

func ParseCommonFolder(config map[string]interface{}, dirName string) error {
	statikFS, err := statik.GetStaticFS(statik.Template)
	if err != nil {
		return err
	}
	f, _ := fs.ReadFile(statikFS, "/common-files/bnr/stest")
	fmt.Println("f : ", f)
	err = fs.Walk(statikFS, dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println("info : ", info.Name())
		if !info.IsDir() {
			fileName := info.Name()
			if strings.Contains(fileName, "-") {
				return fmt.Errorf("filename %v cannot contain '-', kindly remove or replace with '_'", fileName)
			}
			fmt.Println("file : ", path)
			staticFile, err := statikFS.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file %v err: %v", fileName, err)
			}
			fileData, err := ioutil.ReadAll(staticFile)
			if err != nil {
				return fmt.Errorf("could not read file %v err: %v", fileName, err)
			}
			config[fileName[:strings.LastIndex(fileName, ".")]] = string(fileData)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error getting static files : %v", err)
	}
	return nil
}
