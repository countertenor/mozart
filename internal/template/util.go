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
	"github.com/prashantgupta24/mozart/internal/config"
	"github.com/prashantgupta24/mozart/statik"
	"github.com/rakyll/statik/fs"
)

type templateInstance struct {
	scriptName       string
	scriptFileName   string
	templateFileName string
}

//Generate conf files based on input yaml
func Generate(conf *config.Instance, dirToGenerate, templateDir, templateFileExt, generatedDir string) error {
	fmt.Println("")

	templatesToGenerate, err := getTemplatesToGenerate(dirToGenerate, templateDir, templateFileExt)
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

func getTemplatesToGenerate(dirToGenerate, templateDir, templateFileExt string) ([]templateInstance, error) {
	var templates []templateInstance

	statikFS, err := statik.GetStaticFS()
	if err != nil {
		return nil, err
	}

	err = fs.Walk(statikFS, templateDir+dirToGenerate, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == templateFileExt {
			fileName := info.Name()
			scriptFileName := strings.TrimPrefix(path, templateDir)
			templateInstance := templateInstance{
				scriptName:       strings.TrimSuffix(fileName, templateFileExt),
				scriptFileName:   strings.TrimSuffix(scriptFileName, templateFileExt) + ".sh",
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

func generateTemplate(scriptName, scriptFileName, templateFileName, generatedDir string, combinedStruct *config.Instance) error {

	statikFS, err := statik.GetStaticFS()
	if err != nil {
		return err
	}
	templateFile, err := statik.OpenFileFromStaticFS(statikFS, templateFileName)
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

	t := template.Must(template.New(templateFileName).
		Funcs(sprig.TxtFuncMap()).
		Funcs(getLastStringSplit()).
		Funcs(excludeFirstHost()).
		Funcs(getFirstHost()).
		Parse(string(templateFileContents)))

	err = t.Execute(scriptFile, combinedStruct)
	if err != nil {
		os.Remove(fileName)
		return fmt.Errorf("error while generating %v script : %v", scriptName, err)
	}

	fmt.Printf("generated script %-40v location: %v\n", scriptName, fileName)
	return nil

}

//send the first host
func getFirstHost() map[string]interface{} {
	funcMap := template.FuncMap{
		"getFirstHost": func(s string) string {
			return strings.Split(s, " ")[0]
		},
	}
	return funcMap
}

//sends the list of hosts excluding the first one
func excludeFirstHost() map[string]interface{} {
	funcMap := template.FuncMap{
		"excludeFirstHost": func(s string) string {
			return strings.Join(strings.Split(s, " ")[1:], " ")
		},
	}
	return funcMap
}

func getLastStringSplit() map[string]interface{} {
	funcMap := template.FuncMap{
		"split": func(s string) string {
			splitStrings := strings.Split(s, "/")
			return splitStrings[len(splitStrings)-1]
		},
	}
	return funcMap
}
