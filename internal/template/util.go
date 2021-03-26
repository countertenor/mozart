package template

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/countertenor/mozart/statik"
	"github.com/rakyll/statik/fs"
)

type templateInstance struct {
	scriptName       string
	scriptFileName   string
	templateFileName string
}

//Generate conf files based on input yaml
func Generate(conf map[string]interface{}, dirToGenerate, templateDir, generatedDir string) error {
	fmt.Println("")

	templatesToGenerate, err := getTemplatesToGenerate(dirToGenerate, templateDir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorCh := make(chan error)
	for _, template := range templatesToGenerate {
		wg.Add(1)
		go func(template templateInstance) {
			err := generateTemplate(template.scriptName, template.scriptFileName, template.templateFileName, generatedDir, conf)
			if err != nil {
				errorCh <- err
			}
			wg.Done()
		}(template)
	}

	go func() {
		wg.Wait()
		close(errorCh)
	}()
	for err := range errorCh {
		return err
	}
	return nil
}

func getTemplatesToGenerate(dirToGenerate, templateDir string) ([]templateInstance, error) {
	var templates []templateInstance

	statikFS, err := statik.GetStaticFS(statik.Template)
	if err != nil {
		return nil, err
	}

	err = fs.Walk(statikFS, templateDir+dirToGenerate, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileName := info.Name()
			fileExt := filepath.Ext(fileName)
			templateInstance := templateInstance{
				scriptName:       strings.TrimSuffix(fileName, fileExt),
				scriptFileName:   strings.TrimPrefix(path, templateDir),
				templateFileName: path,
			}
			templates = append(templates, templateInstance)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error getting static files : %v", err)
	}
	return templates, nil
}

func generateTemplate(scriptName, scriptFileName, templateFileName, generatedDir string, config map[string]interface{}) error {

	templateFile, err := statik.OpenFileFromStaticFS(statik.Template, templateFileName)
	if err != nil {
		return err
	}
	defer templateFile.Close()
	templateFileContents, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("error reading content from file %v err: %v", templateFileName, err)
	}

	//create script file
	fileName := generatedDir + scriptFileName
	splitVal := strings.Split(fileName, "/")
	dirToCreate := strings.Join(splitVal[0:len(splitVal)-1], "/")
	//fmt.Println("dirToCreate : ", dirToCreate)

	if _, err := os.Stat(dirToCreate); os.IsNotExist(err) {
		err := os.MkdirAll(dirToCreate, 0755)
		if err != nil {
			return fmt.Errorf("error while creating %v directory, err: %v", dirToCreate, err)
		}
	}
	scriptFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while generating %v script : %v", scriptName, err)
	}
	defer scriptFile.Close()

	//make script executable
	err = scriptFile.Chmod(0755)
	if err != nil {
		log.Fatalf("could not make %v file executable, err : %v", fileName, err)
	}

	delims := []string{"{{", "}}"}
	if config["delims"] != nil {
		delimsParsed, parseOk := config["delims"].([]interface{})
		if parseOk && len(delimsParsed) == 2 {
			delims0, parseOk0 := delimsParsed[0].(string)
			delims1, parseOk1 := delimsParsed[1].(string)
			if parseOk0 && parseOk1 {
				delims[0] = delims0
				delims[1] = delims1
			} else {
				return fmt.Errorf("could not parse delims in config file")
			}
		} else {
			return fmt.Errorf("could not parse delims in config file")
		}
	}
	t := template.Must(template.New(templateFileName).
		Funcs(sprig.TxtFuncMap()).
		Funcs(writeFile()).
		Delims(delims[0], delims[1]).
		Parse(string(templateFileContents)))

	err = t.Execute(scriptFile, config)
	if err != nil {
		os.Remove(fileName)
		return fmt.Errorf("error while generating %v script : %v", scriptName, err)
	}

	fmt.Printf("generated script %-40v location: %v\n", scriptName, fileName)
	return nil

}

func writeFile() map[string]interface{} {
	funcMap := template.FuncMap{
		"writeFile": func(filename, fileData string) (string, error) {
			err := ioutil.WriteFile(filename, []byte(fileData), 0644)
			if err != nil {
				return "", fmt.Errorf("could not write to file %v, err : %v", filename, err)
			}
			return "", nil
		},
	}
	return funcMap
}
