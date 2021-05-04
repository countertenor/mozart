package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/countertenor/mozart/static"
)

type templateInstance struct {
	scriptName             string
	scriptFileRelativePath string
	templateFilePath       string
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
			err := generateTemplate(template.scriptName, template.scriptFileRelativePath, template.templateFilePath, generatedDir, conf)
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

	// fmt.Println("dirToGenerate : ", filepath.Join(templateDir, dirToGenerate))
	err := static.Walk(static.ResourceType, filepath.Join(templateDir, dirToGenerate), func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			scriptRelativePath, err := static.GetRelativePath(path, templateDir)
			if err != nil {
				return err
			}
			fileName := info.Name()
			fileExt := filepath.Ext(fileName)
			templateInstance := templateInstance{
				scriptName:             strings.TrimSuffix(fileName, fileExt),
				scriptFileRelativePath: scriptRelativePath,
				templateFilePath:       path,
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

func generateTemplate(scriptName, scriptFileRelativePath, templateFilePath, generatedDir string, config map[string]interface{}) error {

	templateFile, err := static.OpenFileFromStaticFS(static.ResourceType, templateFilePath)
	if err != nil {
		return err
	}
	defer templateFile.Close()
	templateFileContents, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return fmt.Errorf("error reading content from file %v err: %v", templateFilePath, err)
	}

	//create script file
	fullFilePath := filepath.Join(generatedDir, scriptFileRelativePath)
	dirToCreate := filepath.Dir(fullFilePath)
	// fmt.Println("dirToCreate : ", dirToCreate)

	if _, err := os.Stat(dirToCreate); os.IsNotExist(err) {
		err := os.MkdirAll(dirToCreate, 0755)
		if err != nil {
			return fmt.Errorf("error while creating %v directory, err: %v", dirToCreate, err)
		}
	}

	//create script file
	scriptFile, err := os.Create(fullFilePath)
	if err != nil {
		return fmt.Errorf("error while generating %v script : %v", scriptName, err)
	}
	defer scriptFile.Close()

	//make script executable
	err = scriptFile.Chmod(0755)
	if err != nil {
		log.Fatalf("could not make %v file executable, err : %v", fullFilePath, err)
	}

	templatedContents, err := Templatize(config, templateFilePath, templateFileContents)
	if err != nil {
		return fmt.Errorf("error while templatizing %v script : %v", scriptName, err)
	}

	_, err = scriptFile.Write(templatedContents.Bytes())
	if err != nil {
		os.Remove(fullFilePath)
		return fmt.Errorf("error while writing to %v script : %v", scriptName, err)
	}

	fmt.Printf("generated script %-40v location: %v\n", scriptName, fullFilePath)
	return nil

}

//Templatize parses the file contents into a bytes.Buffer
// while applying templating changes
func Templatize(config map[string]interface{}, templateFilePath string, templateFileContents []byte) (bytes.Buffer, error) {
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
				return bytes.Buffer{}, fmt.Errorf("could not parse delims in config file")
			}
		} else {
			return bytes.Buffer{}, fmt.Errorf("could not parse delims in config file")
		}
	}

	t := template.Must(template.New(templateFilePath).
		Funcs(sprig.TxtFuncMap()).
		Funcs(writeFile()).
		Delims(delims[0], delims[1]).
		Parse(string(templateFileContents)))

	var templatedContent bytes.Buffer
	err := t.Execute(&templatedContent, config)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return templatedContent, nil
}

//write contents to file
func writeFile() map[string]interface{} {
	funcMap := template.FuncMap{
		"writeContentToFile": func(fileContent, filename string) (string, error) {
			err := ioutil.WriteFile(filename, []byte(fileContent), 0777)
			if err != nil {
				return "", fmt.Errorf("could not write to file %v, err : %v", filename, err)
			}
			return "", nil
		},
	}
	return funcMap
}
