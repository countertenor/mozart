package template

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/countertenor/mozart/statik"
	"github.com/rakyll/statik/fs"
	"gopkg.in/yaml.v2"
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
			fileExt := fileName[strings.LastIndex(fileName, "."):]
			scriptName := fileName
			if fileExt != ".yaml" {
				// fmt.Println("yaml file found!")
				scriptName = strings.TrimSuffix(fileName, fileExt)
			}
			templateInstance := templateInstance{
				scriptName:       scriptName,
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
	fileExt := scriptFileName[strings.LastIndex(scriptFileName, "."):]
	if fileExt == ".yaml" {
		config1 := make(map[string]interface{})
		err = yaml.Unmarshal(templateFileContents, &config1)
		if err != nil {
			return fmt.Errorf("error while unmarshalling yaml %v: %v", scriptFileName, err)
		}
		for key, val := range config1 {
			config[key] = val
		}
		return nil
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

	err = t.Execute(scriptFile, config)
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
