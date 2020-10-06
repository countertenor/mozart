package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

func main() {

	// f, _ := os.Open("test1.tmpl")

	// t := template.Must(template.New("test1.tmpl").
	// 	Funcs(sprig.TxtFuncMap()).ParseFiles(f.Name()))

	// err = t.Execute("test1.sh", nil)

	// type Inventory struct {
	// 	Material string
	// 	Count    uint
	// }
	// sweaters := Inventory{"wool", 17}
	sweaters := make(map[string]interface{})
	getData(sweaters)

	// tmpl, err := template.New("test").Parse("{{.Values.Required.Value1}} items are made of {{.Material}}")
	tmpl, err := template.New("test1.tmpl").ParseFiles("test1.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters)
	if err != nil {
		panic(err)
	}
}

func getData(config map[string]interface{}) {
	fileData, _ := ioutil.ReadFile("data.yaml")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		return fmt.Errorf("no yaml file called %v found. Run 'init' to generate the sample config yaml file", filename)
	// 	}
	// 	return fmt.Errorf("error while reading conf file %v : %v", filename, err)
	// }

	// working
	// NOTE: the client-go library can only decode json, so we will first convert the yaml to json before unmarshaling
	// jsonData, err := yamlUtil.YAMLToJSON([]byte(fileData))
	// if err != nil {
	// 	// return errors.WithStackTrace(err)
	// 	fmt.Println("err1 : ", err)
	// }
	// fmt.Println("JSON : ", string(jsonData))

	// // var destinationObj interface{}
	// // destinationObj := make(map[string]interface{})
	// err = json.Unmarshal(jsonData, &config)
	// if err != nil {
	// 	// return errors.WithStackTrace(err)
	// 	fmt.Println("err2 : ", err)
	// }
	// // return nil

	// fmt.Println("destinationObj : ", config)

	//new
	err := yaml.Unmarshal(fileData, &config)
	if err != nil {
		// return errors.WithStackTrace(err)
		fmt.Println("err2 : ", err)
	}
}
